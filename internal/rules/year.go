package rules

import (
	"cmp"
	"slices"
	"time"

	"github.com/Optimus-Marketing/litcal/internal/dates"
	"github.com/Optimus-Marketing/litcal/internal/definitions"
	"github.com/Optimus-Marketing/litcal/internal/i18n"
	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

type yearBuilder struct {
	year              int
	catalog           definitions.Catalog
	names             i18n.Names
	easter            time.Time
	christmas         time.Time
	celebrations      []*liturgy.Celebration
	celebrationsByKey map[string]*liturgy.Celebration

	solemnities                    celebrationGroup
	feastsOfTheLord                celebrationGroup
	feasts                         celebrationGroup
	memorials                      celebrationGroup
	sundaysOfAdventLentAndEaster   celebrationGroup
	adventWeekdaysBeforeDecember17 celebrationGroup
	adventWeekdaysFromDecember17   celebrationGroup
	christmasWeekdays              celebrationGroup
	lentWeekdays                   celebrationGroup
}

func CelebrationsOfYear(year int, catalog definitions.Catalog, names i18n.Names) []liturgy.Celebration {
	builder := &yearBuilder{
		year:              year,
		catalog:           catalog,
		names:             names,
		easter:            dates.Easter(year),
		christmas:         dates.Date(year, time.December, 25),
		celebrationsByKey: make(map[string]*liturgy.Celebration),
	}
	return builder.build()
}

func (b *yearBuilder) build() []liturgy.Celebration {
	b.addProperOfTime()
	b.addFixedSolemnities()
	b.addFeastsOfTheLord()
	b.addOrdinarySundays()
	b.addFeasts()
	b.addAdventWeekdays()
	b.addChristmasOctaveWeekdays()
	b.addLentWeekdays()
	b.addImmaculateHeart()
	b.addMemorials(liturgy.RankMemorial)
	b.addDecreedCelebrations(liturgy.RankMemorial)
	b.addConversionOfSaintPaulInPaulineYear()
	b.applyRankDecrees()
	b.addMemorials(liturgy.RankOptionalMemorial)
	b.addDecreedCelebrations(liturgy.RankOptionalMemorial)
	b.applyRenameDecrees()
	b.addChristmasWeekdays()
	b.addWeekdaysAfterEpiphany()
	b.addEasterWeekdays()
	b.addOrdinaryTimeWeekdays()
	b.addSaturdayMemorialsOfMary()
	b.assignSeasons()
	b.addVigils()
	b.assignPsalterWeeks()
	b.assignCycles()
	return b.sortedCelebrations()
}

func (b *yearBuilder) newCelebration(definition definitions.Celebration, day time.Time) *liturgy.Celebration {
	return &liturgy.Celebration{
		Key:                   definition.Key,
		Name:                  b.names.Celebration(definition.Key),
		Date:                  day,
		Rank:                  definition.Rank,
		Colors:                slices.Clone(definition.Colors),
		Commons:               slices.Clone(definition.Commons),
		IsMovable:             definition.IsMovable,
		IsHolyDayOfObligation: definition.IsHolyDayOfObligation,
	}
}

func (b *yearBuilder) newCelebrationOnFixedDate(definition definitions.Celebration) *liturgy.Celebration {
	return b.newCelebration(definition, dates.Date(b.year, definition.Month, definition.Day))
}

func (b *yearBuilder) addFromProperOfTime(key string, day time.Time) *liturgy.Celebration {
	celebration := b.newCelebration(b.catalog.ProperOfTime[key], day)
	b.add(celebration)
	return celebration
}

func (b *yearBuilder) addWeekday(
	key, name string, day time.Time, color liturgy.Color, psalterWeek int,
) *liturgy.Celebration {
	weekday := &liturgy.Celebration{
		Key:         key,
		Name:        name,
		Date:        day,
		Rank:        liturgy.RankWeekday,
		Colors:      []liturgy.Color{color},
		IsMovable:   true,
		PsalterWeek: psalterWeek,
	}
	b.add(weekday)
	return weekday
}

func (b *yearBuilder) add(celebration *liturgy.Celebration) {
	b.insert(celebration)
	switch {
	case celebration.Rank >= liturgy.RankSolemnity:
		b.solemnities.add(celebration)
	case celebration.Rank == liturgy.RankFeastOfTheLord:
		b.feastsOfTheLord.add(celebration)
	case celebration.Rank == liturgy.RankFeast:
		b.feasts.add(celebration)
	case celebration.Rank == liturgy.RankMemorial:
		b.memorials.add(celebration)
	}
}

func (b *yearBuilder) insert(celebration *liturgy.Celebration) {
	if existing, found := b.celebrationsByKey[celebration.Key]; found {
		b.celebrations[slices.Index(b.celebrations, existing)] = celebration
	} else {
		b.celebrations = append(b.celebrations, celebration)
	}
	b.celebrationsByKey[celebration.Key] = celebration
}

func (b *yearBuilder) remove(key string) {
	b.solemnities.remove(key)
	b.feasts.remove(key)
	b.memorials.remove(key)
	if removed, found := b.celebrationsByKey[key]; found {
		delete(b.celebrationsByKey, key)
		b.celebrations = slices.DeleteFunc(b.celebrations, func(celebration *liturgy.Celebration) bool {
			return celebration == removed
		})
	}
}

// A memorial lowered to an optional memorial or a commemoration keeps occupying its date.
func (b *yearBuilder) changeRank(celebration *liturgy.Celebration, rank liturgy.Rank) {
	previous := celebration.Rank
	if previous == rank {
		return
	}
	celebration.Rank = rank
	switch {
	case rank >= liturgy.RankFeastOfTheLord:
		b.solemnities.add(celebration)
		b.feasts.remove(celebration.Key)
		b.memorials.remove(celebration.Key)
	case rank == liturgy.RankFeast:
		b.feasts.add(celebration)
		if previous > liturgy.RankFeast {
			b.solemnities.remove(celebration.Key)
		} else if previous == liturgy.RankMemorial {
			b.memorials.remove(celebration.Key)
		}
	case rank == liturgy.RankMemorial:
		b.memorials.add(celebration)
		if previous > liturgy.RankFeast {
			b.solemnities.remove(celebration.Key)
		} else if previous > liturgy.RankMemorial {
			b.feasts.remove(celebration.Key)
		}
	}
}

func (b *yearBuilder) dateOf(key string) time.Time {
	return b.celebrationsByKey[key].Date
}

func (b *yearBuilder) celebrationsOn(day time.Time) []*liturgy.Celebration {
	var celebrations []*liturgy.Celebration
	for _, celebration := range b.celebrations {
		if celebration.Date.Equal(day) {
			celebrations = append(celebrations, celebration)
		}
	}
	return celebrations
}

func (b *yearBuilder) hasSolemnityOrFeast(day time.Time) bool {
	return b.solemnities.hasCelebrationOn(day) || b.feastsOfTheLord.hasCelebrationOn(day) ||
		b.feasts.hasCelebrationOn(day)
}

func (b *yearBuilder) hasSolemnityFeastOrMemorial(day time.Time) bool {
	return b.hasSolemnityOrFeast(day) || b.memorials.hasCelebrationOn(day)
}

func (b *yearBuilder) sortedCelebrations() []liturgy.Celebration {
	ordered := slices.Clone(b.celebrations)
	slices.SortStableFunc(ordered, func(first, second *liturgy.Celebration) int {
		if order := first.Date.Compare(second.Date); order != 0 {
			return order
		}
		if first.IsVigil() != second.IsVigil() {
			if first.IsVigil() {
				return 1
			}
			return -1
		}
		return cmp.Compare(first.Rank, second.Rank)
	})

	celebrations := make([]liturgy.Celebration, len(ordered))
	for index, celebration := range ordered {
		celebrations[index] = *celebration
	}
	return celebrations
}

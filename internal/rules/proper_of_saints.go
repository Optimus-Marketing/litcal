package rules

import (
	"slices"
	"time"

	"github.com/Optimus-Marketing/litcal/internal/dates"
	"github.com/Optimus-Marketing/litcal/internal/definitions"
	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

func (b *yearBuilder) properOfSaintsWithRank(rank liturgy.Rank) []definitions.Celebration {
	var definitionsWithRank []definitions.Celebration
	for _, definition := range b.catalog.ProperOfSaints {
		if definition.Rank == rank && definition.IsInForce(b.year) {
			definitionsWithRank = append(definitionsWithRank, definition)
		}
	}
	return definitionsWithRank
}

func (b *yearBuilder) addFixedSolemnities() {
	b.addFromProperOfTime("MaryMotherOfGod", dates.Date(b.year, time.January, 1))
	for _, definition := range b.properOfSaintsWithRank(liturgy.RankSolemnity) {
		solemnity := b.newCelebrationOnFixedDate(definition)
		if coinciding := b.solemnities.celebrationOn(solemnity.Date); coinciding != nil {
			solemnity.Date = b.transferredDate(solemnity, coinciding)
		}
		b.add(solemnity)
	}
}

func (b *yearBuilder) transferredDate(solemnity, coinciding *liturgy.Celebration) time.Time {
	day := solemnity.Date
	nextDay := day.AddDate(0, 0, 1)
	switch {
	case solemnity.Key == "StJoseph" && dates.IsWithin(day, b.dateOf("PalmSun"), b.easter):
		return b.easter.AddDate(0, 0, -8)
	case solemnity.Key == "Annunciation" && dates.IsWithin(day, b.dateOf("PalmSun"), b.dateOf("Easter2")):
		return b.easter.AddDate(0, 0, 8)
	case slices.Contains([]string{"Annunciation", "StJoseph", "ImmaculateConception"}, solemnity.Key) &&
		b.sundaysOfAdventLentAndEaster.hasCelebrationOn(day):
		return nextDay
	case solemnity.Key == "StsPeterPaulAp" && coinciding.Key == "SacredHeart" &&
		!b.solemnities.hasCelebrationOn(nextDay):
		return nextDay
	case solemnity.Key == "NativityJohnBaptist" && coinciding.Key == "SacredHeart" &&
		!b.solemnities.hasCelebrationOn(coinciding.Date.AddDate(0, 0, -1)):
		return day.AddDate(0, 0, -1)
	}
	return day
}

func (b *yearBuilder) addFeasts() {
	for _, definition := range b.properOfSaintsWithRank(liturgy.RankFeast) {
		feast := b.newCelebrationOnFixedDate(definition)
		if !dates.IsSunday(feast.Date) && !b.hasSolemnityOrFeast(feast.Date) {
			b.add(feast)
		}
	}
}

func (b *yearBuilder) addImmaculateHeart() {
	immaculateHeart := b.newCelebration(b.catalog.ProperOfTime["ImmaculateHeart"], b.easter.AddDate(0, 0, 69))
	if !b.hasSolemnityFeastOrMemorial(immaculateHeart.Date) {
		b.add(immaculateHeart)
	}
}

func (b *yearBuilder) addMemorials(rank liturgy.Rank) {
	for _, definition := range b.properOfSaintsWithRank(rank) {
		memorial := b.newCelebrationOnFixedDate(definition)
		switch {
		case !dates.IsSunday(memorial.Date) && !b.hasSolemnityFeastOrMemorial(memorial.Date):
			b.add(memorial)
			if b.isPrivilegedWeekday(memorial.Date) {
				b.changeRank(memorial, liturgy.RankCommemoration)
			} else if memorial.Rank == liturgy.RankMemorial {
				b.replaceAdventWeekday(memorial)
			}
		case rank == liturgy.RankMemorial && b.isImmaculateHeartOn(memorial.Date):
			b.changeRank(b.celebrationsByKey["ImmaculateHeart"], liturgy.RankOptionalMemorial)
			memorial.Rank = liturgy.RankOptionalMemorial
			b.add(memorial)
		}
	}
}

func (b *yearBuilder) isPrivilegedWeekday(day time.Time) bool {
	return b.adventWeekdaysFromDecember17.hasCelebrationOn(day) || b.christmasWeekdays.hasCelebrationOn(day) ||
		b.lentWeekdays.hasCelebrationOn(day)
}

func (b *yearBuilder) isImmaculateHeartOn(day time.Time) bool {
	immaculateHeart := b.celebrationsByKey["ImmaculateHeart"]
	return immaculateHeart != nil && immaculateHeart.Date.Equal(day)
}

func (b *yearBuilder) replaceAdventWeekday(memorial *liturgy.Celebration) {
	if !memorial.Date.After(b.dateOf("Advent1")) || !memorial.Date.Before(dates.Date(b.year, time.December, 17)) {
		return
	}
	if weekday := b.adventWeekdaysBeforeDecember17.celebrationOn(memorial.Date); weekday != nil {
		memorial.PsalterWeek = weekday.PsalterWeek
		b.remove(weekday.Key)
	}
}

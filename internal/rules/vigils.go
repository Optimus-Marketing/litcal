package rules

import (
	"slices"
	"time"

	"github.com/Optimus-Marketing/litcal/internal/dates"
	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

var solemnitiesOfTheLordAndMary = []string{
	"Easter", "Christmas", "Ascension", "Pentecost", "Trinity", "CorpusChristi", "SacredHeart", "ChristKing",
	"MaryMotherOfGod", "Annunciation", "ImmaculateConception", "Assumption", "StJoseph", "NativityJohnBaptist",
}

func (b *yearBuilder) addVigils() {
	holyThursday := b.dateOf("HolyThurs")
	secondSundayOfEaster := b.dateOf("Easter2")
	for _, celebration := range slices.Clone(b.celebrations) {
		if !dates.IsSunday(celebration.Date) && celebration.Rank <= liturgy.RankFeast {
			continue
		}
		if !celebration.Date.Before(holyThursday) && celebration.Date.Before(secondSundayOfEaster) {
			continue
		}
		if celebration.Key != "AshWednesday" {
			b.addVigilFor(celebration)
		}
	}
}

func (b *yearBuilder) addVigilFor(celebration *liturgy.Celebration) {
	if !dates.IsSunday(celebration.Date) && celebration.Rank < liturgy.RankSolemnity {
		return
	}
	if celebration.Key == "AllSouls" || !b.canHaveVigil(celebration.Date) {
		celebration.HasVigilMass = false
		celebration.HasFirstVespers = false
		return
	}

	vigil := &liturgy.Celebration{
		Key:                   celebration.Key + "_vigil",
		Name:                  b.names.VigilMass(celebration.Name),
		Date:                  celebration.Date.AddDate(0, 0, -1),
		Rank:                  celebration.Rank,
		Colors:                celebration.Colors,
		Commons:               celebration.Commons,
		IsMovable:             celebration.IsMovable,
		Season:                celebration.Season,
		IsHolyDayOfObligation: celebration.IsHolyDayOfObligation,
		VigilFor:              celebration.Key,
	}
	b.insert(vigil)
	celebration.HasVigilMass = true
	celebration.HasFirstVespers = true
	celebration.HasSecondVespers = true

	coinciding := b.solemnities.celebrationOn(vigil.Date)
	if coinciding == nil {
		return
	}
	if celebration.Rank == liturgy.RankHigherSolemnity {
		if b.canHaveVigil(coinciding.Date) {
			giveVigilPrecedence(celebration, coinciding)
		}
		return
	}
	switch {
	case isSolemnityOfTheLordOrMary(coinciding.Key) && !isSolemnityOfTheLordOrMary(celebration.Key):
		celebration.HasVigilMass = false
		celebration.HasFirstVespers = false
		coinciding.HasSecondVespers = true
		b.remove(vigil.Key)
	case isSolemnityOfTheLordOrMary(celebration.Key) && !isSolemnityOfTheLordOrMary(coinciding.Key),
		b.year == 2022 && slices.Contains([]string{"SacredHeart", "Lent3", "Assumption"}, celebration.Key):
		giveVigilPrecedence(celebration, coinciding)
	}
}

func (b *yearBuilder) canHaveVigil(day time.Time) bool {
	isInHolyWeek := day.After(b.dateOf("PalmSun")) && day.Before(b.easter)
	isInEasterOctave := day.After(b.easter) && day.Before(b.dateOf("Easter2"))
	return !isInHolyWeek && !isInEasterOctave
}

func giveVigilPrecedence(celebration, coinciding *liturgy.Celebration) {
	celebration.HasVigilMass = true
	celebration.HasFirstVespers = true
	coinciding.HasSecondVespers = false
}

func isSolemnityOfTheLordOrMary(key string) bool {
	return slices.Contains(solemnitiesOfTheLordAndMary, key)
}

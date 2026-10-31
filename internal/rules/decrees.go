package rules

import (
	"time"

	"github.com/Optimus-Marketing/litcal/internal/dates"
	"github.com/Optimus-Marketing/litcal/internal/definitions"
	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

func (b *yearBuilder) addDecreedCelebrations(rank liturgy.Rank) {
	for _, decree := range b.catalog.Decrees {
		if decree.Action != definitions.CreateCelebration || decree.Rank != rank || !decree.IsInForce(b.year) {
			continue
		}
		celebration := b.newCelebration(decree.Celebration, b.decreedDate(decree))
		switch rank {
		case liturgy.RankMemorial:
			b.addDisplacingMemorial(celebration)
		case liturgy.RankOptionalMemorial:
			if !b.hasSolemnityFeastOrMemorial(celebration.Date) {
				b.add(celebration)
			}
		}
	}
}

func (b *yearBuilder) decreedDate(decree definitions.Decree) time.Time {
	if relative := decree.RelativeDate; relative != nil {
		return b.dateOf(relative.CelebrationKey).AddDate(0, 0, relative.DaysAfter)
	}
	return dates.Date(b.year, decree.Month, decree.Day)
}

func (b *yearBuilder) addDisplacingMemorial(memorial *liturgy.Celebration) {
	if b.hasSolemnityOrFeast(memorial.Date) {
		return
	}
	for _, displaced := range b.celebrationsOn(memorial.Date) {
		b.remove(displaced.Key)
	}
	b.add(memorial)
}

// During the Pauline Year the Conversion of Saint Paul fell on a Sunday and could be kept as an optional memorial.
func (b *yearBuilder) addConversionOfSaintPaulInPaulineYear() {
	if b.year != 2009 || b.celebrationsByKey["ConversionStPaul"] != nil {
		return
	}
	conversion := b.newCelebrationOnFixedDate(b.catalog.ProperOfSaintsCelebration("ConversionStPaul"))
	conversion.Rank = liturgy.RankOptionalMemorial
	b.add(conversion)
}

func (b *yearBuilder) applyRankDecrees() {
	for _, decree := range b.catalog.Decrees {
		if decree.Action != definitions.ChangeRank || !decree.IsInForce(b.year) {
			continue
		}
		if celebration := b.celebrationsByKey[decree.Key]; celebration != nil {
			b.changeRank(celebration, decree.Rank)
		}
	}
}

func (b *yearBuilder) applyRenameDecrees() {
	for _, decree := range b.catalog.Decrees {
		if decree.Action != definitions.RenameCelebration || !decree.IsInForce(b.year) {
			continue
		}
		if celebration := b.celebrationsByKey[decree.Key]; celebration != nil {
			celebration.Name = b.names.RenamedCelebration(decree.Key)
		}
	}
}

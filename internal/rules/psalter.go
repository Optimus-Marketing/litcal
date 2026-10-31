package rules

import (
	"time"

	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

func psalterWeekOf(week int) int {
	return (week-1)%4 + 1
}

func (b *yearBuilder) assignPsalterWeeks() {
	for _, celebration := range b.celebrations {
		if celebration.PsalterWeek != 0 {
			continue
		}
		switch {
		case celebration.IsVigil():
			if anticipated := b.celebrationsByKey[celebration.VigilFor]; anticipated != nil {
				celebration.PsalterWeek = anticipated.PsalterWeek
			}
		case celebration.Rank == liturgy.RankCommemoration || celebration.Rank == liturgy.RankOptionalMemorial:
			if weekday := b.firstWeekdayOn(celebration.Date); weekday != nil {
				celebration.PsalterWeek = weekday.PsalterWeek
			}
		}
	}
}

func (b *yearBuilder) firstWeekdayOn(day time.Time) *liturgy.Celebration {
	for _, celebration := range b.celebrations {
		if celebration.Rank == liturgy.RankWeekday && celebration.Date.Equal(day) {
			return celebration
		}
	}
	return nil
}

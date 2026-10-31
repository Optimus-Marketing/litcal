package rules

import (
	"github.com/Optimus-Marketing/litcal/internal/dates"
	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

var (
	sundayCycles  = [...]liturgy.Cycle{liturgy.CycleA, liturgy.CycleB, liturgy.CycleC}
	weekdayCycles = [...]liturgy.Cycle{liturgy.CycleI, liturgy.CycleII}
)

func (b *yearBuilder) assignCycles() {
	for _, celebration := range b.celebrations {
		day := celebration.Date
		if celebration.IsVigil() {
			day = day.AddDate(0, 0, 1)
		}
		liturgicalYear := dates.LiturgicalYear(day)
		celebration.SundayCycle = sundayCycles[(liturgicalYear-1)%3]
		celebration.WeekdayCycle = weekdayCycles[(liturgicalYear-1)%2]
	}
}

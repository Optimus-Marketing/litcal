package rules

import (
	"time"

	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

func (b *yearBuilder) assignSeasons() {
	for _, celebration := range b.celebrations {
		celebration.Season = b.seasonOf(celebration)
	}
}

func (b *yearBuilder) seasonOf(celebration *liturgy.Celebration) liturgy.Season {
	day := celebration.Date
	holyThursday := b.dateOf("HolyThurs")
	switch {
	case !day.Before(b.dateOf("Advent1")) && day.Before(b.christmas):
		return liturgy.SeasonAdvent
	case !day.Before(b.christmas) || !day.After(b.dateOf("BaptismLord")):
		return liturgy.SeasonChristmas
	case !day.Before(b.dateOf("AshWednesday")) && day.Before(holyThursday):
		return liturgy.SeasonLent
	case !day.Before(holyThursday) && day.Before(b.easter):
		if celebration.Key == "HolyThursChrism" {
			return liturgy.SeasonLent
		}
		return liturgy.SeasonEasterTriduum
	case !day.Before(b.easter) && !day.After(b.dateOf("Pentecost")):
		return liturgy.SeasonEaster
	default:
		return liturgy.SeasonOrdinaryTime
	}
}

func (b *yearBuilder) isInOrdinaryTime(day time.Time) bool {
	return (day.After(b.dateOf("BaptismLord")) && day.Before(b.dateOf("AshWednesday"))) ||
		(day.After(b.dateOf("Pentecost")) && day.Before(b.dateOf("Advent1")))
}

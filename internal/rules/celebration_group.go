package rules

import (
	"slices"
	"time"

	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

type celebrationGroup []*liturgy.Celebration

func (g *celebrationGroup) add(celebration *liturgy.Celebration) {
	for index, member := range *g {
		if member.Key == celebration.Key {
			(*g)[index] = celebration
			return
		}
	}
	*g = append(*g, celebration)
}

func (g *celebrationGroup) remove(key string) {
	*g = slices.DeleteFunc(*g, func(member *liturgy.Celebration) bool { return member.Key == key })
}

func (g celebrationGroup) celebrationOn(day time.Time) *liturgy.Celebration {
	for _, member := range g {
		if member.Date.Equal(day) {
			return member
		}
	}
	return nil
}

func (g celebrationGroup) hasCelebrationOn(day time.Time) bool {
	return g.celebrationOn(day) != nil
}

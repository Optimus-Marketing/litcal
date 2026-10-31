package rules

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/Optimus-Marketing/litcal/internal/definitions"
)

var mondayToSaturday = []time.Weekday{
	time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday,
}

func CelebrationKeys(catalog definitions.Catalog) []string {
	keys := make(map[string]bool)
	for key := range catalog.ProperOfTime {
		keys[key] = true
	}
	for _, definition := range catalog.ProperOfSaints {
		keys[definition.Key] = true
	}
	for _, decree := range catalog.Decrees {
		if decree.Action == definitions.CreateCelebration {
			keys[decree.Key] = true
		}
	}
	for _, weekday := range mondayToSaturday {
		addKeysForWeeks(keys, adventWeekdayKeyFormat, 1, 3, weekday)
		addKeysForWeeks(keys, lentWeekdayKeyFormat, 1, 5, weekday)
		addKeysForWeeks(keys, easterWeekdayKeyFormat, 2, 7, weekday)
		addKeysForWeeks(keys, ordinaryWeekdayKeyFormat, 1, 34, weekday)
	}
	for _, weekday := range []time.Weekday{time.Thursday, time.Friday, time.Saturday} {
		keys[fmt.Sprintf(afterAshWednesdayKeyFormat, weekday)] = true
	}
	addKeysForDays(keys, adventDecemberWeekdayKeyFormat, 17, 24)
	addKeysForDays(keys, christmasOctaveWeekdayKeyFormat, 26, 31)
	addKeysForDays(keys, christmasWeekdayKeyFormat, 2, 5)
	addKeysForDays(keys, afterEpiphanyWeekdayKeyFormat, 7, 12)
	return slices.Sorted(maps.Keys(keys))
}

func addKeysForWeeks(keys map[string]bool, keyFormat string, firstWeek, lastWeek int, weekday time.Weekday) {
	for week := firstWeek; week <= lastWeek; week++ {
		keys[fmt.Sprintf(keyFormat, week, weekday)] = true
	}
}

func addKeysForDays(keys map[string]bool, keyFormat string, firstDay, lastDay int) {
	for day := firstDay; day <= lastDay; day++ {
		keys[fmt.Sprintf(keyFormat, day)] = true
	}
}

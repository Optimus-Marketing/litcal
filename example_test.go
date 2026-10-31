package litcal_test

import (
	"fmt"
	"log"
	"time"

	"github.com/Optimus-Marketing/litcal"
)

func ExampleCalendar_Year() {
	calendar, err := litcal.New(litcal.English)
	if err != nil {
		log.Fatal(err)
	}
	celebrations, err := calendar.Year(2026)
	if err != nil {
		log.Fatal(err)
	}
	christmasEve := time.Date(2026, time.December, 24, 0, 0, 0, 0, time.UTC)
	for _, celebration := range celebrations {
		if celebration.Date.Equal(christmasEve) {
			fmt.Printf("%s (%s, psalter week %d)\n", celebration.Name, celebration.Season, celebration.PsalterWeek)
		}
	}
	// Output:
	// Thursday of the 4th Week of Advent (ADVENT, psalter week 4)
	// Christmas Vigil Mass (CHRISTMAS, psalter week 0)
}

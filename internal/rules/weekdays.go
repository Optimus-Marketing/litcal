package rules

import (
	"fmt"
	"time"

	"github.com/Optimus-Marketing/litcal/internal/dates"
	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

const (
	adventWeekdayKeyFormat          = "AdventWeekday%d%s"
	adventDecemberWeekdayKeyFormat  = "AdventWeekdayDec%d"
	christmasOctaveWeekdayKeyFormat = "ChristmasWeekdayDec%d"
	afterAshWednesdayKeyFormat      = "%sAfterAshWednesday"
	lentWeekdayKeyFormat            = "LentWeekday%d%s"
	christmasWeekdayKeyFormat       = "ChristmasWeekdayJan%d"
	afterEpiphanyWeekdayKeyFormat   = "DayAfterEpiphanyJan%d"
	easterWeekdayKeyFormat          = "EasterWeekday%d%s"
	ordinaryWeekdayKeyFormat        = "OrdWeekday%d%s"
	saturdayMemorialKeyFormat       = "SatMemBVM%d"
)

func (b *yearBuilder) addAdventWeekdays() {
	firstSunday := b.dateOf("Advent1")
	for day := firstSunday.AddDate(0, 0, 1); day.Before(b.christmas); day = day.AddDate(0, 0, 1) {
		if dates.IsSunday(day) || b.hasSolemnityFeastOrMemorial(day) {
			continue
		}
		week := dates.DaysBetween(firstSunday, day)/7 + 1
		key := fmt.Sprintf(adventWeekdayKeyFormat, week, day.Weekday())
		if day.Month() == time.December && day.Day() >= 17 {
			key = fmt.Sprintf(adventDecemberWeekdayKeyFormat, day.Day())
		}
		weekday := b.addWeekday(key, b.names.AdventWeekday(day, week), day, liturgy.ColorPurple, psalterWeekOf(week))
		if day.Day() < 17 {
			b.adventWeekdaysBeforeDecember17.add(weekday)
		} else {
			b.adventWeekdaysFromDecember17.add(weekday)
		}
	}
}

func (b *yearBuilder) addChristmasOctaveWeekdays() {
	for dayOfOctave := 2; dayOfOctave <= 7; dayOfOctave++ {
		day := b.christmas.AddDate(0, 0, dayOfOctave-1)
		if dates.IsSunday(day) || b.hasSolemnityFeastOrMemorial(day) {
			continue
		}
		key := fmt.Sprintf(christmasOctaveWeekdayKeyFormat, day.Day())
		name := b.names.ChristmasOctaveDay(dayOfOctave)
		b.christmasWeekdays.add(b.addWeekday(key, name, day, liturgy.ColorWhite, 0))
	}
}

func (b *yearBuilder) addLentWeekdays() {
	firstSunday := b.dateOf("Lent1")
	for day := b.dateOf("AshWednesday").AddDate(0, 0, 1); day.Before(b.dateOf("PalmSun")); day = day.AddDate(0, 0, 1) {
		if dates.IsSunday(day) || b.solemnities.hasCelebrationOn(day) {
			continue
		}
		if !day.After(firstSunday) {
			key := fmt.Sprintf(afterAshWednesdayKeyFormat, day.Weekday())
			b.addWeekday(key, b.names.WeekdayAfterAshWednesday(day), day, liturgy.ColorPurple, 0)
			continue
		}
		week := dates.DaysBetween(firstSunday, day)/7 + 1
		key := fmt.Sprintf(lentWeekdayKeyFormat, week, day.Weekday())
		name := b.names.LentWeekday(day, week)
		b.lentWeekdays.add(b.addWeekday(key, name, day, liturgy.ColorPurple, psalterWeekOf(week)))
	}
}

func (b *yearBuilder) addChristmasWeekdays() {
	for day := dates.Date(b.year, time.January, 2); day.Before(b.dateOf("Epiphany")); day = day.AddDate(0, 0, 1) {
		if dates.IsSunday(day) || b.hasSolemnityFeastOrMemorial(day) {
			continue
		}
		key := fmt.Sprintf(christmasWeekdayKeyFormat, day.Day())
		b.addWeekday(key, b.names.ChristmasWeekday(day), day, liturgy.ColorWhite, 0)
	}
}

func (b *yearBuilder) addWeekdaysAfterEpiphany() {
	for day := b.dateOf("Epiphany").AddDate(0, 0, 1); day.Before(b.dateOf("BaptismLord")); day = day.AddDate(0, 0, 1) {
		if b.hasSolemnityFeastOrMemorial(day) {
			continue
		}
		key := fmt.Sprintf(afterEpiphanyWeekdayKeyFormat, day.Day())
		b.addWeekday(key, b.names.ChristmasWeekday(day), day, liturgy.ColorWhite, 0)
	}
}

func (b *yearBuilder) addEasterWeekdays() {
	for day := b.easter.AddDate(0, 0, 1); day.Before(b.dateOf("Pentecost")); day = day.AddDate(0, 0, 1) {
		if dates.IsSunday(day) || b.hasSolemnityFeastOrMemorial(day) {
			continue
		}
		week := dates.DaysBetween(b.easter, day)/7 + 1
		key := fmt.Sprintf(easterWeekdayKeyFormat, week, day.Weekday())
		b.addWeekday(key, b.names.EasterWeekday(day, week), day, liturgy.ColorWhite, psalterWeekOf(week))
	}
}

func (b *yearBuilder) addOrdinaryTimeWeekdays() {
	baptism := b.dateOf("BaptismLord")
	firstSunday := dates.NextSunday(baptism)
	week := 1
	for day := baptism.AddDate(0, 0, 1); day.Before(b.dateOf("AshWednesday")); day = day.AddDate(0, 0, 1) {
		if b.hasSolemnityFeastOrMemorial(day) {
			continue
		}
		if day.After(firstSunday) {
			week = dates.DaysBetween(firstSunday, day)/7 + 2
		}
		b.addOrdinaryTimeWeekday(day, week)
	}

	firstSundayOfAdvent := b.dateOf("Advent1")
	for day := b.dateOf("Pentecost").AddDate(0, 0, 1); day.Before(firstSundayOfAdvent); day = day.AddDate(0, 0, 1) {
		if b.hasSolemnityFeastOrMemorial(day) {
			continue
		}
		b.addOrdinaryTimeWeekday(day, 34-dates.DaysBetween(day, firstSundayOfAdvent)/7)
	}
}

func (b *yearBuilder) addOrdinaryTimeWeekday(day time.Time, week int) {
	key := fmt.Sprintf(ordinaryWeekdayKeyFormat, week, day.Weekday())
	b.addWeekday(key, b.names.OrdinaryTimeWeekday(day, week), day, liturgy.ColorGreen, psalterWeekOf(week))
}

func (b *yearBuilder) addSaturdayMemorialsOfMary() {
	firstOfJanuary := dates.Date(b.year, time.January, 1)
	firstSaturday := firstOfJanuary.AddDate(0, 0, (int(time.Saturday)-int(firstOfJanuary.Weekday())+7)%7)
	saturdayNumber := 0
	for saturday := firstSaturday; saturday.Year() == b.year; saturday = saturday.AddDate(0, 0, 7) {
		if !b.isInOrdinaryTime(saturday) || b.hasSolemnityFeastOrMemorial(saturday) {
			continue
		}
		saturdayNumber++
		b.add(&liturgy.Celebration{
			Key:       fmt.Sprintf(saturdayMemorialKeyFormat, saturdayNumber),
			Name:      b.names.SaturdayMemorialOfMary(),
			Date:      saturday,
			Rank:      liturgy.RankOptionalMemorial,
			Colors:    []liturgy.Color{liturgy.ColorWhite},
			Commons:   []string{"Blessed Virgin Mary"},
			IsMovable: true,
		})
	}
}

package rules

import (
	"fmt"
	"time"

	"github.com/Optimus-Marketing/litcal/internal/dates"
	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

var easterOctaveDayKeys = []string{
	"MonOctaveEaster", "TueOctaveEaster", "WedOctaveEaster", "ThuOctaveEaster", "FriOctaveEaster", "SatOctaveEaster",
}

func (b *yearBuilder) addProperOfTime() {
	fourthSundayOfAdvent := dates.PreviousSunday(b.christmas)

	b.addFromProperOfTime("HolyThurs", b.easter.AddDate(0, 0, -3))
	b.addFromProperOfTime("GoodFri", b.easter.AddDate(0, 0, -2))
	b.addFromProperOfTime("EasterVigil", b.easter.AddDate(0, 0, -1))
	b.addFromProperOfTime("Easter", b.easter)

	b.addFromProperOfTime("Christmas", b.christmas)
	b.addFromProperOfTime("Epiphany", dates.Date(b.year, time.January, 6))
	for day := 2; day <= 5; day++ {
		if sunday := dates.Date(b.year, time.January, day); dates.IsSunday(sunday) {
			b.addFromProperOfTime("Christmas2", sunday)
			break
		}
	}

	b.addFromProperOfTime("Ascension", b.easter.AddDate(0, 0, 39))
	b.addSundayOfSeason("Easter", 7, b.easter.AddDate(0, 0, 42))
	b.addFromProperOfTime("Pentecost", b.easter.AddDate(0, 0, 49))

	for week := 1; week <= 4; week++ {
		b.addSundayOfSeason("Advent", week, fourthSundayOfAdvent.AddDate(0, 0, -7*(4-week)))
	}
	for week := 1; week <= 5; week++ {
		b.addSundayOfSeason("Lent", week, b.easter.AddDate(0, 0, -7*(7-week)))
	}
	b.addFromProperOfTime("PalmSun", b.easter.AddDate(0, 0, -7))
	for week := 2; week <= 6; week++ {
		b.addSundayOfSeason("Easter", week, b.easter.AddDate(0, 0, 7*(week-1)))
	}
	b.addFromProperOfTime("Trinity", b.easter.AddDate(0, 0, 56))
	b.addFromProperOfTime("CorpusChristi", b.easter.AddDate(0, 0, 60))
	if b.year >= 2000 {
		secondSundayOfEaster := b.celebrationsByKey["Easter2"]
		secondSundayOfEaster.Name = b.names.DivineMercySunday(secondSundayOfEaster.Name)
	}

	b.addFromProperOfTime("AshWednesday", b.easter.AddDate(0, 0, -46))
	b.addFromProperOfTime("MonHolyWeek", b.easter.AddDate(0, 0, -6))
	b.addFromProperOfTime("TueHolyWeek", b.easter.AddDate(0, 0, -5))
	b.addFromProperOfTime("WedHolyWeek", b.easter.AddDate(0, 0, -4))
	b.addFromProperOfTime("HolyThursChrism", b.easter.AddDate(0, 0, -3))
	for offset, key := range easterOctaveDayKeys {
		b.addFromProperOfTime(key, b.easter.AddDate(0, 0, offset+1))
	}

	b.addFromProperOfTime("SacredHeart", b.easter.AddDate(0, 0, 68))
	b.addFromProperOfTime("ChristKing", fourthSundayOfAdvent.AddDate(0, 0, -28))
}

func (b *yearBuilder) addSundayOfSeason(season string, week int, day time.Time) {
	sunday := b.addFromProperOfTime(fmt.Sprintf("%s%d", season, week), day)
	sunday.PsalterWeek = psalterWeekOf(week)
	b.sundaysOfAdventLentAndEaster.add(sunday)
}

func (b *yearBuilder) addFeastsOfTheLord() {
	b.addFromProperOfTime("BaptismLord", dates.NextSunday(dates.Date(b.year, time.January, 6)))

	for _, definition := range b.properOfSaintsWithRank(liturgy.RankFeastOfTheLord) {
		b.add(b.newCelebrationOnFixedDate(definition))
	}

	holyFamily := dates.NextSunday(b.christmas)
	if dates.IsSunday(b.christmas) {
		holyFamily = dates.Date(b.year, time.December, 30)
	}
	b.addFromProperOfTime("HolyFamily", holyFamily)
}

func (b *yearBuilder) addOrdinarySundays() {
	ashWednesday := b.dateOf("AshWednesday")
	sunday := dates.NextSunday(b.dateOf("BaptismLord"))
	for week := 2; sunday.Before(ashWednesday); week++ {
		b.addOrdinarySunday(week, sunday)
		sunday = sunday.AddDate(0, 0, 7)
	}

	trinity := b.dateOf("Trinity")
	sunday = b.dateOf("ChristKing").AddDate(0, 0, -7)
	for week := 33; sunday.After(trinity); week-- {
		b.addOrdinarySunday(week, sunday)
		sunday = sunday.AddDate(0, 0, -7)
	}
}

func (b *yearBuilder) addOrdinarySunday(week int, sunday time.Time) {
	if b.solemnities.hasCelebrationOn(sunday) || b.feastsOfTheLord.hasCelebrationOn(sunday) {
		return
	}
	celebration := b.addFromProperOfTime(fmt.Sprintf("OrdSunday%d", week), sunday)
	celebration.PsalterWeek = psalterWeekOf(week)
}

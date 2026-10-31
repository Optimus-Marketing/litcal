package dates

import "time"

func Date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func IsSunday(day time.Time) bool {
	return day.Weekday() == time.Sunday
}

func NextSunday(day time.Time) time.Time {
	return day.AddDate(0, 0, 7-int(day.Weekday()))
}

func PreviousSunday(day time.Time) time.Time {
	daysBack := int(day.Weekday())
	if daysBack == 0 {
		daysBack = 7
	}
	return day.AddDate(0, 0, -daysBack)
}

func DaysBetween(from, to time.Time) int {
	return int(to.Sub(from).Hours() / 24)
}

func IsWithin(day, from, to time.Time) bool {
	return !day.Before(from) && !day.After(to)
}

func Easter(year int) time.Time {
	yearInMetonicCycle := year % 19
	century, yearOfCentury := year/100, year%100
	leapCenturies, centuryRemainder := century/4, century%4
	moonCorrection := (century - (century+8)/25 + 1) / 3
	epact := (19*yearInMetonicCycle + century - leapCenturies - moonCorrection + 15) % 30
	leapYears, yearRemainder := yearOfCentury/4, yearOfCentury%4
	weekday := (32 + 2*centuryRemainder + 2*leapYears - epact - yearRemainder) % 7
	correction := (yearInMetonicCycle + 11*epact + 22*weekday) / 451
	month := (epact + weekday - 7*correction + 114) / 31
	day := (epact+weekday-7*correction+114)%31 + 1
	return Date(year, time.Month(month), day)
}

func FirstSundayOfAdvent(year int) time.Time {
	return PreviousSunday(Date(year, time.December, 25)).AddDate(0, 0, -21)
}

func LiturgicalYear(day time.Time) int {
	if day.Before(FirstSundayOfAdvent(day.Year())) {
		return day.Year()
	}
	return day.Year() + 1
}

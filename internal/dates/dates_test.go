package dates

import (
	"testing"
	"time"
)

func TestEaster(t *testing.T) {
	for year, want := range map[int]string{
		1970: "1970-03-29",
		2000: "2000-04-23",
		2008: "2008-03-23",
		2019: "2019-04-21",
		2024: "2024-03-31",
		2025: "2025-04-20",
		2026: "2026-04-05",
		2027: "2027-03-28",
		2038: "2038-04-25",
		2285: "2285-03-22",
	} {
		if got := Easter(year).Format(time.DateOnly); got != want {
			t.Errorf("Easter(%d) = %s, want %s", year, got, want)
		}
	}
}

func TestFirstSundayOfAdvent(t *testing.T) {
	for year, want := range map[int]string{
		2022: "2022-11-27",
		2023: "2023-12-03",
		2024: "2024-12-01",
		2025: "2025-11-30",
		2026: "2026-11-29",
	} {
		if got := FirstSundayOfAdvent(year).Format(time.DateOnly); got != want {
			t.Errorf("FirstSundayOfAdvent(%d) = %s, want %s", year, got, want)
		}
	}
}

func TestLiturgicalYearStartsOnTheFirstSundayOfAdvent(t *testing.T) {
	for day, want := range map[string]int{
		"2026-01-01": 2026,
		"2026-11-28": 2026,
		"2026-11-29": 2027,
		"2026-12-31": 2027,
	} {
		date, _ := time.Parse(time.DateOnly, day)
		if got := LiturgicalYear(date); got != want {
			t.Errorf("LiturgicalYear(%s) = %d, want %d", day, got, want)
		}
	}
}

func TestSundays(t *testing.T) {
	thursday, sunday := Date(2026, time.October, 29), Date(2026, time.November, 1)
	tests := []struct {
		name string
		got  time.Time
		want string
	}{
		{"next Sunday after a Thursday", NextSunday(thursday), "2026-11-01"},
		{"next Sunday after a Sunday", NextSunday(sunday), "2026-11-08"},
		{"previous Sunday before a Thursday", PreviousSunday(thursday), "2026-10-25"},
		{"previous Sunday before a Sunday", PreviousSunday(sunday), "2026-10-25"},
	}
	for _, test := range tests {
		if got := test.got.Format(time.DateOnly); got != test.want {
			t.Errorf("%s = %s, want %s", test.name, got, test.want)
		}
	}
	if !IsSunday(sunday) || IsSunday(thursday) {
		t.Error("IsSunday mistakes the weekday")
	}
}

func TestDaysBetweenAndIsWithin(t *testing.T) {
	from, to := Date(2026, time.March, 1), Date(2026, time.April, 5)
	if days := DaysBetween(from, to); days != 35 {
		t.Errorf("DaysBetween = %d, want 35", days)
	}
	if !IsWithin(from, from, to) || !IsWithin(to, from, to) || IsWithin(to.AddDate(0, 0, 1), from, to) {
		t.Error("IsWithin must include both ends and nothing after them")
	}
}

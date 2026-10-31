package litcal_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Optimus-Marketing/litcal"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata/golden")

var goldenYears = []struct {
	language litcal.Language
	year     int
}{
	{litcal.English, 2008},
	{litcal.English, 2009},
	{litcal.English, 2016},
	{litcal.English, 2021},
	{litcal.English, 2022},
	{litcal.English, 2024},
	{litcal.English, 2026},
	{litcal.Spanish, 2026},
	{litcal.Spanish, 2027},
}

var rankNames = map[litcal.Rank]string{
	litcal.RankWeekday:          "weekday",
	litcal.RankCommemoration:    "commemoration",
	litcal.RankOptionalMemorial: "optional-memorial",
	litcal.RankMemorial:         "memorial",
	litcal.RankFeast:            "feast",
	litcal.RankFeastOfTheLord:   "feast-of-the-lord",
	litcal.RankSolemnity:        "solemnity",
	litcal.RankHigherSolemnity:  "higher-solemnity",
}

func newCalendar(t testing.TB, language litcal.Language) *litcal.Calendar {
	t.Helper()
	calendar, err := litcal.New(language)
	if err != nil {
		t.Fatalf("New(%q): %v", language, err)
	}
	return calendar
}

func celebrationsOfYear(t testing.TB, calendar *litcal.Calendar, year int) []litcal.Celebration {
	t.Helper()
	celebrations, err := calendar.Year(year)
	if err != nil {
		t.Fatalf("Year(%d): %v", year, err)
	}
	return celebrations
}

func find(t *testing.T, celebrations []litcal.Celebration, key string) litcal.Celebration {
	t.Helper()
	index := slices.IndexFunc(celebrations, func(celebration litcal.Celebration) bool { return celebration.Key == key })
	if index < 0 {
		t.Fatalf("no celebration %q", key)
	}
	return celebrations[index]
}

func TestYearMatchesGolden(t *testing.T) {
	for _, golden := range goldenYears {
		name := fmt.Sprintf("%s-%d", golden.language, golden.year)
		t.Run(name, func(t *testing.T) {
			got := toTSV(celebrationsOfYear(t, newCalendar(t, golden.language), golden.year))
			path := filepath.Join("testdata", "golden", name+".tsv")
			if *update {
				if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("%v (run go test -update to create it)", err)
			}
			if got != string(want) {
				t.Errorf("%s differs from the calendar; run go test -update after an intended change\n%s",
					path, firstDifference(string(want), got))
			}
		})
	}
}

func TestNotableRules(t *testing.T) {
	english := newCalendar(t, litcal.English)
	tests := []struct {
		name string
		year int
		key  string
		date string
		rank litcal.Rank
	}{
		{"Saint Joseph leaves Holy Week", 2008, "StJoseph", "2008-03-15", litcal.RankSolemnity},
		{"Annunciation after the Easter octave", 2024, "Annunciation", "2024-04-08", litcal.RankSolemnity},
		{"Immaculate Conception off Advent Sunday", 2024, "ImmaculateConception", "2024-12-09", litcal.RankSolemnity},
		{"Saint John the Baptist yields to the Sacred Heart", 2022, "NativityJohnBaptist", "2022-06-23", litcal.RankSolemnity},
		{"Saint Mary Magdalene is a memorial before 2016", 2015, "StMaryMagdalene", "2015-07-22", litcal.RankMemorial},
		{"Saint Mary Magdalene is a feast since 2016", 2016, "StMaryMagdalene", "2016-07-22", litcal.RankFeast},
		{"Conversion of Saint Paul in the Pauline Year", 2009, "ConversionStPaul", "2009-01-25", litcal.RankOptionalMemorial},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			celebration := find(t, celebrationsOfYear(t, english, test.year), test.key)
			if date := celebration.Date.Format(time.DateOnly); date != test.date || celebration.Rank != test.rank {
				t.Errorf("%s on %s as %s, want %s as %s",
					test.key, date, rankNames[celebration.Rank], test.date, rankNames[test.rank])
			}
		})
	}
}

func TestRenamedCelebrations(t *testing.T) {
	english := newCalendar(t, litcal.English)
	for year, want := range map[int]string{2020: "Saint Martha", 2021: "Saints Martha, Mary and Lazarus"} {
		if name := find(t, celebrationsOfYear(t, english, year), "StMartha").Name; name != want {
			t.Errorf("StMartha in %d is %q, want %q", year, name, want)
		}
	}
}

func TestPsalterWeekAndCycles(t *testing.T) {
	weekday := find(t, celebrationsOfYear(t, newCalendar(t, litcal.Spanish), 2026), "OrdWeekday30Thursday")
	if date := weekday.Date.Format(time.DateOnly); date != "2026-10-29" {
		t.Fatalf("OrdWeekday30Thursday on %s, want 2026-10-29", date)
	}
	if weekday.PsalterWeek != 2 || weekday.SundayCycle != litcal.CycleA || weekday.WeekdayCycle != litcal.CycleII {
		t.Errorf("psalter week %d and cycles %s/%s, want 2 and A/II",
			weekday.PsalterWeek, weekday.SundayCycle, weekday.WeekdayCycle)
	}
}

func TestYearRange(t *testing.T) {
	calendar := newCalendar(t, litcal.English)
	for _, year := range []int{1969, 10000} {
		if _, err := calendar.Year(year); !errors.Is(err, litcal.ErrYearOutOfRange) {
			t.Errorf("Year(%d) error = %v, want ErrYearOutOfRange", year, err)
		}
	}
	for _, year := range []int{1970, 9999} {
		if _, err := calendar.Year(year); err != nil {
			t.Errorf("Year(%d): %v", year, err)
		}
	}
}

func TestNewRejectsAnUnsupportedLanguage(t *testing.T) {
	if _, err := litcal.New("fr"); !errors.Is(err, litcal.ErrUnsupportedLanguage) {
		t.Errorf("New(fr) error = %v, want ErrUnsupportedLanguage", err)
	}
}

func TestEveryYearIsComplete(t *testing.T) {
	if testing.Short() {
		t.Skip("builds every year from 1970 to 9999")
	}
	for _, language := range []litcal.Language{litcal.English, litcal.Spanish} {
		calendar := newCalendar(t, language)
		celebrationKeys := calendar.CelebrationKeys()
		celebrationCommons := calendar.CelebrationCommons()
		for year := 1970; year <= 9999; year++ {
			celebrations := celebrationsOfYear(t, calendar, year)
			if err := checkYear(year, celebrations, celebrationKeys, celebrationCommons); err != nil {
				t.Fatalf("%s %d: %v", language, year, err)
			}
		}
	}
}

func TestCalendarIsSafeForConcurrentUse(t *testing.T) {
	calendar := newCalendar(t, litcal.English)
	want := toTSV(celebrationsOfYear(t, calendar, 2026))
	var group sync.WaitGroup
	results := make([]string, 8)
	for index := range results {
		group.Add(1)
		go func() {
			defer group.Done()
			celebrations, _ := calendar.Year(2026)
			results[index] = toTSV(celebrations)
		}()
	}
	group.Wait()
	for _, got := range results {
		if got != want {
			t.Fatal("concurrent calls returned different calendars")
		}
	}
}

func checkYear(year int, celebrations []litcal.Celebration, celebrationKeys, celebrationCommons []string) error {
	daysWithCelebration := make(map[string]bool)
	seenKeys := make(map[string]bool)
	for index, celebration := range celebrations {
		if seenKeys[celebration.Key] {
			return fmt.Errorf("%s appears twice", celebration.Key)
		}
		seenKeys[celebration.Key] = true
		if index > 0 && celebration.Date.Before(celebrations[index-1].Date) {
			return fmt.Errorf("%s is out of order", celebration.Key)
		}
		if celebration.Name == "" || len(celebration.Colors) == 0 {
			return fmt.Errorf("%s has no name or no color", celebration.Key)
		}
		if celebration.PsalterWeek < 0 || celebration.PsalterWeek > 4 {
			return fmt.Errorf("%s has psalter week %d", celebration.Key, celebration.PsalterWeek)
		}
		if celebration.IsVigil() {
			continue
		}
		daysWithCelebration[celebration.Date.Format(time.DateOnly)] = true
		_, isListed := slices.BinarySearch(celebrationKeys, celebration.Key)
		if !isListed && !strings.HasPrefix(celebration.Key, "SatMemBVM") {
			return fmt.Errorf("CelebrationKeys does not list %s", celebration.Key)
		}
		for _, common := range celebration.Commons {
			if _, isListed := slices.BinarySearch(celebrationCommons, common); !isListed && common != "Proper" {
				return fmt.Errorf("CelebrationCommons does not list %q of %s", common, celebration.Key)
			}
		}
	}
	for day := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC); day.Year() == year; day = day.AddDate(0, 0, 1) {
		if !daysWithCelebration[day.Format(time.DateOnly)] {
			return fmt.Errorf("no celebration on %s", day.Format(time.DateOnly))
		}
	}
	return nil
}

func toTSV(celebrations []litcal.Celebration) string {
	var lines strings.Builder
	lines.WriteString("# date\tkey\trank\tcolors\tcommons\tseason\tpsalter\tcycles\tvigil for\tflags\tname\n")
	for _, celebration := range celebrations {
		colors := make([]string, len(celebration.Colors))
		for index, color := range celebration.Colors {
			colors[index] = string(color)
		}
		fmt.Fprintf(&lines, "%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s/%s\t%s\t%s\t%s\n",
			celebration.Date.Format(time.DateOnly), celebration.Key, rankNames[celebration.Rank],
			strings.Join(colors, ","), strings.Join(celebration.Commons, "|"), celebration.Season,
			celebration.PsalterWeek, celebration.SundayCycle, celebration.WeekdayCycle, celebration.VigilFor,
			flagsOf(celebration), celebration.Name)
	}
	return lines.String()
}

func flagsOf(celebration litcal.Celebration) string {
	var flags []string
	for _, candidate := range []struct {
		name  string
		isSet bool
	}{
		{"movable", celebration.IsMovable},
		{"vigil-mass", celebration.HasVigilMass},
		{"first-vespers", celebration.HasFirstVespers},
		{"second-vespers", celebration.HasSecondVespers},
		{"obligation", celebration.IsHolyDayOfObligation},
	} {
		if candidate.isSet {
			flags = append(flags, candidate.name)
		}
	}
	return strings.Join(flags, ",")
}

func firstDifference(want, got string) string {
	wantLines, gotLines := strings.Split(want, "\n"), strings.Split(got, "\n")
	for index := range max(len(wantLines), len(gotLines)) {
		var wantLine, gotLine string
		if index < len(wantLines) {
			wantLine = wantLines[index]
		}
		if index < len(gotLines) {
			gotLine = gotLines[index]
		}
		if wantLine != gotLine {
			return fmt.Sprintf("line %d\nwant: %s\ngot:  %s", index+1, wantLine, gotLine)
		}
	}
	return ""
}

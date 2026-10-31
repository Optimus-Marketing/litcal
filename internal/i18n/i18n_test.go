package i18n

import (
	"errors"
	"maps"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"
)

var misusedWords = map[string]*regexp.Regexp{
	"en": regexp.MustCompile(`\b(bishops?|priests?|deacon|abbot|popes?|martyrs?|virgin|religious|apostles?|evangelist|doctors?|hermit|monk)\b`),
	"es": regexp.MustCompile(`^En (la|el) |\b(Obispos?|Presbíteros?|Sacerdote|Diácono|Abad|Papas?|Mártir|Mártires|Religios[oa]|Apóstol|Apóstoles|Evangelista|Doctora?|Doctores|Ermitaño|Monje|Esposo)\b|, Virgen( y|$)`),
}

func load(t *testing.T, language string) Names {
	t.Helper()
	names, err := Load(language)
	if err != nil {
		t.Fatalf("Load(%q): %v", language, err)
	}
	return names
}

func TestLoadRejectsAnUnknownLanguage(t *testing.T) {
	if _, err := Load("fr"); !errors.Is(err, ErrUnknownLanguage) {
		t.Errorf("Load(fr) error = %v, want ErrUnknownLanguage", err)
	}
}

func TestEveryLanguageNamesTheSameCelebrations(t *testing.T) {
	english, spanish := load(t, "en"), load(t, "es")
	if !slices.Equal(slices.Sorted(maps.Keys(english.celebrations)), slices.Sorted(maps.Keys(spanish.celebrations))) {
		t.Error("English and Spanish name different celebrations")
	}
	if !slices.Equal(slices.Sorted(maps.Keys(english.renamed)), slices.Sorted(maps.Keys(spanish.renamed))) {
		t.Error("English and Spanish rename different celebrations")
	}
}

func TestNames(t *testing.T) {
	english, spanish := load(t, "en"), load(t, "es")
	thursday := time.Date(2026, time.December, 24, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		got, want string
	}{
		{english.Celebration("Easter"), "Easter Sunday"},
		{english.AdventWeekday(thursday, 4), "Thursday of the 4th Week of Advent"},
		{english.WeekdayAfterAshWednesday(thursday), "Thursday after Ash Wednesday"},
		{english.VigilMass("Christmas"), "Christmas Vigil Mass"},
		{english.ChristmasOctaveDay(2), "2nd Day of the Octave of Christmas"},
		{english.RenamedCelebration("StMartha"), "Saints Martha, Mary and Lazarus"},
		{spanish.VigilMass("Navidad"), "Navidad (misa vespertina)"},
		{spanish.ChristmasOctaveDay(6), "Día VI dentro de la Octava de la Natividad del Señor"},
		{spanish.OrdinaryTimeWeekday(thursday, 21), "Jueves de la XXI semana del tiempo ordinario"},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("got %q, want %q", test.got, test.want)
		}
	}
}

func TestNamesAreWellFormed(t *testing.T) {
	for language, misused := range misusedWords {
		for _, name := range everyName(load(t, language)) {
			first, _ := utf8.DecodeRuneInString(name)
			switch {
			case name != strings.TrimSpace(name) || strings.Contains(name, "  ") || strings.Contains(name, " ,"):
				t.Errorf("%s name %q has stray spaces", language, name)
			case strings.HasSuffix(name, "."):
				t.Errorf("%s name %q ends with a period", language, name)
			case unicode.IsLower(first):
				t.Errorf("%s name %q starts in lowercase", language, name)
			case strings.ContainsAny(name, "{}"):
				t.Errorf("%s name %q has an unfilled placeholder", language, name)
			case misused.MatchString(name):
				t.Errorf("%s name %q writes %q", language, name, misused.FindString(name))
			}
		}
	}
}

func everyName(names Names) []string {
	all := slices.Concat(slices.Collect(maps.Values(names.celebrations)), slices.Collect(maps.Values(names.renamed)))
	for weekday := range 7 {
		day := time.Date(2026, time.November, 1+weekday, 0, 0, 0, 0, time.UTC)
		all = append(all, names.WeekdayAfterAshWednesday(day), names.ChristmasWeekday(day))
		for week := 1; week <= len(names.weekOrdinals); week++ {
			all = append(all, names.AdventWeekday(day, week), names.LentWeekday(day, week),
				names.EasterWeekday(day, week), names.OrdinaryTimeWeekday(day, week))
		}
	}
	for dayOfOctave := 1; dayOfOctave <= len(names.dayOrdinals); dayOfOctave++ {
		all = append(all, names.ChristmasOctaveDay(dayOfOctave))
	}
	return append(all,
		names.SaturdayMemorialOfMary(),
		names.VigilMass(names.Celebration("Christmas")),
		names.DivineMercySunday(names.Celebration("Easter2")),
	)
}

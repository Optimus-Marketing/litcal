package i18n

import (
	"errors"
	"maps"
	"slices"
	"testing"
	"time"
)

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
		{spanish.VigilMass("Navidad"), "Navidad Misa de la Víspera"},
		{spanish.ChristmasOctaveDay(2), "Segundo Día de la Octava de Navidad"},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("got %q, want %q", test.got, test.want)
		}
	}
}

package rules

import (
	"slices"
	"strings"
	"testing"

	"github.com/Optimus-Marketing/litcal/internal/definitions"
)

func TestPsalterWeekRepeatsEveryFourWeeks(t *testing.T) {
	for week, want := range map[int]int{1: 1, 4: 4, 5: 1, 30: 2, 34: 2} {
		if got := psalterWeekOf(week); got != want {
			t.Errorf("psalterWeekOf(%d) = %d, want %d", week, got, want)
		}
	}
}

func TestCelebrationKeys(t *testing.T) {
	catalog, err := definitions.Load()
	if err != nil {
		t.Fatal(err)
	}
	keys := CelebrationKeys(catalog)
	if !slices.IsSorted(keys) || len(slices.Compact(slices.Clone(keys))) != len(keys) {
		t.Fatal("CelebrationKeys must be sorted and unique")
	}
	for _, key := range []string{"Easter", "StJoseph", "MaryMotherChurch", "AdventWeekdayDec24", "OrdWeekday34Saturday"} {
		if _, isListed := slices.BinarySearch(keys, key); !isListed {
			t.Errorf("CelebrationKeys does not list %s", key)
		}
	}
	for _, key := range keys {
		if strings.HasSuffix(key, "_vigil") || strings.HasPrefix(key, "SatMemBVM") {
			t.Errorf("CelebrationKeys lists %s", key)
		}
	}
}

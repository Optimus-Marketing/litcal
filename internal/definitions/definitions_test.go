package definitions

import (
	"slices"
	"testing"

	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

func TestLoad(t *testing.T) {
	catalog, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if easter := catalog.ProperOfTime["Easter"]; easter.Rank != liturgy.RankHigherSolemnity {
		t.Errorf("Easter rank = %d, want %d", easter.Rank, liturgy.RankHigherSolemnity)
	}
	if len(catalog.ProperOfSaints) == 0 || len(catalog.Decrees) == 0 {
		t.Error("the proper of saints and the decrees must not be empty")
	}
	if magdalene := catalog.ProperOfSaintsCelebration("StMaryMagdalene"); magdalene.Key != "StMaryMagdalene" {
		t.Error("ProperOfSaintsCelebration does not find Saint Mary Magdalene")
	}
	if unknown := catalog.ProperOfSaintsCelebration("Unknown"); unknown.Key != "" {
		t.Error("ProperOfSaintsCelebration returns a celebration for an unknown key")
	}
}

func TestIsInForce(t *testing.T) {
	celebration := Celebration{InForceSince: 2016, InForceUntil: 2020}
	for year, want := range map[int]bool{2015: false, 2016: true, 2020: true, 2021: false} {
		if got := celebration.IsInForce(year); got != want {
			t.Errorf("IsInForce(%d) = %v, want %v", year, got, want)
		}
	}
	if !(Celebration{}).IsInForce(9999) {
		t.Error("a celebration without years of force is always in force")
	}
}

func TestCelebrationCommonsLeaveOutProper(t *testing.T) {
	catalog, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	commons := catalog.CelebrationCommons()
	if !slices.IsSorted(commons) || slices.Contains(commons, properCommon) {
		t.Errorf("CelebrationCommons must be sorted and leave out %q: %v", properCommon, commons)
	}
	if !slices.Contains(commons, "Blessed Virgin Mary") {
		t.Error(`CelebrationCommons does not list "Blessed Virgin Mary"`)
	}
}

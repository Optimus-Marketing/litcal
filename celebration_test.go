package litcal

import (
	"reflect"
	"testing"

	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

func TestCelebrationHasTheFieldsTheRulesBuild(t *testing.T) {
	public, built := reflect.TypeFor[Celebration](), reflect.TypeFor[liturgy.Celebration]()
	if public.NumField() != built.NumField() {
		t.Fatalf("Celebration has %d fields and the rules build %d", public.NumField(), built.NumField())
	}
	for index := range public.NumField() {
		if public.Field(index).Name != built.Field(index).Name {
			t.Errorf("field %d is %s, the rules build %s", index, public.Field(index).Name, built.Field(index).Name)
		}
	}
}

func TestValuesMatchTheOnesTheRulesUse(t *testing.T) {
	ranks := map[Rank]liturgy.Rank{
		RankWeekday:          liturgy.RankWeekday,
		RankCommemoration:    liturgy.RankCommemoration,
		RankOptionalMemorial: liturgy.RankOptionalMemorial,
		RankMemorial:         liturgy.RankMemorial,
		RankFeast:            liturgy.RankFeast,
		RankFeastOfTheLord:   liturgy.RankFeastOfTheLord,
		RankSolemnity:        liturgy.RankSolemnity,
		RankHigherSolemnity:  liturgy.RankHigherSolemnity,
	}
	for rank, built := range ranks {
		if Rank(built) != rank {
			t.Errorf("rank %d is %d in the rules", rank, built)
		}
	}

	strings := map[string]string{
		string(ColorWhite):          string(liturgy.ColorWhite),
		string(ColorRed):            string(liturgy.ColorRed),
		string(ColorGreen):          string(liturgy.ColorGreen),
		string(ColorPurple):         string(liturgy.ColorPurple),
		string(ColorRose):           string(liturgy.ColorRose),
		string(SeasonAdvent):        string(liturgy.SeasonAdvent),
		string(SeasonChristmas):     string(liturgy.SeasonChristmas),
		string(SeasonLent):          string(liturgy.SeasonLent),
		string(SeasonEasterTriduum): string(liturgy.SeasonEasterTriduum),
		string(SeasonEaster):        string(liturgy.SeasonEaster),
		string(SeasonOrdinaryTime):  string(liturgy.SeasonOrdinaryTime),
		string(CycleA):              string(liturgy.CycleA),
		string(CycleB):              string(liturgy.CycleB),
		string(CycleC):              string(liturgy.CycleC),
		string(CycleI):              string(liturgy.CycleI),
		string(CycleII):             string(liturgy.CycleII),
	}
	for public, built := range strings {
		if public != built {
			t.Errorf("%q is %q in the rules", public, built)
		}
	}
}

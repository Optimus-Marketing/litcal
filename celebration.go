package litcal

import (
	"time"

	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

// Celebration is a liturgical day, or the vigil Mass of one, in the General Roman Calendar.
// PsalterWeek is 0 when the four-week psalter does not apply.
type Celebration struct {
	Key                   string
	Name                  string
	Date                  time.Time
	Rank                  Rank
	Colors                []Color
	Commons               []string
	IsMovable             bool
	Season                Season
	PsalterWeek           int
	SundayCycle           Cycle
	WeekdayCycle          Cycle
	VigilFor              string
	HasVigilMass          bool
	HasFirstVespers       bool
	HasSecondVespers      bool
	IsHolyDayOfObligation bool
}

// IsVigil reports whether the celebration is the vigil Mass of the celebration named by VigilFor.
func (c Celebration) IsVigil() bool {
	return c.VigilFor != ""
}

// Rank orders celebrations by liturgical precedence.
type Rank int

const (
	RankWeekday Rank = iota
	RankCommemoration
	RankOptionalMemorial
	RankMemorial
	RankFeast
	RankFeastOfTheLord
	RankSolemnity
	RankHigherSolemnity
)

type Color string

const (
	ColorWhite  Color = "white"
	ColorRed    Color = "red"
	ColorGreen  Color = "green"
	ColorPurple Color = "purple"
	ColorRose   Color = "rose"
)

type Season string

const (
	SeasonAdvent        Season = "ADVENT"
	SeasonChristmas     Season = "CHRISTMAS"
	SeasonLent          Season = "LENT"
	SeasonEasterTriduum Season = "EASTER_TRIDUUM"
	SeasonEaster        Season = "EASTER"
	SeasonOrdinaryTime  Season = "ORDINARY_TIME"
)

// Cycle is a cycle of the Lectionary: A, B or C for Sundays and I or II for weekdays.
type Cycle string

const (
	CycleA  Cycle = "A"
	CycleB  Cycle = "B"
	CycleC  Cycle = "C"
	CycleI  Cycle = "I"
	CycleII Cycle = "II"
)

func newCelebration(celebration liturgy.Celebration) Celebration {
	return Celebration{
		Key:                   celebration.Key,
		Name:                  celebration.Name,
		Date:                  celebration.Date,
		Rank:                  Rank(celebration.Rank),
		Colors:                colorsOf(celebration.Colors),
		Commons:               celebration.Commons,
		IsMovable:             celebration.IsMovable,
		Season:                Season(celebration.Season),
		PsalterWeek:           celebration.PsalterWeek,
		SundayCycle:           Cycle(celebration.SundayCycle),
		WeekdayCycle:          Cycle(celebration.WeekdayCycle),
		VigilFor:              celebration.VigilFor,
		HasVigilMass:          celebration.HasVigilMass,
		HasFirstVespers:       celebration.HasFirstVespers,
		HasSecondVespers:      celebration.HasSecondVespers,
		IsHolyDayOfObligation: celebration.IsHolyDayOfObligation,
	}
}

func colorsOf(colors []liturgy.Color) []Color {
	if colors == nil {
		return nil
	}
	converted := make([]Color, len(colors))
	for index, color := range colors {
		converted[index] = Color(color)
	}
	return converted
}

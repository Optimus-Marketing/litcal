package liturgy

import "time"

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

func (c Celebration) IsVigil() bool {
	return c.VigilFor != ""
}

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

type Cycle string

const (
	CycleA  Cycle = "A"
	CycleB  Cycle = "B"
	CycleC  Cycle = "C"
	CycleI  Cycle = "I"
	CycleII Cycle = "II"
)

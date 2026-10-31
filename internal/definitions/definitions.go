package definitions

import (
	"embed"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/Optimus-Marketing/litcal/internal/liturgy"
)

//go:embed data/*.json
var data embed.FS

const properCommon = "Proper"

type Celebration struct {
	Key                   string          `json:"key"`
	Month                 time.Month      `json:"month"`
	Day                   int             `json:"day"`
	Rank                  liturgy.Rank    `json:"rank"`
	Colors                []liturgy.Color `json:"colors"`
	Commons               []string        `json:"commons"`
	IsMovable             bool            `json:"movable"`
	IsHolyDayOfObligation bool            `json:"holyDayOfObligation"`
	InForceSince          int             `json:"since"`
	InForceUntil          int             `json:"until"`
}

func (c Celebration) IsInForce(year int) bool {
	return year >= c.InForceSince && (c.InForceUntil == 0 || year <= c.InForceUntil)
}

type DecreeAction string

const (
	CreateCelebration DecreeAction = "create"
	ChangeRank        DecreeAction = "rank"
	RenameCelebration DecreeAction = "rename"
)

type Decree struct {
	Celebration
	Action       DecreeAction  `json:"action"`
	RelativeDate *RelativeDate `json:"after"`
}

type RelativeDate struct {
	CelebrationKey string `json:"key"`
	DaysAfter      int    `json:"days"`
}

type Catalog struct {
	ProperOfTime   map[string]Celebration
	ProperOfSaints []Celebration
	Decrees        []Decree
}

func Load() (Catalog, error) {
	var properOfTime, properOfSaints []Celebration
	var decrees []Decree
	if err := decode("data/proper_of_time.json", &properOfTime); err != nil {
		return Catalog{}, err
	}
	if err := decode("data/proper_of_saints.json", &properOfSaints); err != nil {
		return Catalog{}, err
	}
	if err := decode("data/decrees.json", &decrees); err != nil {
		return Catalog{}, err
	}

	properOfTimeByKey := make(map[string]Celebration, len(properOfTime))
	for _, celebration := range properOfTime {
		properOfTimeByKey[celebration.Key] = celebration
	}
	return Catalog{ProperOfTime: properOfTimeByKey, ProperOfSaints: properOfSaints, Decrees: decrees}, nil
}

func (c Catalog) ProperOfSaintsCelebration(key string) Celebration {
	for _, celebration := range c.ProperOfSaints {
		if celebration.Key == key {
			return celebration
		}
	}
	return Celebration{}
}

func (c Catalog) CelebrationCommons() []string {
	commons := make(map[string]bool)
	for _, celebration := range c.ProperOfSaints {
		for _, common := range celebration.Commons {
			commons[common] = true
		}
	}
	for _, decree := range c.Decrees {
		for _, common := range decree.Commons {
			commons[common] = true
		}
	}
	delete(commons, properCommon)
	return slices.Sorted(maps.Keys(commons))
}

func decode(path string, target any) error {
	content, err := data.ReadFile(path)
	if err != nil {
		return fmt.Errorf("litcal: read %s: %w", path, err)
	}
	if err := json.Unmarshal(content, target); err != nil {
		return fmt.Errorf("litcal: decode %s: %w", path, err)
	}
	return nil
}

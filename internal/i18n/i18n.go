package i18n

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"
)

//go:embed data/*.json
var data embed.FS

var ErrUnknownLanguage = errors.New("i18n: unknown language")

type Names struct {
	celebrations map[string]string
	renamed      map[string]string
	weekdays     [7]string
	weekOrdinals []string
	dayOrdinals  []string
	patterns     patterns
}

type patterns struct {
	AdventWeekday      string `json:"adventWeekday"`
	ChristmasOctaveDay string `json:"christmasOctaveDay"`
	ChristmasWeekday   string `json:"christmasWeekday"`
	AfterAshWednesday  string `json:"afterAshWednesday"`
	LentWeekday        string `json:"lentWeekday"`
	EasterWeekday      string `json:"easterWeekday"`
	OrdinaryWeekday    string `json:"ordinaryWeekday"`
	SaturdayMemorial   string `json:"saturdayMemorial"`
	VigilMass          string `json:"vigilMass"`
	DivineMercySunday  string `json:"divineMercySunday"`
}

func Load(language string) (Names, error) {
	path := "data/" + language + ".json"
	content, err := data.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Names{}, fmt.Errorf("%w: %q", ErrUnknownLanguage, language)
	}
	if err != nil {
		return Names{}, fmt.Errorf("litcal: read %s: %w", path, err)
	}
	var file struct {
		Celebrations map[string]string `json:"names"`
		Renamed      map[string]string `json:"renamed"`
		Weekdays     [7]string         `json:"weekdays"`
		WeekOrdinals []string          `json:"weekOrdinals"`
		DayOrdinals  []string          `json:"dayOrdinals"`
		Patterns     patterns          `json:"patterns"`
	}
	if err := json.Unmarshal(content, &file); err != nil {
		return Names{}, fmt.Errorf("litcal: decode %s: %w", path, err)
	}
	return Names{
		celebrations: file.Celebrations,
		renamed:      file.Renamed,
		weekdays:     file.Weekdays,
		weekOrdinals: file.WeekOrdinals,
		dayOrdinals:  file.DayOrdinals,
		patterns:     file.Patterns,
	}, nil
}

func (n Names) Celebration(key string) string {
	return n.celebrations[key]
}

func (n Names) RenamedCelebration(key string) string {
	return n.renamed[key]
}

func (n Names) AdventWeekday(day time.Time, week int) string {
	return n.withWeekdayAndWeek(n.patterns.AdventWeekday, day, week)
}

func (n Names) LentWeekday(day time.Time, week int) string {
	return n.withWeekdayAndWeek(n.patterns.LentWeekday, day, week)
}

func (n Names) EasterWeekday(day time.Time, week int) string {
	return n.withWeekdayAndWeek(n.patterns.EasterWeekday, day, week)
}

func (n Names) OrdinaryTimeWeekday(day time.Time, week int) string {
	return n.withWeekdayAndWeek(n.patterns.OrdinaryWeekday, day, week)
}

func (n Names) WeekdayAfterAshWednesday(day time.Time) string {
	return n.withWeekday(n.patterns.AfterAshWednesday, day)
}

func (n Names) ChristmasWeekday(day time.Time) string {
	return n.withWeekday(n.patterns.ChristmasWeekday, day)
}

func (n Names) ChristmasOctaveDay(dayOfOctave int) string {
	return strings.ReplaceAll(n.patterns.ChristmasOctaveDay, "{ordinal}", n.dayOrdinals[dayOfOctave-1])
}

func (n Names) SaturdayMemorialOfMary() string {
	return n.patterns.SaturdayMemorial
}

func (n Names) VigilMass(celebrationName string) string {
	return strings.ReplaceAll(n.patterns.VigilMass, "{name}", celebrationName)
}

func (n Names) DivineMercySunday(secondSundayOfEasterName string) string {
	return strings.ReplaceAll(n.patterns.DivineMercySunday, "{name}", secondSundayOfEasterName)
}

func (n Names) withWeekdayAndWeek(pattern string, day time.Time, week int) string {
	replacer := strings.NewReplacer("{weekday}", n.weekdays[day.Weekday()], "{ordinal}", n.weekOrdinals[week-1])
	return replacer.Replace(pattern)
}

func (n Names) withWeekday(pattern string, day time.Time) string {
	return strings.ReplaceAll(pattern, "{weekday}", n.weekdays[day.Weekday()])
}

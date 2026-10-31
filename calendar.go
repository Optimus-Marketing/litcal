// Package litcal computes the General Roman Calendar of the Roman rite for the years 1970 to 9999.
package litcal

import (
	"errors"
	"fmt"

	"github.com/Optimus-Marketing/litcal/internal/definitions"
	"github.com/Optimus-Marketing/litcal/internal/i18n"
	"github.com/Optimus-Marketing/litcal/internal/rules"
)

const (
	firstYear = 1970
	lastYear  = 9999
)

// ErrYearOutOfRange is returned for a year the calendar does not cover.
var ErrYearOutOfRange = errors.New("litcal: year out of range")

var ErrUnsupportedLanguage = errors.New("litcal: unsupported language")

// Language selects the language of the celebration names.
type Language string

const (
	English Language = "en"
	Spanish Language = "es"
)

// Calendar computes the General Roman Calendar with the celebration names in one language.
type Calendar struct {
	catalog definitions.Catalog
	names   i18n.Names
}

// New returns a calendar whose celebration names are in the given language.
func New(language Language) (*Calendar, error) {
	catalog, err := definitions.Load()
	if err != nil {
		return nil, err
	}
	names, err := i18n.Load(string(language))
	if errors.Is(err, i18n.ErrUnknownLanguage) {
		return nil, fmt.Errorf("%w %q", ErrUnsupportedLanguage, language)
	}
	if err != nil {
		return nil, err
	}
	return &Calendar{catalog: catalog, names: names}, nil
}

// Year returns the celebrations of a civil year sorted by date and rank. It includes the vigil Mass
// of January 1, celebrated on December 31 of the previous year.
func (c *Calendar) Year(year int) ([]Celebration, error) {
	if year < firstYear || year > lastYear {
		return nil, fmt.Errorf("%w: %d", ErrYearOutOfRange, year)
	}
	built := rules.CelebrationsOfYear(year, c.catalog, c.names)
	celebrations := make([]Celebration, len(built))
	for index, celebration := range built {
		celebrations[index] = newCelebration(celebration)
	}
	return celebrations, nil
}

// CelebrationKeys returns, sorted, the key of every celebration the calendar can place on a day in any year. It
// leaves out the vigil Masses, which belong to the celebration they precede, and the Saturday memorials of Mary,
// whose keys only number the Saturdays of a year.
func (c *Calendar) CelebrationKeys() []string {
	return rules.CelebrationKeys(c.catalog)
}

// CelebrationCommons returns, sorted, every common a celebration can use. "Proper", which marks a celebration with
// texts of its own, is not a common and is left out.
func (c *Calendar) CelebrationCommons() []string {
	return c.catalog.CelebrationCommons()
}

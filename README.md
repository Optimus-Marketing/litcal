# litcal

The General Roman Calendar in Go: which celebrations fall on each day from 1970 to 9999, with their rank, colors,
season, psalter week, Lectionary cycles and vigil Masses, and their names in English and Spanish. It only uses the
standard library and needs Go 1.24 or later.

## Installation

```sh
go get github.com/Optimus-Marketing/litcal
```

## Usage

```go
import "github.com/Optimus-Marketing/litcal"

calendar, err := litcal.New(litcal.Spanish)
if err != nil {
	return err
}
celebrations, err := calendar.Year(2026)
```

`Year` returns the celebrations of a civil year sorted by date and rank, plus the vigil Mass of January 1 on the
previous December 31. A `Calendar` is safe for concurrent use. `CelebrationKeys` and `CelebrationCommons` list every
celebration key and every common, for example to keep a registry of them in a database.

## Layout

| Path | Contents |
|---|---|
| `calendar.go`, `celebration.go` | The public API |
| `internal/rules` | How a year is built: proper of time, proper of saints, weekdays, decrees, vigils, seasons, psalter weeks and cycles |
| `internal/definitions` | The definition of each celebration and the decrees that changed them, in `data/` |
| `internal/i18n` | The names and naming patterns of each language, in `data/` |
| `internal/dates` | Easter, Advent, Sundays and the liturgical year |
| `internal/liturgy` | The celebration the rules build |
| `testdata/golden` | The expected calendar of several years, one line per celebration |

A new language needs its file in `internal/i18n/data` and a `Language` constant.

## Development

```sh
make test      # tests with the race detector, without the sweep of every year
make test-all  # every test, including every year from 1970 to 9999
make lint      # gofmt and go vet
make golden    # rewrite testdata/golden after an intended change
```

## Source

The rules and data are derived from [LiturgicalCalendarAPI](https://github.com/Liturgical-Calendar/LiturgicalCalendarAPI)
(Apache License 2.0). The Spanish celebration names follow the Calendario Litúrgico-Pastoral of the Spanish Episcopal
Conference, and some English names were corrected to the Liturgical Calendar of the United States Conference of Catholic
Bishops. Apart from names, its English calendars from 1970 to 2068 match this package except for two decrees it declares
but does not apply:

- Saint Mary Magdalene is a feast since 2016.
- July 29 is the memorial of Saints Martha, Mary and Lazarus since 2021.

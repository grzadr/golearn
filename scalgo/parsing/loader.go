// ReadFileContent reads the content of a Markdown file and returns it as a byte slice.
package parsing

import (
	"fmt"
	"strconv"
	"strings"
)

type RecordUnit interface {
	ConvertToBaseUnit(value float64) float64
	ConvertFromBaseUnit(value float64) float64
	getUnit() int
}

type TimeUnits int

const (
	Nanosecond TimeUnits = iota
	Microsecond
	Millisecond
	Second
	Minute
	Hour
	Day
	Week
	Month
	Year
	Decade
	Century
	Millennium
)

var timeUnitMap = map[string]TimeUnits{
	// Nanosecond
	"nanosecond":  Nanosecond,
	"nanoseconds": Nanosecond,
	"ns":          Nanosecond,

	// Microsecond
	"microsecond":  Microsecond,
	"microseconds": Microsecond,
	"μs":           Microsecond,
	"us":           Microsecond,

	// Millisecond
	"millisecond":  Millisecond,
	"milliseconds": Millisecond,
	"ms":           Millisecond,

	// Second
	"second":  Second,
	"seconds": Second,
	"sec":     Second,
	"s":       Second,

	// Minute
	"minute":  Minute,
	"minutes": Minute,
	"min":     Minute,
	"m":       Minute,

	// Hour
	"hour":  Hour,
	"hours": Hour,
	"hr":    Hour,
	"h":     Hour,

	// Day
	"day":  Day,
	"days": Day,
	"d":    Day,

	// Week
	"week":  Week,
	"weeks": Week,
	"wk":    Week,
	"w":     Week,

	// Month
	"month":  Month,
	"months": Month,
	"mo":     Month,

	// Year
	"year":  Year,
	"years": Year,
	"yr":    Year,
	"y":     Year,

	// Decade
	"decade":  Decade,
	"decades": Decade,

	// Century
	"century":   Century,
	"centuries": Century,

	// Millennium
	"millennium": Millennium,
	"millennia":  Millennium,
}

// ParseTimeUnit converts a string to a TimeUnits constant
func parseTimeUnit(s string) (TimeUnits, bool) {
	unit, found := timeUnitMap[strings.ToLower(strings.TrimSpace(s))]
	return unit, found
}

type TimeUnit struct {
	Unit TimeUnits
}

func NewTimeUnit(input string) (RecordUnit, error) {
	unit, found := parseTimeUnit(input)

	if !found {
		return nil, fmt.Errorf("Unknown unit %s", input)
	}

	return &TimeUnit{
		Unit: unit,
	}, nil
}

func (t *TimeUnit) getUnit() int {
	return int(t.Unit)
}

func (t *TimeUnit) ConvertToBaseUnit(value float64) float64 {
	// TODO implement conversion
	return 0
}

func (t *TimeUnit) ConvertFromBaseUnit(value float64) float64 {
	// TODO implement conversion
	return 0
}

func newRecordUnit(unit string) (RecordUnit, error) {
	if unit == "" {
		return nil, nil
	}

	new_functions := []func(input string) (RecordUnit, error){
		NewTimeUnit,
	}

	for _, f := range new_functions {
		recordUnit, err := f(unit)
		if err == nil {
			return recordUnit, nil
		}
	}

	return nil, fmt.Errorf("no matching unit found for %s", unit)
}

type Record struct {
	Label string
	Value float64
	Unit  RecordUnit
}

func splitInputString(input string) (label string, value float64, unit string, err error) {
	label, value_with_unit, found := strings.Cut(strings.TrimSpace(input), ":")

	if !found {
		return "", 0, "", fmt.Errorf("No colon found in input string")
	}

	if label == "" {
		return "", 0, "", fmt.Errorf("No label found in input string")
	}

	value_str, unit, _ := strings.Cut(strings.TrimSpace(value_with_unit), " ")

	value, err = strconv.ParseFloat(value_str, 64)

	if err != nil {
		return "", 0, "", err
	}

	return label, value, unit, nil
}

// NewRecord creates a new Record with the given label, value, and unit string
func NewRecord(label string, value float64, unit_str string) (*Record, error) {
	unit, err := newRecordUnit(unit_str)
	if err != nil {
		return nil, fmt.Errorf("failed to create record: %w", err)
	}
	return &Record{
		Label: label,
		Value: value,
		Unit:  unit,
	}, nil
}

func NewRecordFromString(input string) (*Record, error) {
	label, value, unit, err := splitInputString(input)

	if err != nil {
		return nil, err
	}

	record, err := NewRecord(label, value, unit)

	if err != nil {
		return nil, err
	}

	return record, nil
}

// ReadFileContent reads the content of a Markdown file and returns it as a byte slice.
package parsing

import (
	"fmt"
	"strings"
)

type ElementUnit interface {
	ConvertToBaseUnit(value float64) float64
	ConvertFromBaseUnit(value float64) float64
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

type TimeUnit struct {
	Unit  TimeUnits
	Value float64
}

func (t TimeUnit) ConvertToBaseUnit(value float64) float64 {
	return 0
}

func (t TimeUnit) ConvertFromBaseUnit(value float64) float64 {
	return 0
}

type Element struct {
	Label string
	Value float64
	Unit  ElementUnit
}

func NewFromString(input string) (Element, error) {
	before, after, found := strings.Cut(input, ":")

	if !found {
		return Element{}, fmt.Errorf("No colon found in input string")
	}

	label := strings.TrimSpace(before)



}

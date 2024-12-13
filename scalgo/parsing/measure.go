package parsing

import (
	"fmt"
	"strconv"
	"strings"
)

type Measure struct {
	Value float64
	Unit  Unit
}

func (m *Measure) getBaseValue() float64 {
	return m.Value * m.Unit.multiplier
}

func (m *Measure) isEmpty() bool {
	return m.Value == 0 && m.Unit.isEmpty()
}

func (m *Measure) Scale() Measure {
	return Measure{
		Value: m.Value,
		Unit:  m.Unit,
	}
}

func newMeasure(unit_files *UnitFiles, measure string) (Measure, error) {

	value_s, label, found := strings.Cut(strings.TrimSpace(measure), " ")

	if !found {
		return Measure{}, fmt.Errorf("failed to split %s by space", measure)
	}

	_, unit, found := unit_files.findUnit(strings.TrimSpace(label))

	if !found {
		return Measure{}, fmt.Errorf("Unit %s was not found", measure)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(value_s), 64)

	if err != nil {
		return Measure{}, fmt.Errorf("Failed to convert %s to f64", measure)
	}

	return Measure{
		Value: value,
		Unit:  unit,
	}, nil
}

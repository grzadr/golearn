package parsing

import (
	"fmt"
	"math"
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

func (m *Measure) scale(ref *Measure, scale *Measure) Measure {
	return Measure{
		Value: (m.getBaseValue() / ref.getBaseValue()) * scale.getBaseValue(),
		Unit:  Unit{},
	}
}

func (m *Measure) str(units *UnitSlice, max_units int) string {
	result := ""
	// prev_unit := Unit{}
	used_units := 0
	leftover := m.getBaseValue()
	for _, unit := range *units {
		if leftover == 0.0 {
			break
		}
		// if leftover >= m.getBaseValue() {
		// 	prev_unit = unit
		// 	continue
		// }

		part := math.Floor(leftover / unit.multiplier)
		if part <= 0.0 {
			continue
		}

		used_units++

		if used_units == max_units {
			result += fmt.Sprintf(", %.02f %s", leftover/unit.multiplier, unit.Name)
			break
		}

		leftover = leftover - (part * unit.multiplier)

		result += fmt.Sprintf(", %d %s", int(part), unit.Name)
	}

	return result[2:]
}

func splitMeasureString(measure string) (string, string, error) {
	value, label, found := strings.Cut(strings.TrimSpace(measure), " ")

	if !found {
		return "", "", fmt.Errorf("failed to split %s by space", measure)
	}

	return value, label, nil
}

func convertUnitValueString(str string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(str), 64)

	if err != nil {
		return value, fmt.Errorf("Failed to convert %s to f64", str)
	}

	return value, nil
}

func newMeasure(unit_files *UnitFiles, measure_str string) (
	measure Measure, unit_file string, err error) {

	value_str, label, err := splitMeasureString(measure_str)

	if err != nil {
		return measure, unit_file, err
	}

	unit_file, unit, found := unit_files.findUnit(strings.TrimSpace(label))

	if !found {
		return measure, unit_file, fmt.Errorf("Unit %s was not found", measure_str)
	}

	value, err := convertUnitValueString(value_str)

	if err != nil {
		return measure, unit_file, err
	}

	measure = Measure{
		Value: value,
		Unit:  unit,
	}

	return measure, unit_file, nil
}

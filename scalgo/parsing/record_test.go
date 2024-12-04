package parsing

import (
	"errors"
	"fmt"
	"testing"
)

func helperCompareRecords(ref *Record, other *Record) error {
	detected := make([]error, 0, 16)

	if ref.Label != other.Label {
		detected = append(detected, fmt.Errorf("Expected label %s, got %s", ref.Label, other.Label))
	}

	if ref.Measure.Value != other.Measure.Value {
		detected = append(detected, fmt.Errorf("Expected Value %f, got %f", ref.Measure.Value, other.Measure.Value))
	}

	if ref.Measure.Unit.Name != other.Measure.Unit.Name {
		detected = append(detected, fmt.Errorf("Expected name %s, got %s", ref.Measure.Unit.Name, other.Measure.Unit.Name))
	}

	if ref.Measure.Unit.multiplier != other.Measure.Unit.multiplier {
		detected = append(detected, fmt.Errorf("Expected multiplier %f, got %f", ref.Measure.Unit.multiplier, other.Measure.Unit.multiplier))
	}

	if len(detected) > 0 {
		return errors.Join(detected...)
	}

	return nil
}

func TestNewRecord(t *testing.T) {
	query := "Some label: 1 minute"

	expected := Record{
		Label: "Some label",
		Measure: Measure{
			Value: 1,
			Unit: Unit{
				Name:       "minute",
				multiplier: 60,
			},
		},
	}

	record, err := newRecord(query, embedded_units)

	if err != nil {
		t.Errorf("newRecord returned an error: %v", err)
	}

	if err := helperCompareRecords(&expected, &record); err != nil {
		t.Errorf("Found differences from expected: %v", err)
	}
}

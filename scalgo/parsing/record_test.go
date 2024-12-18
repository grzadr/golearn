package parsing

import (
	"fmt"
	"testing"
)

func helperCompareRecords(ref *Record, other *Record) []error {
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

	return detected
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

	expected_unit_file := "time"

	record, unit_file, err := newRecord(query, EmbeddedUnits)

	if err != nil {
		t.Errorf("newRecord returned an error: %v", err)
	}

	registerErrors(
		helperCompareRecords(&expected, &record),
		"Detected difference between Records",
		t,
	)

	if unit_file != expected_unit_file {
		t.Errorf(
			"Expected unit_file `%s`, got %s",
			expected_unit_file,
			unit_file,
		)
	}
}

func TestNewRecordError(t *testing.T) {
	cases := []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "Missing : ",
			input: "Label 1 hour",
			err:   "Record `Label 1 hour` missing `: `",
		},
		{
			name:  "Missing Label",
			input: ": 1 hour",
			err:   "Record `: 1 hour` missing label",
		},
		{
			name:  "Missing Label",
			input: "Label: hour",
			err:   "failed to split hour by space",
		},
	}

	for _, c := range cases {
		record, unit_file, err := newRecord(c.input, EmbeddedUnits)

		if !record.isEmpty() {
			t.Errorf("Test %s: Expected nil, got %v", c.name, record)
		}

		if err == nil || err.Error() != c.err {
			t.Errorf("Test %s: Expected error %s, got %v", c.name, c.err, err)
		}

		if unit_file != "" {
			t.Errorf("Expected empty unit_file, got %s", unit_file)
		}

	}
}

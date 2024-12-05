package parsing

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	// "testing/fstest"
)

const EnlistmentTestVarBasicString string = `@scale 1 year
Item 2: 50 minutes

Item 3: 10 minutes
Item 1: 100 minutes

`

var EnlistmentTestBasicRecordSlice []Record = []Record{
	{
		Label: "Item 1",
		Measure: Measure{
			Value: 100,
			Unit: Unit{
				Name:       "minute",
				multiplier: 60,
			},
		},
	},
	{
		Label: "Item 2",
		Measure: Measure{
			Value: 10,
			Unit: Unit{
				Name:       "minute",
				multiplier: 60,
			},
		},
	},
	{
		Label: "Item 3",
		Measure: Measure{
			Value: 50,
			Unit: Unit{
				Name:       "minute",
				multiplier: 60,
			},
		},
	},
}

var EnlistmentTestVarBasicObj Enlistment = Enlistment{
	options: Options{
		sort:     true,
		reversed: false,
		scale: &Measure{
			Value: 1,
			Unit: Unit{
				Name:       "year",
				multiplier: 31536000,
			},
		},
	},
	records: EnlistmentTestBasicRecordSlice,
	ref:     &EnlistmentTestBasicRecordSlice[0],
}

func helperCompareRecordSlice(ref *[]Record, sub *[]Record) []error {
	var errs []error

	// Guard against nil pointers
	if ref == nil || sub == nil {
		return append(errs, fmt.Errorf("nil slice pointer provided"))
	}

	// Compare lengths
	if len(*ref) != len(*sub) {
		errs = append(errs, fmt.Errorf("slices have different lengths: ref=%d, sub=%d",
			len(*ref), len(*sub)))
		// Return early as order comparison wouldn't make sense
		return errs
	}

	for i := range len(*ref) {
		fmt.Printf("Record %d", i)
		errs = append(errs, helperCompareRecords(&(*ref)[i], &(*sub)[i]))
	}

	return errs

}

func helperCompareOptions(ref Options, sub Options) []error {
	var errs []error

	// Compare basic fields
	if ref.sort != sub.sort {
		errs = append(errs, fmt.Errorf("sort field mismatch: expected %v, got %v",
			ref.sort, sub.sort))
	}

	if ref.reversed != sub.reversed {
		errs = append(errs, fmt.Errorf("reversed field mismatch: expected %v, got %v",
			ref.reversed, sub.reversed))
	}

	// Compare scale field and its nested structures
	if ref.scale == nil && sub.scale != nil {
		errs = append(errs, fmt.Errorf("scale mismatch: expected nil, got non-nil"))
	} else if ref.scale != nil && sub.scale == nil {
		errs = append(errs, fmt.Errorf("scale mismatch: expected non-nil, got nil"))
	} else if ref.scale != nil && sub.scale != nil {
		if ref.scale.Value != sub.scale.Value {
			errs = append(errs, fmt.Errorf("scale value mismatch: expected %v, got %v",
				ref.scale.Value, sub.scale.Value))
		}

		if ref.scale.Unit.Name != sub.scale.Unit.Name {
			errs = append(errs, fmt.Errorf("scale unit name mismatch: expected %q, got %q",
				ref.scale.Unit.Name, sub.scale.Unit.Name))
		}

		if ref.scale.Unit.multiplier != sub.scale.Unit.multiplier {
			errs = append(errs, fmt.Errorf("scale unit multiplier mismatch: expected %v, got %v",
				ref.scale.Unit.multiplier, sub.scale.Unit.multiplier))
		}
	}

	return errs
}

func helperCompareEnlistment(ref *Enlistment, sub *Enlistment) error {
	errs := make([]error, 0, 16)
	if sub == nil {
		return fmt.Errorf("Expected enlistment to not be nil")
	}

	option_errs := helperCompareOptions(ref.options, sub.options)

	if len(option_errs) > 0 {
		errs = append(errs, option_errs...)
	}

	records_errs := helperCompareRecordSlice(&ref.records, &sub.records)

	if len(records_errs) > 0 {
		errs = append(errs, records_errs...)
	}

	if ref_err := helperCompareRecords(ref.ref, sub.ref); ref_err != nil {
		errs = append(errs, ref_err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func TestNewRecordEnlistmentFromReader(t *testing.T) {
	// Given
	reader := strings.NewReader(EnlistmentTestVarBasicString)

	// When
	enlistment, err := NewRecordEnlistmentFromReader(reader, embedded_units)

	// Then
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := helperCompareEnlistment(&EnlistmentTestVarBasicObj, enlistment); err != nil {
		t.Errorf("Enlistment differ from reference:\n%v", err)
	}

	// Add more specific assertions based on your expected output
}

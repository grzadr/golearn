package parsing

import (
	"fmt"
	"strings"
	"testing"
	// "testing/fstest"
)

const EnlistmentTestVarBasicString string = `@scale 1 year
Item 2: 15 minutes

Item 3: 60 seconds
Item 1: 1 hour

`

var EnlistmentTestBasicRecordSlice []Record = []Record{
	{
		Label: "Item 1",
		Measure: Measure{
			Value: 1,
			Unit: Unit{
				Name:       "hour",
				multiplier: 3600,
			},
		},
	},
	{
		Label: "Item 2",
		Measure: Measure{
			Value: 15,
			Unit: Unit{
				Name:       "minute",
				multiplier: 60,
			},
		},
	},
	{
		Label: "Item 3",
		Measure: Measure{
			Value: 60,
			Unit: Unit{
				Name:       "second",
				multiplier: 1,
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

	// Guard against nil pointers
	if ref == nil || sub == nil {
		return []error{fmt.Errorf("nil slice pointer provided")}
	}

	// Compare lengths
	if len(*ref) != len(*sub) {
		return []error{fmt.Errorf("slices have different lengths: ref=%d, sub=%d",
			len(*ref), len(*sub))}
	}

	errs := make([]error, 0, len(*ref)*6)

	for i := range len(*ref) {
		fmt.Printf("Record %d", i)
		errs = append(errs, helperCompareRecords(&(*ref)[i], &(*sub)[i])...)
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

func helperCompareEnlistment(ref *Enlistment, sub *Enlistment) []error {

	if sub == nil {
		return []error{fmt.Errorf("Expected enlistment to not be nil")}
	}

	errs := make([]error, 0, 16)

	errs = append(errs, helperCompareOptions(ref.options, sub.options)...)
	errs = append(errs, helperCompareRecordSlice(&ref.records, &sub.records)...)
	errs = append(errs, helperCompareRecords(ref.ref, sub.ref)...)

	return errs
}

func TestNewRecordEnlistmentFromReader(t *testing.T) {
	reader := strings.NewReader(EnlistmentTestVarBasicString)

	enlistment, err := NewRecordEnlistmentFromReader(reader, embedded_units)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	registerErrors(
		helperCompareEnlistment(&EnlistmentTestVarBasicObj, enlistment),
		"Enlistment differ from reference", t)
}

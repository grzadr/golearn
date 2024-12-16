package parsing

import (
	"fmt"
	"strings"
	"testing"
	"testing/fstest"
)

func TestIsTrue(t *testing.T) {
	test_cases := []struct {
		has_value bool
		value     string
		expected  bool
	}{
		{
			has_value: false,
			value:     "true",
			expected:  true,
		},
		{
			has_value: true,
			value:     "true",
			expected:  true,
		},
		{
			has_value: false,
			value:     "false",
			expected:  true,
		},
		{
			has_value: true,
			value:     "false",
			expected:  false,
		},
		{
			has_value: true,
			value:     "",
			expected:  false,
		},
	}

	for _, tc := range test_cases {
		if res := isTrue(tc.has_value, tc.value); res != tc.expected {
			t.Errorf("Result %t differs from %+v", res, tc)
		}
	}
}

func TestIsFalse(t *testing.T) {
	test_cases := []struct {
		value    string
		expected bool
	}{
		{
			value:    "false",
			expected: true,
		},
		{
			value:    "true",
			expected: false,
		},
		{
			value:    "",
			expected: false,
		},
	}

	for _, tc := range test_cases {
		if res := isFalse(tc.value); res != tc.expected {
			t.Errorf("Result %t differs from %+v", res, tc)
		}
	}
}

func TestSetFlag(t *testing.T) {
	test_cases := []struct {
		has_value bool
		value     string
		flag      bool
		expected  bool
		error_msg string
	}{
		{
			has_value: false,
			value:     "true",
			flag:      false,
			expected:  true,
		},
		{
			has_value: true,
			value:     "true",
			flag:      false,
			expected:  true,
		},
		{
			has_value: false,
			value:     "false",
			flag:      false,
			expected:  true,
		},
		{
			has_value: true,
			value:     "false",
			flag:      true,
			expected:  false,
		},
		{
			has_value: true,
			value:     "",
			flag:      true,
			expected:  false,
			error_msg: "Unknown value: ",
		},
	}

	for _, tc := range test_cases {
		flag := &tc.flag
		if len(tc.error_msg) > 0 {
			if err := setFlag(flag, tc.has_value, tc.value); err == nil || err.Error() != tc.error_msg {
				t.Errorf("Expected error %s, got %v", tc.error_msg, err)
			}
		} else {
			err := setFlag(flag, tc.has_value, tc.value)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			if *flag != tc.expected {
				t.Errorf("Expected %+v, got %t", tc, *flag)
			}
		}
	}
}

const EnlistmentTestVarBasicString string = `@scale 1 year
Item 2: 15 minutes

Item 3: 60 seconds
Item 1: 1 hour

`

const EnlistmentTestVarReversedString string = `@scale 1 year
Item 2: 15 minutes
@reverse
Item 3: 60 seconds
Item 1: 1 hour

`

const EnlistmentTestVarUnsortedString string = `@scale 1 year
Item 2: 15 minutes
@sort false
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

var EnlistmentTestReverseRecordSlice []Record = []Record{
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
		Label: "Item 1",
		Measure: Measure{
			Value: 1,
			Unit: Unit{
				Name:       "hour",
				multiplier: 3600,
			},
		},
	},
}

var EnlistmentTestUnsortedRecordSlice []Record = []Record{
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
}

var EnlistmentTestScale Measure = Measure{
	Value: 1,
	Unit: Unit{
		Name:       "year",
		multiplier: 31536000,
	},
}

var EnlistmentTestVarBasicObj Enlistment = Enlistment{
	options: EnlistmentOptions{
		sort:     true,
		reversed: false,
		scale:    &EnlistmentTestScale,
	},
	records:   EnlistmentTestBasicRecordSlice,
	ref:       &EnlistmentTestBasicRecordSlice[0],
	unit_file: "time",
}

var EnlistmentTestVarReverseObj Enlistment = Enlistment{
	options: EnlistmentOptions{
		sort:     true,
		reversed: true,
		scale:    &EnlistmentTestScale,
	},
	records: EnlistmentTestReverseRecordSlice,
	ref:     &EnlistmentTestReverseRecordSlice[0],
}

var EnlistmentTestVarUnsortedObj Enlistment = Enlistment{
	options: EnlistmentOptions{
		sort:     false,
		reversed: false,
		scale:    &EnlistmentTestScale,
	},
	records: EnlistmentTestUnsortedRecordSlice,
	ref:     &EnlistmentTestUnsortedRecordSlice[2],
}

func helperCompareRecordSlice(ref *RecordSlice, sub *RecordSlice) []error {

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

	for i, p := range IterZip(*ref, *sub) {
		fmt.Printf("Record %d\n", i)
		errs = append(errs, helperCompareRecords(&p.First, &p.Second)...)
	}

	return errs

}

func helperCompareOptions(ref EnlistmentOptions, sub EnlistmentOptions) []error {
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

	enlistment, err := NewEnlistmentFromReader(reader, embedded_units)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	registerErrors(
		helperCompareEnlistment(&EnlistmentTestVarBasicObj, enlistment),
		"Enlistment differ from reference", t)
}

func TestNewRecordEnlistmentReverse(t *testing.T) {
	reader := strings.NewReader(EnlistmentTestVarReversedString)

	enlistment, err := NewEnlistmentFromReader(reader, embedded_units)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	registerErrors(
		helperCompareEnlistment(&EnlistmentTestVarReverseObj, enlistment),
		"Enlistment differ from reference", t)
}

func TestNewRecordEnlistmentUnsorted(t *testing.T) {
	reader := strings.NewReader(EnlistmentTestVarUnsortedString)

	enlistment, err := NewEnlistmentFromReader(reader, embedded_units)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	registerErrors(
		helperCompareEnlistment(&EnlistmentTestVarUnsortedObj, enlistment),
		"Enlistment differ from reference", t)
}

var enlistmentTestFS = fstest.MapFS{
	"test_enlistment.txt": {
		Data: []byte(EnlistmentTestVarBasicString),
	},
}

func TestNewRecordEnlistmentFromFile(t *testing.T) {
	enlistment, err := NewEnlistmentFromFile(
		enlistmentTestFS,
		"test_enlistment.txt",
		embedded_units,
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	registerErrors(
		helperCompareEnlistment(&EnlistmentTestVarBasicObj, enlistment),
		"Enlistment differ from reference", t)
}

var FixtureScaledRecordSlice = RecordSlice{
	Record{
		Label: "Item 1",
		Measure: Measure{
			Value: 31536000,
			Unit: Unit{
				Name:       "second",
				multiplier: 1,
			},
		},
	},
	Record{
		Label: "Item 2",
		Measure: Measure{
			Value: 7884000,
			Unit: Unit{
				Name:       "second",
				multiplier: 1,
			},
		},
	},
	Record{
		Label: "Item 3",
		Measure: Measure{
			Value: 525600,
			Unit: Unit{
				Name:       "second",
				multiplier: 1,
			},
		},
	},
}

func TestEnlistmentScaleRecords(t *testing.T) {
	result, err := EnlistmentTestVarBasicObj.scaleRecords(embedded_units)

	if err != nil {
		t.Error(err)
		return
	}

	registerErrors(
		helperCompareRecordSlice(&FixtureScaledRecordSlice, &result),
		"Enlistment.scaleRecords failed tests",
		t,
	)
}

var FixtureExpectedRecordSliceStr = []string{
	"Item 1: 1 year",
	"Item 2: 3 months",
	"Item 3: 6 day, 2 hour",
}

func TestRecordSliceStr(t *testing.T) {
	FixtureScaledRecordSlice.str()
}

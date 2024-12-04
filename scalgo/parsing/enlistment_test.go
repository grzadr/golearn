package parsing

import (
	"strings"
	"testing"
	// "testing/fstest"
)

const EnlistmentTestVarBasicString string = `@scale 1 year
Item 2: 10
Item 1: 100 years

Item 3: 50
`

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
	records: []Record{
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
	},
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
	if enlistment == nil {
		t.Error("expected enlistment to not be nil")
	}

	// Add more specific assertions based on your expected output
}

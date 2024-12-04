package parsing


import (
	"testing"
	"strings"
	"testing/fstest"
)

const EnlistmentTestVarBasicString string = `@scale 1 year
Item 1: 100 years
Item 2: 10
Item 3: 50
`

const EnlistmentTestVarBasicObj Enlistment = Enlistment{
	options: Options {
		Sort: true,
		Reversed: false,
		Scale: Measure {
			Value: 1,
			Unit: {
				Name: "year",
				value: 31536000,
			},
		},
	},
	Records 
}

func TestNewRecordEnlistmentFromReader(t *testing.T) {
    // Given
    content :=
    reader := strings.NewReader(content)

    // When
    enlistment, err := NewRecordEnlistmentFromReader(reader)

    // Then
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if enlistment == nil {
        t.Error("expected enlistment to not be nil")
    }

    // Add more specific assertions based on your expected output
}

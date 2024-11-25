package parsing

import (
	"testing"
)

var testUnitEntryJsonData = []byte(`{
    "meter": {
        "value": 1.0,
        "aliases": ["m", "meters"]
    },
    "kilometer": {
        "value": 1000.0,
        "aliases": ["km", "kilometers"]
    }
}`)

var testExpectedData = []UnitEntry{
	{ // Removed UnitEntry{} syntax as it's redundant
		Name:    "meter",
		Value:   1.0,
		Aliases: []string{"m", "meters"},
	},
	{
		Name:    "kilometer",
		Value:   1000.0,
		Aliases: []string{"km", "kilometers"},
	},
}

func helpCompareUnitEntry(a, b *UnitEntry) bool {
	if a.Name != b.Name || a.Value != b.Value || len(a.Aliases) != len(b.Aliases) {
		return false
	}

	for i, alias := range a.Aliases {
		if alias != b.Aliases[i] {
			return false
		}
	}

	return true
}

func TestIterUnitEntries_Success(t *testing.T) {
	for i, next := range IterUnitEntries(testUnitEntryJsonData) {
		if next.Err != nil {
			t.Fatalf("Received unexpected error %v", next.Err)
		}

		if !helpCompareUnitEntry(&next.Entry, &testExpectedData[i]) {
			t.Errorf("Entry %d: expected %+v, got %+v", i, testExpectedData[i], next.Entry)
		}
	}
}

func TestIterUnitEntries_EarlyTermination(t *testing.T) {
	count := 0
	for i, next := range IterUnitEntries(testUnitEntryJsonData) {
		count = i
		if next.Err != nil {
			t.Errorf("Received unexpected error %v", next.Err)
		}
		break
	}

	if count != 0 {
		t.Errorf("Expected to terminate at 0, but terminated at %d", count)
	}
}

func TestIterUnitEntries_InvalidJSON(t *testing.T) {
	testCases := []struct {
		name    string
		input   []byte
		wantErr string
	}{
		{
			name:  "invalid opening delimiter",
			input: []byte(`["not an object"]`),
			// Exact match for the error message
			wantErr: "expected {, got [",
		},
		{
			name:    "corrupted JSON",
			input:   []byte(`{invalid`),
			wantErr: "reading key: invalid character 'i'",
		},
		{
			name:    "non-string key",
			input:   []byte(`{123: {}}`),
			wantErr: "reading key: invalid character '1'",
		},
		{
			name:    "invalid value structure",
			input:   []byte(`{"key": "not an object"}`),
			wantErr: "decoding value for \"key\": json: cannot unmarshal string into Go value of type parsing.UnitEntry",
		},
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: "reading token: EOF",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotErr error
			for _, next := range IterUnitEntries(tc.input) {
				gotErr = next.Err
				break // We only need to check the first error
			}

			if gotErr == nil {
				t.Fatal("expected error, got nil")
			}
			if gotErr.Error() != tc.wantErr { // Changed to exact match
				t.Errorf("expected error %q, got %q", tc.wantErr, gotErr.Error())
			}
		})
	}
}

func TestIterUnitEntries_ComplexErrors(t *testing.T) {
	testCases := []struct {
		name    string
		input   []byte
		wantErr string
	}{
		{
			name: "missing value field",
			input: []byte(`{
                "meter": {
                    "aliases": ["m"]
                }
            }`),
			wantErr: "invalid entry \"meter\": value field is required",
		},
		{
			name: "missing aliases field",
			input: []byte(`{
                "meter": {
                    "value": 1.0
                }
            }`),
			wantErr: "invalid entry \"meter\": aliases field is required",
		},
		{
			name: "invalid value type",
			input: []byte(`{
                "meter": {
                    "value": "not a number",
                    "aliases": ["m"]
                }
            }`),
			wantErr: "decoding value for \"meter\": json: cannot unmarshal string into Go struct field UnitEntry.value of type float64",
		},
		{
			name: "invalid aliases type",
			input: []byte(`{
                "meter": {
                    "value": 1.0,
                    "aliases": "not an array"
                }
            }`),
			wantErr: "decoding value for \"meter\": json: cannot unmarshal string into Go struct field UnitEntry.aliases of type []string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotErr error
			for _, next := range IterUnitEntries(tc.input) {
				gotErr = next.Err
				break
			}

			if gotErr == nil {
				t.Fatal("expected error, got nil")
			}
			if gotErr.Error() != tc.wantErr {
				t.Errorf("expected error %q, got %q", tc.wantErr, gotErr.Error())
			}
		})
	}
}

func TestUnitEntry_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		entry   UnitEntry
		wantErr string
	}{
		{
			name: "valid entry",
			entry: UnitEntry{
				Name:    "meter",
				Value:   1.0,
				Aliases: []string{"m"},
			},
			wantErr: "",
		},
		{
			name: "missing value",
			entry: UnitEntry{
				Name:    "meter",
				Aliases: []string{"m"},
			},
			wantErr: "value field is required",
		},
		{
			name: "missing aliases",
			entry: UnitEntry{
				Name:  "meter",
				Value: 1.0,
			},
			wantErr: "aliases field is required",
		},
		{
			name:    "zero value struct",
			entry:   UnitEntry{},
			wantErr: "value field is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.entry.Validate()
			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tc.wantErr {
				t.Errorf("expected error %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}


// Benchmark to ensure performance
func BenchmarkIterUnitEntries(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, next := range IterUnitEntries(testUnitEntryJsonData) {
			if next.Err != nil {
				b.Fatal(next.Err)
			}
		}
	}
}

// package parsing

// import (
// 	"testing"
// )

// var testUnitEntryJsonData = []byte(`{
//     "meter": {
//         "value": 1.0,
//         "aliases": ["m", "meters"]
//     },
//     "kilometer": {
//         "value": 1000.0,
//         "aliases": ["km", "kilometers"]
//     }
// }`)

// var testExpectedData = []UnitEntry{
// 	UnitEntry{
// 		Name:    "meter",
// 		Value:   1.0,
// 		Aliases: []string{"m", "meters"},
// 	},
// 	UnitEntry{
// 		Name:    "kilometer",
// 		Value:   1000.0,
// 		Aliases: []string{"km", "kilometers"},
// 	},
// }

// func helpCompareUnitEntry(a *UnitEntry, b *UnitEntry) bool {
// 	if a.Name != b.Name || a.Value != b.Value || len(a.Aliases) != len(b.Aliases) {
// 		return false
// 	}

// 	a_aliases := a.Aliases
// 	b_aliases := b.Aliases

// 	for i := 0; i < len(a_aliases); i++ {
// 		if a_aliases[i] != b_aliases[i] {
// 			return false
// 		}
// 	}

// 	return true
// }

// func TestIterUnitEntry(t *testing.T) {
// 	for i, next := range IterUnitEntries(testUnitEntryJsonData) {
// 		if next.Err != nil {
// 			t.Errorf("Received unexpected error %v", next.Err)
// 		}

// 		if !helpCompareUnitEntry(&next.Entry, &testExpectedData[i]) {
// 			t.Errorf("Expected entry %v, got %v", next.Entry, testExpectedData[i])
// 			return
// 		}
// 	}
// }

// func TestIterUnitEntry_EarlyTermination(t *testing.T) {
// 	last_i := 0
// 	for i, next := range IterUnitEntries(testUnitEntryJsonData) {
// 		last_i = i
// 		if next.Err != nil {
// 			t.Errorf("Received unexpected error %v", next.Err)
// 		}

// 		if !helpCompareUnitEntry(&next.Entry, &testExpectedData[i]) {
// 			t.Errorf("Expected entry %v, got %v", next.Entry, testExpectedData[i])
// 			return
// 		}

// 		break
// 	}

// 	if last_i > 0 {
// 		t.Errorf("Expected to terminate at 0, but terminated at %d", last_i)
// 	}
// }

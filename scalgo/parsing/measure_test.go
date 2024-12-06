package parsing

import (
	"strings"
	"testing"
)

var TestUnitFiles = func() UnitFiles {
	files, err := newUnitFiles(testFS, "units")

	if err != nil {
		panic(err)
	}

	return files
}()

func TestNewMeasureInvalid(t *testing.T) {
	// Using table-driven tests for better organization and coverage
	tests := []struct {
		name        string
		input       string
		wantValue   float64
		wantUnit    string
		wantBaseVal float64
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid kilometer measurement",
			input:       "42 kilometers",
			wantValue:   42.0,
			wantUnit:    "kilometer",
			wantBaseVal: 1000.0,
			wantErr:     false,
		},
		{
			name:        "missing space between value and unit",
			input:       "42kilometers",
			wantErr:     true,
			errContains: "failed to split",
		},
		{
			name:        "unknown unit",
			input:       "42 lightyears",
			wantErr:     true,
			errContains: "was not found",
		},
		{
			name:        "invalid numeric value",
			input:       "4x2 kilometers",
			wantErr:     true,
			errContains: "Failed to convert",
		},

		{
			name:        "extra whitespace handling",
			input:       "  42   kilometers  ",
			wantValue:   42.0,
			wantUnit:    "kilometer",
			wantBaseVal: 1000.0,
			wantErr:     false,
		},
		{
			name:        "zero value",
			input:       "0 kilometers",
			wantValue:   0.0,
			wantUnit:    "kilometer",
			wantBaseVal: 1000.0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			measure, err := newMeasure(&TestUnitFiles, tt.input)

			// Error case handling
			if tt.wantErr {
				if err == nil {
					t.Errorf("newMeasure() expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("newMeasure() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			// Success case validation
			if err != nil {
				t.Errorf("newMeasure() unexpected error: %v", err)
				return
			}

			if measure.Value != tt.wantValue {
				t.Errorf("Value = %v, want %v", measure.Value, tt.wantValue)
			}

			if measure.Unit.Name != tt.wantUnit {
				t.Errorf("Unit.Name = %v, want %v", measure.Unit.Name, tt.wantUnit)
			}

			if measure.Unit.multiplier != tt.wantBaseVal {
				t.Errorf("Unit.value = %v, want %v", measure.Unit.multiplier, tt.wantBaseVal)
			}
		})
	}
}

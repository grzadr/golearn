package parsing

import (
	"testing"
)

func TestParseTimeUnitValid(t *testing.T) {
	unit, found := parseTimeUnit("nanosecond")

	if found != true {
		t.Fatalf("Unexpected error: %t", found)
	}

	if unit != Nanosecond {
		t.Errorf("Expected unit to be %d, got %d", Nanosecond, unit)
	}
}

func TestParseTimeUnitInvalid(t *testing.T) {
	unit, found := parseTimeUnit("invalid")

	if found != false {
		t.Fatalf("Unexpected error: %t", found)
	}

	if unit != 0 {
		t.Errorf("Expected unit to be %d, got %d", 0, unit)
	}
}

func TestNewTimeUnitValid(t *testing.T) {
	unit, err := NewTimeUnit("nanosecond")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if TimeUnits(unit.getUnit()) != Nanosecond {
		t.Errorf("Expected unit to be %d, got %d", Nanosecond, unit.getUnit())
	}
}

func TestNewTimeUnitInvalid(t *testing.T) {
	_, err := NewTimeUnit("invalid")

	if err == nil {
		t.Error("Expected an error for invalid unit, but got nil")
	}

	if err.Error() != "Unknown unit invalid" {
		t.Errorf("Expected error message 'Unknown unit invalid', got '%s'", err.Error())
	}
}

func TestTimeUnitgetUnit(t *testing.T) {
	unit := TimeUnit{Unit: Nanosecond}

	if unit.getUnit() != int(Nanosecond) {
		t.Errorf("Expected unit to be %d, got %d", Nanosecond, unit.getUnit())
	}
}

func TestNewRecordUnitValid(t *testing.T) {
	unit, err := newRecordUnit("nanosecond")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if TimeUnits(unit.(*TimeUnit).getUnit()) != Nanosecond {
		t.Errorf("Expected unit to be %d, got %d", Nanosecond, unit.(*TimeUnit).getUnit())
	}
}

func TestNewRecordUnitEmpty(t *testing.T) {
	unit, err := newRecordUnit("")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if unit != nil {
		t.Errorf("Expected unit to be nil, got %v", unit)
	}
}

func TestNewRecordUnitInvalid(t *testing.T) {
	_, err := newRecordUnit("invalid")

	if err == nil {
		t.Error("Expected an error for invalid unit, but got nil")
	}

	if err.Error() != "no matching unit found for invalid" {
		t.Errorf("Expected error message 'no matching unit found for invalid', got '%s'", err.Error())
	}
}

func TestLoaderSplitInputString(t *testing.T) {
	input := "Label 1: 3.14 years"
	label, value, unit, err := splitInputString(input)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if label != "Label 1" {
		t.Errorf("Expected label to be '%s', got '%s'", "Label 1", label)
	}

	if value != 3.14 {
		t.Errorf("Expected value to be %f, got %f", 3.14, value)
	}

	if unit != "years" {
		t.Errorf("Expected unit to be '%s', got '%s'", "years", unit)
	}
}

func TestLoaderSplitInputStringNoColon(t *testing.T) {
	input := "Label 1 3.14 years"
	_, _, _, err := splitInputString(input)

	if err == nil {
		t.Error("Expected an error for no colon, but got nil")
	}

	if err.Error() != "No colon found in input string" {
		t.Errorf("Expected error message 'No colon found in input string', got '%s'", err.Error())
	}
}

func TestLoaderSplitInputStringNoLabel(t *testing.T) {
	input := ": 3.14 years"
	_, _, _, err := splitInputString(input)

	if err == nil {
		t.Error("Expected an error for no label, but got nil")
	}

	if err.Error() != "No label found in input string" {
		t.Errorf("Expected error message 'No label found in input string', got '%s'", err.Error())
	}
}

func TestLoaderSplitInputStringInvalidValue(t *testing.T) {
	input := "Label 1: not a number years"
	_, _, _, err := splitInputString(input)

	if err == nil {
		t.Error("Expected an error for invalid value, but got nil")
	}

	if err.Error() != "strconv.ParseFloat: parsing \"not\": invalid syntax" {
		t.Errorf("Expected error message 'strconv.ParseFloat: parsing \"not\": invalid syntax', got '%s'", err.Error())
	}
}

func TestNewRecord(t *testing.T) {
	label := "Label 1"
	value := 3.14
	unit_str := "years"
	record, err := NewRecord(label, value, unit_str)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if record.Label != label {
		t.Errorf("Expected label to be '%s', got '%s'", label, record.Label)
	}

	if record.Value != value {
		t.Errorf("Expected value to be %f, got %f", value, record.Value)
	}

	if record.Unit.(*TimeUnit).getUnit() != int(Year) {
		t.Errorf("Expected unit to be %d, got %d", Year, record.Unit.(*TimeUnit).getUnit())
	}
}

func TestNewRecordInvalidUnit(t *testing.T) {
	_, err := NewRecord("Label 1", 3.14, "invalid")

	if err == nil {
		t.Error("Expected an error for invalid unit, but got nil")
	}

	if err.Error() != "failed to create record: no matching unit found for invalid" {
		t.Errorf("Expected error message 'no matching unit found for invalid', got '%s'", err.Error())
	}
}

func TestNewRecordFromString(t *testing.T) {
	input := "Label 1: 3.14 years"
	record, err := NewRecordFromString(input)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if record.Label != "Label 1" {
		t.Errorf("Expected label to be '%s', got '%s'", "Label 1", record.Label)
	}

	if record.Value != 3.14 {
		t.Errorf("Expected value to be %f, got %f", 3.14, record.Value)
	}

	if record.Unit.(*TimeUnit).getUnit() != int(Year) {
		t.Errorf("Expected unit to be %d, got %d", Year, record.Unit.(*TimeUnit).getUnit())
	}
}

func TestNewRecordFromStringInvalid(t *testing.T) {
	input := "Label 1 3.14 years"
	_, err := NewRecordFromString(input)

	if err == nil {
		t.Error("Expected an error for no colon, but got nil")
	}

	if err.Error() != "No colon found in input string" {
		t.Errorf("Expected error message 'No colon found in input string', got '%s'", err.Error())
	}
}

func TestNewRecordFromStringInvalidValue(t *testing.T) {
	input := "Label 1: 3.14 invalid"
	_, err := NewRecordFromString(input)

	if err == nil {
		t.Error("Expected an error for invalid value, but got nil")
	}

	if err.Error() != "failed to create record: no matching unit found for invalid" {
		t.Errorf("Expected error message 'failed to create record: no matching unit found for invalid', got '%s'", err.Error())
	}
}

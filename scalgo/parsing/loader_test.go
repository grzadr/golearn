package parsing

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestParseTimeUnitValid(t *testing.T) {
	unit, found := parseTimeUnit("second")

	if found != true {
		t.Fatalf("Unexpected error: %t", found)
	}

	if unit != Second {
		t.Errorf("Expected unit to be %d, got %d", Second, unit)
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
	unit, err := NewTimeUnit("second")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if TimeUnits(unit.getUnit()) != Second {
		t.Errorf("Expected unit to be %d, got %d", Second, unit.getUnit())
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
	unit := TimeUnit{Unit: Second}

	if unit.getUnit() != int64(Second) {
		t.Errorf("Expected unit to be %d, got %d", Second, unit.getUnit())
	}
}

func TestTimeUnitConvertToUnit(t *testing.T) {
	unit := TimeUnit{Unit: Minute}
	result := unit.ConvertToUnit(1.0, int64(Second))
	expected := 60.0

	if result != expected {
		t.Errorf("Expected value to be %f, got %f", expected, result)
	}
}

func TestTimeUnitConvertToBaseUnit(t *testing.T) {
	unit := TimeUnit{Unit: Minute}

	if unit.ConvertToBaseUnit(1.0) != 60.0 {
		t.Errorf("Expected value to be 60.0, got %f", unit.ConvertToBaseUnit(1.0))
	}
}

func TestNewRecordUnitValid(t *testing.T) {
	unit, err := newRecordUnit("second")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if TimeUnits(unit.(*TimeUnit).getUnit()) != Second {
		t.Errorf("Expected unit to be %d, got %d", Second, unit.(*TimeUnit).getUnit())
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

	if record.Unit.getUnit() != int64(Year) {
		t.Errorf("Expected unit to be %d, got %d", Year, record.Unit.(*TimeUnit).getUnit())
	}

	if record.BaseValue != 3.14*float64(Year) {
		t.Errorf("Expected base value to be %f, got %f", 3.14*float64(Year), record.BaseValue)
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

	if record.Unit.(*TimeUnit).getUnit() != int64(Year) {
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

func TestNewRecordSliceFromReader(t *testing.T) {
	input := `Label 1: 3.14 years
Label 2: 42 days
Label 3: 1.5 hours`
	reader := strings.NewReader(input)

	enlistment, err := NewRecordEnlistmentFromReader(reader)

	records := enlistment.Records

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	// Check the first record
	if records[0].Label != "Label 1" || records[0].Value != 3.14 || records[0].Unit.getUnit() != int64(Year) {
		t.Errorf("First record doesn't match expected values")
	}

	// Check the second record
	if records[1].Label != "Label 2" || records[1].Value != 42 || records[1].Unit.getUnit() != int64(Day) {
		t.Errorf("Second record doesn't match expected values")
	}

	// Check the third record
	if records[2].Label != "Label 3" || records[2].Value != 1.5 || records[2].Unit.getUnit() != int64(Hour) {
		t.Errorf("Third record doesn't match expected values")
	}
}

func TestNewRecordSliceFromReaderEmptyInput(t *testing.T) {
	reader := strings.NewReader("")

	enlistment, err := NewRecordEnlistmentFromReader(reader)

	records := enlistment.Records

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(records) != 0 {
		t.Errorf("Expected 0 records, got %d", len(records))
	}
}

func TestNewRecordSliceFromReaderInvalidInput(t *testing.T) {
	input := `Label 1: 3.14 years
Invalid line
Label 3: 1.5 hours`
	reader := strings.NewReader(input)

	_, err := NewRecordEnlistmentFromReader(reader)

	if err == nil {
		t.Error("Expected an error for invalid input, but got nil")
	}
}

func TestNewRecordSliceFromFile(t *testing.T) {
	// Create a temporary file for testing
	content := `Label 1: 3.14 years
Label 2: 42 days
Label 3: 1.5 hours`
	tmpfile, err := os.CreateTemp("", "test_record_file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	enlistment, err := NewRecordEnlistmentFromFile(tmpfile.Name())

	records := enlistment.Records

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	// Check the first record
	if records[0].Label != "Label 1" || records[0].Value != 3.14 || records[0].Unit.getUnit() != int64(Year) {
		t.Errorf("First record doesn't match expected values")
	}
}

func TestNewRecordSliceFromFileNonExistentFile(t *testing.T) {
	_, err := NewRecordEnlistmentFromFile("non_existent_file.txt")

	if err == nil {
		t.Error("Expected an error for non-existent file, but got nil")
	}
}

// errorReader is a custom io.Reader that always returns an error
type errorReader struct{}

func (er errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}

func TestNewRecordSliceFromReaderScannerError(t *testing.T) {
	reader := errorReader{}

	_, err := NewRecordEnlistmentFromReader(reader)

	if err == nil {
		t.Error("Expected an error due to scanner failure, but got nil")
	}

	if err.Error() != "forced read error" {
		t.Errorf("Expected error message 'forced read error', got '%s'", err.Error())
	}
}

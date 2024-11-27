package parsing

import (
	"testing"
)

var TestUnitFiles = func() UnitFiles {
	files, err := newUnitFiles(testFS, "units")

	if err != nil {
		panic(err)
	}

	return files
}()

func TestNewMeasure(t *testing.T) {
	measure, err := newMeasure(&TestUnitFiles, "42 kilometers")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if measure.Value != 42000.0 {
		t.Errorf("Expected Value to equal %f, got %f", 42000.0, measure.Value)
	}

	if measure.Unit.value != 1000.0 {
		t.Errorf("Expected Unit.value to equal %f, got %f", 1000.0, measure.Unit.value)
	}

	if measure.Unit.Name != "kilometer" {
		t.Errorf("Expected Unit.Name to equal %s, got %s", "kilometer", measure.Unit.Name)
	}
}

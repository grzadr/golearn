package parsing

import (
	"testing"
	"testing/fstest"
)

var testJsonData = []byte(`{
    "meter": {
        "value": 1.0,
        "aliases": ["m", "meters"]
    },
    "kilometer": {
        "value": 1000.0,
        "aliases": ["km", "kilometers"]
    }
}`)

var testFS = fstest.MapFS{
	"units/test_unit.json": {
		Data: testJsonData,
	},
}

func TestLoadUnitEntriesFromJson(t *testing.T) {
	units := loadUnitEntriesFromJson(testJsonData)

	// Test number of entries
	if len(units) != 2 {
		t.Errorf("Expected 2 units, got %d", len(units))
	}

	// Test meter entry
	meter, exists := units["meter"]
	if !exists {
		t.Error("Expected 'meter' entry to exist")
	}
	if meter.Value != 1.0 {
		t.Errorf("Expected meter value to be 1.0, got %f", meter.Value)
	}
	if len(meter.Aliases) != 2 {
		t.Errorf("Expected 2 aliases for meter, got %d", len(meter.Aliases))
	}
	if meter.Aliases[0] != "m" || meter.Aliases[1] != "meters" {
		t.Errorf("Incorrect aliases for meter: %v", meter.Aliases)
	}

	// Test kilometer entry
	km, exists := units["kilometer"]
	if !exists {
		t.Error("Expected 'kilometer' entry to exist")
	}
	if km.Value != 1000.0 {
		t.Errorf("Expected kilometer value to be 1000.0, got %f", km.Value)
	}
	if len(km.Aliases) != 2 {
		t.Errorf("Expected 2 aliases for kilometer, got %d", len(km.Aliases))
	}
	if km.Aliases[0] != "km" || km.Aliases[1] != "kilometers" {
		t.Errorf("Incorrect aliases for kilometer: %v", km.Aliases)
	}
}

func TestLoadUnitEntriesFromJson_InvalidJSON(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with invalid JSON, but function did not panic")
		}
	}()

	invalidJson := []byte(`{invalid json}`)
	loadUnitEntriesFromJson(invalidJson)
}

func TestLoadUnitEntriesFromFS(t *testing.T) {
	units := loadUnitEntriesFromFS(testFS, "units/test_unit.json")

	// Verify that units were loaded correctly
	if len(units) != 2 {
		t.Errorf("Expected 2 units, got %d", len(units))
	}

	// Basic verification of content
	_, exists := units["meter"]
	if !exists {
		t.Error("Expected 'meter' entry to exist")
	}
}

func TestLoadUnitEntriesFromFS_InvalidPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with invalid path, but function did not panic")
		}
	}()

	loadUnitEntriesFromFS(testFS, "nonexistent/path.json")
}

func TestLoadUnitEntriesFromUnitsPath(t *testing.T) {
	units := loadUnitEntriesFromDefaults("units/time.json")

	// Verify that units were loaded correctly
	if len(units) != 10 {
		t.Errorf("Expected 10 units, got %d", len(units))
	}

	// Basic verification of content
	_, exists := units["Second"]
	if !exists {
		t.Error("Expected 'Second' entry to exist")
	}
}

func TestLoadUnitEntriesFromUnitsPath_InvalidPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with invalid path, but function did not panic")
		}
	}()

	loadUnitEntriesFromDefaults("nonexistent/path.json")
}

func loadUnitEntriesFilesFromDirectory

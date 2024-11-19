package parsing

import (
	"fmt"
	"path"
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

var testFSDirPath = "units"

var testFS = fstest.MapFS{
	path.Join(testFSDirPath, "test_unit.json"): {
		Data: testJsonData,
	},
	path.Join(testFSDirPath, "empty.json"): {
		Data: []byte("{}"),
	},
}

func validateUnitEntry(units UnitEntries) error {
	if len(units) != 2 {
		return fmt.Errorf("Expected 2 units, got %d", len(units))
	}

	meter, exists := units["meter"]
	if !exists {
		return fmt.Errorf("expected 'meter' entry to exist")
	}
	if meter.Value != 1.0 {
		return fmt.Errorf("expected meter value to be 1.0, got %f", meter.Value)
	}
	if len(meter.Aliases) != 2 {
		return fmt.Errorf("expected 2 aliases for meter, got %d", len(meter.Aliases))
	}
	if meter.Aliases[0] != "m" || meter.Aliases[1] != "meters" {
		return fmt.Errorf("incorrect aliases for meter: %v", meter.Aliases)
	}

	km, exists := units["kilometer"]
	if !exists {
		return fmt.Errorf("expected 'kilometer' entry to exist")
	}
	if km.Value != 1000.0 {
		return fmt.Errorf("expected kilometer value to be 1000.0, got %f", km.Value)
	}
	if len(km.Aliases) != 2 {
		return fmt.Errorf("expected 2 aliases for kilometer, got %d", len(km.Aliases))
	}
	if km.Aliases[0] != "km" || km.Aliases[1] != "kilometers" {
		return fmt.Errorf("incorrect aliases for kilometer: %v", km.Aliases)
	}

	return nil
}

func validateEmptyUnitEntry(units UnitEntries) error {
	if len(units) > 0 {
		return fmt.Errorf("Expected 0 units, got %d", len(units))
	}

	return nil
}

func TestLoadUnitEntriesFromJson(t *testing.T) {
	units := loadUnitEntriesFromJson(testJsonData)

	if err := validateUnitEntry(units); err != nil {
		t.Error(err)
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
	units := loadUnitEntriesFromFS(unitsFS, "units/time.json")

	// Verify that units were loaded correctly
	if len(units) != 10 {
		t.Errorf("Expected 10 units, got %d", len(units))
	}

	// Basic verification of content
	second_entry, exists := units["second"]
	if !exists {
		t.Error("Expected 'second' entry to exist")
	}

	if second_entry.Value != 1.0 {
		t.Errorf("Expected value 1.0, got %f instead", second_entry.Value)
	}
}

func TestLoadUnitEntriesFromUnitsPath_InvalidPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with invalid path, but function did not panic")
		}
	}()

	loadUnitEntriesFromFS(unitsFS, "nonexistent/path.json")
}

func TestLoadUnitEntriesFilesFromDirectory(t *testing.T) {
	entries_files := loadUnitEntriesFilesFromDirectory(testFS, testFSDirPath)

	expected_files := 2

	if num_files := len(entries_files); num_files != expected_files {
		t.Errorf("Expected %d UnitEntriesFiles, loaded %d", expected_files, num_files)
	}

	test_unit_entry, found := entries_files["test_unit"]

	if !found {
		t.Error("Expected to find `test_unit`")
	}

	if err := validateUnitEntry(test_unit_entry); err != nil {
		t.Error(err)
	}

	test_unit_entry, found = entries_files["empty"]

	if !found {
		t.Error("Expected to find `empty`")
	}

	if err := validateEmptyUnitEntry(test_unit_entry); err != nil {
		t.Error(err)
	}
}

func TestLoadUnitEntriesFilesFromEmbedded(t *testing.T) {
	entries_files := loadUnitEntriesFilesFromEmbedded()

	expected_content := map[string]int{
		"time":   10,
		"length": 21,
	}

	if num_files := len(entries_files); num_files != len(expected_content) {
		t.Errorf("Expected %d UnitEntriesFiles, loaded %d", len(expected_content), num_files)
	}

	for filename, expected := range expected_content {
		entry, found := entries_files[filename]

		if !found {
			t.Errorf("Expected to find `%s`", filename)
		}

		if entries_n := len(entry); entries_n != expected {
			t.Errorf("Expected %d UnitEntriesFiles, loaded %d", expected, entries_n)
		}
	}
}

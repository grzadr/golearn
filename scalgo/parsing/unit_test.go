package parsing

import (
	"encoding/json"
	"fmt"
	"path"
	"testing"
	"testing/fstest"
)

type TestEntryRecords map[string]TestEntry

var testJsonMap = TestEntryRecords{
	"meter": {
		Value:   1.0,
		Aliases: []string{"m", "meters"},
	},
	"kilometer": {
		Value:   1000.0,
		Aliases: []string{"km", "kilometers"},
	},
}

var testJsonData []byte = func() []byte {
	data, err := json.Marshal(testJsonMap)
	if err != nil {
		panic(err)
	}
	return data
}()

var testFSDirPath = "units"

var testFS = fstest.MapFS{
	path.Join(testFSDirPath, "test_unit.json"): {
		Data: testJsonData,
	},
	path.Join(testFSDirPath, "empty.json"): {
		Data: []byte("{}"),
	},
}

func helpVerifyLengths(len_t int, len_r int) error {
	if len_r == 0 {
		return fmt.Errorf("UnitRecords has 0 elements")
	} else if len_t != len_r {
		return fmt.Errorf("Expected UnitRecords of length %d, got %d instead", len_t, len_r)
	}

	return nil
}

func helpCompareTestUnitEntries(name string, t TestEntry, u Unit) []error {
	errors := make([]error, 0, 2)

	if name != u.Name {
		errors = append(errors, fmt.Errorf("Unit %s contains wrong name %s", name, u.Name))
	}

	if t.Value != u.value {
		errors = append(errors, fmt.Errorf("Value for Unit %s expected to be %f, got %f instead", name, t.Value, u.value))
	}

	return errors
}

func helpVerifyEntryExists(name string, alias string, ref TestEntry, records *UnitRecords) []error {
	errors := make([]error, 0, 2)

	r_entry, found := (*records)[alias]

	if !found {
		errors = append(errors, fmt.Errorf("Unit %s/%s is missing from UnitRecords", name, alias))
		return errors
	}

	if err := helpCompareTestUnitEntries(name, ref, r_entry); err != nil {
		errors = append(errors, err...)
	}

	return errors
}

func helpCompareUnitRecords(expected *TestEntryRecords, records *UnitRecords) []error {
	errors := make([]error, 0, 16)

	if err := helpVerifyLengths(len(*expected), len(*expected)); err != nil {
		errors = append(errors, err)
	}

	for name, ref := range *expected {
		errors = append(errors, helpVerifyEntryExists(name, name, ref, records)...)

		for _, alias := range ref.Aliases {
			errors = append(errors, helpVerifyEntryExists(name, alias, ref, records)...)
		}

	}

	return errors
}

func TestNewUnitRecords(t *testing.T) {
	records, err := newUnitRecords(testJsonData)

	if err != nil {
		t.Error(err)
		return
	}

	errors := helpCompareUnitRecords(&testJsonMap, &records)

	if num_errors := len(errors); num_errors > 0 {
		t.Errorf("Detected %d errors", num_errors)
	}

	for i, err := range errors {
		t.Error(err)
		if i > 9 {
			break
		}
	}
}

func TestNewUnitRecords_Error(t *testing.T) {
	records, err := newUnitRecords([]byte(`{
		"meter": {
			"aliases": ["m", "meters"]
		},
	}`))

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "invalid entry \"meter\": positive non-zero value field is required" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	if length := len(records); length > 0 {
		t.Errorf("Expcted empty slice, got %d elements instead", length)
	}
}

func TestUnitRecords_FindUnit(t *testing.T) {
	records, err := newUnitRecords(testJsonData)

	if err != nil {
		t.Error(err)
		return
	}

	existing_name := "meter"

	unit, found := records.findUnit(existing_name)

	if !found {
		t.Errorf("Expected to find \"%s\" Unit", existing_name)
	}

	ref := testJsonMap[existing_name]

	errors := helpCompareTestUnitEntries(existing_name, ref, unit)

	if num := len(errors); num > 0 {
		t.Errorf("Found %d errors", num)
		for _, err := range errors {
			t.Error(err)
		}
	}
}

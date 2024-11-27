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
	path.Join(testFSDirPath, "text.txt"): {
		Data: []byte("FILE"),
	},
}

func helpVerifyLengths(len_t int, len_r int) error {
	if len_t != len_r {
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

func TestNewUnitRecordsFromFS(t *testing.T) {
	records, err := newUnitRecordsFromFS(testFS, path.Join(testFSDirPath, "test_unit.json"))

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	if len(records) == 0 {
		t.Error("Units is empty")
		return
	}

	errors := helpCompareUnitRecords(&testJsonMap, &records)

	if num_errors := len(errors); num_errors > 0 {
		t.Errorf("Found %d errors", num_errors)
		for _, err := range errors {
			t.Error(err)
		}
	}
}

func TestNewUnitRecordsFromFS_Error(t *testing.T) {
	records, err := newUnitRecordsFromFS(testFS, path.Join(testFSDirPath, "invalid.json"))

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "failed to read units/invalid.json: open units/invalid.json: file does not exist" {
		t.Errorf("Expected different error message: %v", err)
	}

	if num := len(records); num > 0 {
		t.Errorf("Expected empty records, got %d instead", num)
	}
}

func TestNewUnitFiles(t *testing.T) {
	files, err := newUnitFiles(testFS, "units")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if l := len(files); l != 2 {
		t.Errorf("Excepted 2 files, got %d instead", l)
	}

	records, found := files["test_unit"]

	if !found {
		t.Errorf("Expected to find \"test_unit\"")
	}

	errors := helpCompareUnitRecords(&testJsonMap, &records)

	if num_errors := len(errors); num_errors > 0 {
		t.Errorf("Found %d errors", num_errors)
		for _, err := range errors {
			t.Error(err)
		}
	}

	records, found = files["empty"]

	if !found {
		t.Errorf("Expected to find \"empty\"")
	}

	if l := len(records); l > 0 {
		t.Errorf("Expected empty records, got %d instead", l)
	}

	records, found = files["text"]

	if found {
		t.Errorf("Expected to not find \"text\" records")
	}
}

func TestNewUnitFiles_WrongPath(t *testing.T) {
	_, err := newUnitFiles(testFS, "wrong")

	if err == nil {
		t.Error("Expected an error to occur")
	} else if err.Error() != "open wrong: file does not exist" {
		t.Errorf("Expected different error: %v", err)
	}

}

func TestNewUnitFiles_WrongJSON(t *testing.T) {
	var wrongFS = fstest.MapFS{
		path.Join(testFSDirPath, "wrong.json"): {
			Data: []byte(`{
				"meter": {
					"value": 0.0
				},
			}`),
		},
	}

	_, err := newUnitFiles(wrongFS, testFSDirPath)

	if err == nil {
		t.Error("Expected an error to occur")
	} else if err.Error() != "invalid entry \"meter\": positive non-zero value field is required" {
		t.Errorf("Expected different error: %v", err)
	}
}

func TestLoadUnitEntriesFilesFromEmbedded(t *testing.T) {
	entries_files, err := newUnitFilesFromEmbedded()

	if err != nil {
		t.Errorf("Function returned unexpected error: %v", err)
	}

	expected_content := map[string]int{
		"time":   32,
		"length": 65,
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

func TestNewUnitFilesFindUnit(t *testing.T) {
	files, err := newUnitFiles(testFS, "units")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	file_name, unit, found := files.findUnit("kilometers")

	if !found {
		t.Error("Expected to find kilometers")
	}

	if unit.Name != "kilometer" || unit.value != 1000.0 {
		t.Errorf("Wrong unit retrieved %v", unit)
	}

	if file_name != "test_unit" {
		t.Errorf("Expected file name to be \"test_unit\", got %s", file_name)
	}
}

func TestNewUnitFilesFindUnit_NotFound(t *testing.T) {
	files, err := newUnitFiles(testFS, "units")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	file_name, unit, found := files.findUnit("light years")

	if found {
		t.Error("Not expected to find \"light years\"")
	}

	if unit.Name != "" || unit.value != 0.0 {
		t.Errorf("Wrong unit retrieved %v", unit)
	}

	if file_name != "" {
		t.Errorf("Expected file name to be \"\", got %s", file_name)
	}
}

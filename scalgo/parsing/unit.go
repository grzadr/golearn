//go:build !test
// +build !test

package parsing

import (
	"embed"
	// "encoding/json"
	"fmt"
	"io/fs"
)

//go:embed units/*.json
var unitsFS embed.FS

const UNITS_PATH = "units"

type Unit struct {
	Name       string
	multiplier float64
}

func (u *Unit) isEmpty() bool {
	return u.multiplier == 0.0
}

type UnitRecords map[string]Unit

type UnitFiles map[string]UnitRecords

var embedded_units *UnitFiles

func newUnitRecords(json_data []byte) (UnitRecords, error) {
	result := make(UnitRecords)

	for _, next := range IterUnitEntries(json_data) {
		if next.Err != nil {
			return make(UnitRecords), next.Err
		}

		unit := Unit{Name: next.Entry.Name, multiplier: next.Entry.Value}

		result[unit.Name] = unit

		for _, alias := range next.Entry.Aliases {
			result[alias] = unit
		}
	}

	return result, nil
}

func (r *UnitRecords) findUnit(alias string) (Unit, bool) {
	unit, found := (*r)[alias]

	return unit, found
}

func newUnitRecordsFromFS(fsys fs.FS, entries_path string) (UnitRecords, error) {

	data, err := fs.ReadFile(fsys, entries_path)

	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", entries_path, err)
	}

	return newUnitRecords(data)

}

func newUnitFiles(fsys fs.FS, dir_path string) (UnitFiles, error) {
	files := make(UnitFiles)

	for walk_entry, err := range walkFS(fsys, dir_path) {
		if err != nil {
			return files, err
		}

		if !walk_entry.isJSONFile() {
			continue
		}

		unit_entry, err := newUnitRecordsFromFS(fsys, walk_entry.Path)

		if err != nil {
			return nil, err
		}

		files[walk_entry.Name] = unit_entry
	}

	return files, nil
}

func (u *UnitFiles) findUnit(alias string) (
	unit_file string,
	unit Unit,
	found bool,
) {
	for unit_file, records := range *u {
		unit, found := records.findUnit(alias)

		if found {
			return unit_file, unit, true
		}
	}

	return "", Unit{}, false
}

func newUnitFilesFromEmbedded() (UnitFiles, error) {
	return newUnitFiles(unitsFS, UNITS_PATH)
}

func init() {
	files, err := newUnitFilesFromEmbedded()

	if err != nil {
		panic(err)
	}

	embedded_units = &files
}

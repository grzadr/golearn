//go:build !test
// +build !test

package parsing

import (
	"embed"
	"fmt"
	"io/fs"
	"maps"
	"sort"
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

type UnitsMap map[string]Unit

type UnitRecords struct {
	units     UnitsMap
	base_unit Unit
}

func makeUnitRecords() UnitRecords {
	return UnitRecords{
		units:     make(UnitsMap),
		base_unit: Unit{},
	}
}

func (ur *UnitRecords) setUnit(unit Unit) {
	ur.units[unit.Name] = unit

	if unit.multiplier == 1.0 {
		ur.base_unit = unit
	}
}

func (ur *UnitRecords) setUnitAliased(alias string, unit Unit) {
	ur.units[alias] = unit
}

func (ur *UnitRecords) findUnit(alias string) (unit Unit, found bool) {
	unit, found = ur.units[alias]

	return unit, found
}

func (ur *UnitRecords) Len() int {
	return len(ur.units)
}

type UnitSlice []Unit

type UnitFiles map[string]UnitRecords

var EmbeddedUnits *UnitFiles

func newUnitRecords(json_data []byte) (records UnitRecords, err error) {
	records = makeUnitRecords()

	for _, next := range IterUnitEntries(json_data) {
		if next.Err != nil {
			return records, next.Err
		}

		unit := Unit{Name: next.Entry.Name, multiplier: next.Entry.Value}

		records.setUnit(unit)

		for _, alias := range next.Entry.Aliases {
			records.setUnitAliased(alias, unit)
		}
	}

	return records, nil
}

func (ur *UnitRecords) makeOrderedSliceUpTo(last_value float64) (
	units UnitSlice,
) {
	if ur.Len() == 0 {
		return units
	}

	units = make(UnitSlice, 0, ur.Len())

	seen := make(map[string]struct{}, ur.Len())

	for unit := range maps.Values(ur.units) {
		if _, found := seen[unit.Name]; found {
			continue
		}
		seen[unit.Name] = struct{}{}
		if unit.multiplier <= last_value {
			units = append(units, unit)
		}
	}

	sort.Slice(
		units,
		func(i, j int) bool {
			return units[i].multiplier > units[j].multiplier
		},
	)

	return units
}

func newUnitRecordsFromFS(fsys fs.FS, entries_path string) (UnitRecords, error) {
	data, err := fs.ReadFile(fsys, entries_path)

	if err != nil {
		return makeUnitRecords(), fmt.Errorf("failed to read %s: %w", entries_path, err)
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

	EmbeddedUnits = &files
}

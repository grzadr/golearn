package parsing

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	// "path"
	// "strings"
)

//go:embed units/*.json
var unitsFS embed.FS

var units_path = "units"

// Unit represents a single unit entry with its value and aliases
type UnitEntry struct {
	Value   float64  `json:"value"`
	Aliases []string `json:"aliases"`
}

type UnitEntries map[string]UnitEntry

type UnitEntriesFiles map[string]UnitEntries

func loadUnitEntriesFromJson(json_data []byte) UnitEntries {
	units := make(UnitEntries)
	if err := json.Unmarshal(json_data, &units); err != nil {
		panic(fmt.Sprintf("failed to unmarshal json data: %v", err))
	}

	return units
}

func loadUnitEntriesFromFS(fsys fs.FS, entries_path string) UnitEntries {
	data, err := fs.ReadFile(fsys, entries_path)
	if err != nil {
		panic(fmt.Sprintf("failed to read %s: %v", entries_path, err))
	}

	return loadUnitEntriesFromJson(data)
}


func loadUnitEntriesFilesFromDirectory(fsys fs.FS, dir_path string) UnitEntriesFiles{
	entries = make(UnitEntriesFiles)

	return entries
}



type Unit struct {
    unit float64
}

type Units map[string]Unit
type UnitAliases map[string]*Unit

type UnitsRecord struct {
    units Units
    aliases UnitAliases
}

type UnitsRegistry map[string]UnitsRecord

var registry *UnitsRegistry

func


// type Units map[string]UnitEntries

// type Unit float64

// // Units holds all units of a specific type mapped by their canonical names

// // UnitRegistry holds all unit types (time, distance, weight, etc.)
// type UnitRegistry struct {
// 	units    map[string]Units      // maps unit type to its units (e.g., "time" -> TimeUnits)
// 	aliasMap map[string]aliasEntry // quick lookup for all aliases across all unit types
// }

// // aliasEntry helps identify which unit type an alias belongs to
// type aliasEntry struct {
// 	unitType      string // e.g., "time", "distance"
// 	canonicalName string // e.g., "Hour", "Kilometer"
// }

// var units_path = "units"

// // registry is the singleton instance of UnitRegistry
// var registry *UnitRegistry

// func makeUnitRegistry() *UnitRegistry {
// 	return &UnitRegistry{
// 		units:    make(map[string]Units),
// 		aliasMap: make(map[string]aliasEntry),
// 	}
// }

// func readDirEntries(units_path string) []fs.DirEntry {
// 	entries, err := unitsFS.ReadDir(units_path)
// 	if err != nil {
// 		panic(fmt.Sprintf("failed to read units directory: %v", err))
// 	}

// 	return entries
// }

// func extractUnitType(entry fs.DirEntry) string {
// 	return strings.TrimSuffix(entry.Name(), ".json")
// }

// func readUnitFile(units_path string, entry fs.DirEntry) Units {
// 	// Read and parse the file
// 	data, err := unitsFS.ReadFile(path.Join(units_path, entry.Name()))
// 	if err != nil {
// 		panic(fmt.Sprintf("failed to read %s: %v", entry.Name(), err))
// 	}

// 	units := make(Units)
// 	if err := json.Unmarshal(data, &units); err != nil {
// 		panic(fmt.Sprintf("failed to unmarshal %s: %v", entry.Name(), err))
// 	}

// 	return units
// }

// func initializeUnitRegistry() {
//     registry = makeUnitRegistry()

//     for _, entry := range readDirEntries(units_path) {
//         if entry.IsDir() || path.Ext(entry.Name()) != ".json" {
//             continue
//         }

//         unitType := extractUnitType(entry)

//         units := readUnitFile(units_path, entry)

//         registry.units[unitType] = units

//         // Build the alias map for this unit type
//         for canonicalName, unit := range units {
//             for _, alias := range unit.Aliases {
//                 registry.aliasMap[strings.ToLower(alias)] = aliasEntry{
//                     unitType:      unitType,
//                     canonicalName: canonicalName,
//                 }
//             }
//         }
//     }
// }

func init() {
	// initializeUnitRegistry()
}

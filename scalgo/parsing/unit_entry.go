package parsing

import (
    "bytes"
    "encoding/json"
    "fmt"
    "iter"
)

// UnitEntry represents a single unit definition.
// Fields are exported to work with json.Decoder
type UnitEntry struct {
    Name    string   `json:"-"`     // Filled from JSON key
    Value   float64  `json:"value"`
    Aliases []string `json:"aliases"`
}

// NextUnitEntry wraps a UnitEntry with potential error.
// This is cleaner than having a separate error field.
type NextUnitEntry struct {
    Entry UnitEntry
    Err   error
}

// IterUnitEntries returns an iterator over unit entries in JSON data.
// The iterator yields an index and Result for each entry.
func IterUnitEntries(jsonData []byte) iter.Seq2[int, NextUnitEntry] {
    return func(yield func(int, NextUnitEntry) bool) {
        decoder := json.NewDecoder(bytes.NewReader(jsonData))

        // Check for opening delimiter
        if err := expectToken(decoder, json.Delim('{')); err != nil {
            yield(0, NextUnitEntry{Err: err})
            return
        }

        // Iterate through entries
        for i := 0; decoder.More(); i++ {
            entry, err := parseNextEntry(decoder)
            if err != nil {
                yield(i, NextUnitEntry{Err: err})
                return
            }

            if !yield(i, NextUnitEntry{Entry: entry}) {
                return
            }
        }
    }
}

// expectToken checks for an expected JSON token.
func expectToken(decoder *json.Decoder, expected json.Delim) error {
    token, err := decoder.Token()
    if err != nil {
        return fmt.Errorf("reading token: %w", err)
    }

    delim, ok := token.(json.Delim)
    if !ok || delim != expected {
        return fmt.Errorf("expected %v, got %v", expected, token)
    }
    return nil
}

// parseNextEntry reads the next unit entry from the decoder.
func parseNextEntry(decoder *json.Decoder) (UnitEntry, error) {
    // Read key (unit name)
    key, err := decoder.Token()
    if err != nil {
        return UnitEntry{}, fmt.Errorf("reading key: %w", err)
    }

    name, ok := key.(string)
    if !ok {
        return UnitEntry{}, fmt.Errorf("expected string key, got %T", key)
    }

    // Read and decode the value
    var entry UnitEntry
    if err := decoder.Decode(&entry); err != nil {
        return UnitEntry{}, fmt.Errorf("decoding value for %q: %w", name, err)
    }

    entry.Name = name
    return entry, nil
}

// // "io/fs"

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"iter"
// )

// type newUnitEntry struct {
// 	Name    string   `json:"-,"`
// 	Value   float64  `json:"value"`
// 	Aliases []string `json:"aliases"`
// }

// type nextUnitEntry struct {
// 	entry newUnitEntry
// 	err   error
// }

// type unitEntrySlice []newUnitEntry

// func iterUnitEntry(json_data []byte) iter.Seq2[int, nextUnitEntry] {
// 	return func(yield func(int, nextUnitEntry) bool) {

// 		i := 0

// 		decoder := json.NewDecoder(bytes.NewReader(json_data))
// 		if _, err := decoder.Token(); err != nil {
// 			yield(i, nextUnitEntry{newUnitEntry{}, fmt.Errorf("expected opening delimiter: %w", err)})
// 			return
// 		}

// 		for decoder.More() {
// 			// Read the key (unit name)
// 			key, err := decoder.Token()
// 			if err != nil {
// 				yield(i, nextUnitEntry{newUnitEntry{}, fmt.Errorf("error reading key: %w", err)})
// 				return
// 			}

// 			// Read the value (UnitEntry)
// 			entry := newUnitEntry{
// 				Name: key.(string),
// 			}
// 			if err := decoder.Decode(&entry); err != nil {
// 				yield(i, nextUnitEntry{newUnitEntry{}, fmt.Errorf("error reading value for %q: %w", key, err)})
// 				return
// 			}

// 			if !yield(i, nextUnitEntry{entry, nil}) {
// 				return
// 			}

// 			i++
// 		}
// 	}
// }

// func loadUnitEntriesStreaming(jsonData []byte) error {
//     decoder := json.NewDecoder(bytes.NewReader(jsonData))

//     // Read opening delimiter '{'
//     if _, err := decoder.Token(); err != nil {
//         return fmt.Errorf("expected opening delimiter: %w", err)
//     }

//     // While the array contains values
//     for decoder.More() {
//         // Read the key (unit name)
//         key, err := decoder.Token()
//         if err != nil {
//             return fmt.Errorf("error reading key: %w", err)
//         }

//         // Read the value (UnitEntry)
//         var entry UnitEntry
//         if err := decoder.Decode(&entry); err != nil {
//             return fmt.Errorf("error reading value for %q: %w", key, err)
//         }

//         // Process each entry as it's decoded
//         processEntry(key.(string), entry)
//     }

//     return nil
// }

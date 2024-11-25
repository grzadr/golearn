package parsing

import (
	"testing"
)

var testUnitEntryJsonData = []byte(`{
    "meter": {
        "value": 1.0,
        "aliases": ["m", "meters"]
    },
    "kilometer": {
        "value": 1000.0,
        "aliases": ["km", "kilometers"]
    }
}`)

var testExpectedData = []UnitEntry{
	UnitEntry{
		Name: "meter",
		Value: 1.0,
		Aliases: []string{"m", "meters"},
	},
	UnitEntry{
		Name: "kilometer",
		Value: 1000.0,
		Aliases: []string{"km", "kilometers"},
	},
}


func helpCompareUnitEntry(a *UnitEntry, b *UnitEntry) bool {
	if a.Name != b.Name || a.Value != b.Value || len(a.Aliases) != len(b.Aliases) {
		return false
	}

	a_aliases := a.Aliases
	b_aliases := b.Aliases

	for i := 0; i < len(a_aliases); i++ {
		if a_aliases[i] != b_aliases[i] {
			return false
		}
	}

	return true
}



func TestIterUnitEntry(t *testing.T) {
	for i, next := range IterUnitEntries(testUnitEntryJsonData) {
		if next.Err != nil {
			t.Errorf("Received unexpected error %v", next.Err)
		}

		if !helpCompareUnitEntry(&next.Entry, &testExpectedData[i]) {
			t.Errorf("Expected entry %v, got %v", next.Entry, testExpectedData[i])
			return
		}
	}
}

func TestIterUnitEntry_EarlyTermination(t *testing.T) {
	last_i := 0
	for i, next := range IterUnitEntries(testUnitEntryJsonData) {
		last_i = i
		if next.Err != nil {
			t.Errorf("Received unexpected error %v", next.Err)
		}

		if !helpCompareUnitEntry(&next.Entry, &testExpectedData[i]) {
			t.Errorf("Expected entry %v, got %v", next.Entry, testExpectedData[i])
			return
		}

		break
	}

	if last_i > 0 {
		t.Errorf("Expected to terminate at 0, but terminated at %d", last_i)
	}
}

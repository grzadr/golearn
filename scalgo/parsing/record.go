package parsing

import (
	"fmt"
	"strings"
)

type Record struct {
	Label string
	Measure
}

func (r *Record) isEmpty() bool {
	return len(r.Label) == 0 && r.Measure.isEmpty()
}

func newRecord(query string, unit_files *UnitFiles) (Record, error) {
	label, measure_str, found := strings.Cut(query, ": ")

	if !found {
		return Record{}, fmt.Errorf("Record `%s` missing `: `", query)
	}

	if len(label) == 0 {
		return Record{}, fmt.Errorf("Record `%s` missing label", query)
	}

	measure, err := newMeasure(unit_files, measure_str)

	if err != nil {
		return Record{}, err
	}
	return Record{Label: strings.TrimSpace(label), Measure: measure}, nil
}

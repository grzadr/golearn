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

func (r *Record) scale(ref *Record, scale *Measure) Record {
	return Record {
		Label: r.Label,
		Measure: r.Measure.scale(&ref.Measure, scale),
	}
}

func splitRecordString(str string) (label, measure_str string, err error) {
	var found bool
	label, measure_str, found = strings.Cut(str, ": ")

	if !found {
		return label, measure_str, fmt.Errorf("Record `%s` missing `: `", str)
	}

	label = strings.TrimSpace(label)

	if len(label) == 0 {
		return label, measure_str, fmt.Errorf("Record `%s` missing label", str)
	}

	return label, measure_str, nil
}

func newRecord(record_str string, unit_files *UnitFiles) (record Record, unit_file string, err error) {
	label, measure_str, err := splitRecordString(record_str)

	if err != nil {
		return record, unit_file, err
	}

	var measure Measure

	measure, unit_file, err = newMeasure(unit_files, measure_str)

	if err != nil {
		return record, unit_file, err
	}

	record = Record{
		Label:   label,
		Measure: measure,
	}
	return record, unit_file, nil
}

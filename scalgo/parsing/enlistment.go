package parsing

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"math"
	"sort"
	"strings"
)

type EnlistmentOptions struct {
	sort            bool
	reversed        bool
	scale           *Measure
	scale_unit_file string
}

func isTrue(has_value bool, value string) bool {
	return !has_value || value == "true"
}

func isFalse(value string) bool {
	return value == "false"
}

func setFlag(flag *bool, has_value bool, value string) error {
	if isTrue(has_value, value) {
		*flag = true
	} else if isFalse(value) {
		*flag = false
	} else {
		return fmt.Errorf("Unknown value: %s", value)
	}
	return nil
}
func (o *EnlistmentOptions) setSort(has_value bool, value string) error {
	return setFlag(&o.sort, has_value, value)
}
func (o *EnlistmentOptions) setReversed(has_value bool, value string) error {
	return setFlag(&o.reversed, has_value, value)
}

func (o *EnlistmentOptions) setScale(has_value bool, value string, unit_files *UnitFiles) error {
	if !has_value {
		return fmt.Errorf("Missing measure")
	}
	scale, unit_file, err := newMeasure(unit_files, value)

	if err != nil {
		return fmt.Errorf("Wrong scale measure: %w", err)
	}

	o.scale = &scale
	o.scale_unit_file = unit_file

	return nil
}

func (o *EnlistmentOptions) hasScale() bool {
	return o.scale != nil
}

func newOptions() EnlistmentOptions {
	return EnlistmentOptions{sort: true,
		reversed: false,
		scale:    &Measure{}}
}

type RecordSlice []Record

func (rs *RecordSlice) Str(units *UnitRecords, ref Record, max_units int) iter.Seq[string] {
	return func(yield func(string) bool) {
		picked_units := units.makeOrderedSliceUpTo(ref.getBaseValue())

		for _, record := range *rs {
			if !yield(record.str(&picked_units, max_units)) {
				return
			}
		}
	}
}

type Enlistment struct {
	options   EnlistmentOptions
	records   RecordSlice
	ref       *Record
	unit_file string
}

const CommentPrefix = "#"
const SettingPrefix = "@"

func NewEnlistment() *Enlistment {
	return &Enlistment{
		records: make(RecordSlice, 0, 32),
		ref:     nil,
		options: newOptions(),
	}
}

func prepareOption(option string) (name string, value string, has_value bool) {
	name, value, has_value = strings.Cut(option, " ")

	value = strings.TrimSpace(value)

	if name != "@scale" {
		value = strings.ToLower(value)
	}

	return
}

func (e *Enlistment) applyOption(option string, unit_files *UnitFiles) error {
	name, value, has_value := prepareOption(option)

	var err error

	switch name {
	case "@scale":
		err = e.options.setScale(has_value, value, unit_files)
	case "@sort":
		err = e.options.setSort(has_value, value)
	case "@reverse":
		err = e.options.setReversed(has_value, value)
	default:
		err = fmt.Errorf("Unknown option: %s", option)
	}

	return err
}

func (e *Enlistment) enabledSort() bool {
	return e.options.sort
}

func (e *Enlistment) enabledReversed() bool {
	return e.options.reversed
}

func (e *Enlistment) SortRecords() {
	if len(e.records) == 0 {
		return
	}

	sort.Slice(e.records, func(i, j int) bool {
		if e.enabledReversed() {
			return e.records[i].getBaseValue() < e.records[j].getBaseValue()
		}
		return e.records[i].getBaseValue() > e.records[j].getBaseValue()
	})

	e.options.sort = true
}

func (e *Enlistment) detectIssue() error {
	detected := make([]error, 0, 2)
	if len(e.records) == 0 {
		detected = append(detected, fmt.Errorf("Enlisting is missing records"))
	}

	if !e.options.hasScale() {
		detected = append(detected, fmt.Errorf("Enlisting is missing @scale"))
	}

	if len(detected) > 0 {
		return errors.Join(detected...)
	}

	return nil
}

func (e *Enlistment) findRefRecord() (ref *Record) {
	if e.ref != nil {
		return e.ref
	}

	if e.enabledSort() {
		return &e.records[0]
	}

	compareFunc := math.Max

	if e.enabledReversed() {
		compareFunc = math.Min
	}

	ref = &e.records[0]

	for _, record := range e.records[1:] {
		compared := compareFunc(float64(ref.getBaseValue()), float64(record.getBaseValue()))

		if ref.Value != compared {
			ref = &record
		}
	}

	return ref
}

func (e *Enlistment) appendLine(line string, unit_files *UnitFiles) error {
	record, unit_file, err := newRecord(line, unit_files)
	if err != nil {
		return err
	}

	if e.unit_file == "" {
		e.unit_file = unit_file
	} else if e.unit_file != unit_file {
		return fmt.Errorf(
			"Expected unit from %s for line `%s`, got %s",
			e.unit_file,
			line,
			unit_file,
		)
	}

	e.records = append(e.records, record)

	return nil
}

func (e *Enlistment) getScaledRecords(unit_files *UnitFiles) (
	records RecordSlice, err error,
) {
	scale := e.options.scale
	records = make(RecordSlice, 0, len(e.records))

	units, found := (*unit_files)[e.unit_file]

	if !found {
		return records, fmt.Errorf("Failed to find unit file %s", e.unit_file)
	}

	for _, r := range e.records {
		scaled := r.scale(e.ref, scale)
		scaled.Measure.Unit = units.base_unit
		records = append(records, scaled)
	}

	return records, nil
}

func (e *Enlistment) ScaleRecords(unit_files *UnitFiles, max_units int) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		scaled, err := e.getScaledRecords(unit_files)

		if err != nil {
			yield("", err)
			return
		}
		units := (*unit_files)[e.unit_file]
		for s := range scaled.Str(&units, *e.ref, max_units) {
			if !yield(s, nil) {
				return
			}
		}
	}
}

func scanReaderIntoEnlistment(reader io.Reader, unit_files *UnitFiles) (enlistment *Enlistment, err error) {
	enlistment = NewEnlistment()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, CommentPrefix) {
			continue
		}
		if strings.HasPrefix(line, SettingPrefix) {
			// Parse settings
			if err := enlistment.applyOption(line, unit_files); err != nil {
				return enlistment, err
			}
			continue
		}

		if err := enlistment.appendLine(line, unit_files); err != nil {
			return enlistment, err
		}
	}

	if err := scanner.Err(); err != nil {
		return enlistment, err
	}

	return enlistment, nil
}

func NewEnlistmentFromReader(reader io.Reader, unit_files *UnitFiles) (*Enlistment, error) {
	enlistment, err := scanReaderIntoEnlistment(reader, unit_files)
	if err != nil {
		return enlistment, err
	}

	if err := enlistment.detectIssue(); err != nil {
		return enlistment, err
	}

	if enlistment.options.sort {
		enlistment.SortRecords()
	}

	enlistment.ref = enlistment.findRefRecord()

	return enlistment, nil
}

func NewEnlistmentFromFile(
	fsys fs.FS,
	filename string,
	unit_files *UnitFiles,
) (*Enlistment, error) {
	file, err := fsys.Open(filename)

	if err != nil {
		return &Enlistment{}, err
	}
	defer file.Close()

	return NewEnlistmentFromReader(file, unit_files)
}

package parsing

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

type Options struct {
	sort     bool
	reversed bool
	scale    *Measure
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
func (o *Options) setSort(has_value bool, value string) error {
	return setFlag(&o.sort, has_value, value)
}
func (o *Options) setReversed(has_value bool, value string) error {
	return setFlag(&o.reversed, has_value, value)
}

func (o *Options) setScale(has_value bool, value string, unit_files *UnitFiles) error {
	if !has_value {
		return fmt.Errorf("Missing measure")
	}
	scale, err := newMeasure(unit_files, value)

	if err != nil {
		fmt.Errorf("Wrong scale measure: %w", err)
	}

	o.scale = &scale

	return nil
}

func (o *Options) hasScale() bool {
	return o.scale != nil
}

func newOptions() Options {
	return Options{sort: true,
		reversed: false,
		scale:    &Measure{}}
}

type Enlistment struct {
	options Options
	records []Record
	ref     *Record
}

const CommentPrefix = "#"
const SettingPrefix = "@"

func NewEnlistment() *Enlistment {
	return &Enlistment{
		records: make([]Record, 0, 32),
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

	switch name {
	case "@scale":
		e.options.setScale(has_value, value, unit_files)
	case "@sort":
		e.options.setSort(has_value, value)
	case "@reverse":
		e.options.setReversed(has_value, value)
	default:
		fmt.Errorf("Unknown option: %s", option)
	}

	return nil
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
			return e.records[i].Value < e.records[j].Value
		}
		return e.records[i].Value > e.records[j].Value
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
		compared := compareFunc(float64(ref.Value), float64(record.Value))

		if ref.Value != compared {
			ref = &record
		}
	}

	return ref
}

func ScanReaderIntoEnlistment(reader io.Reader, unit_files *UnitFiles) (*Enlistment, error) {
	enlistment := NewEnlistment()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, CommentPrefix) {
			continue
		}
		if strings.HasPrefix(line, SettingPrefix) {
			// Parse settings
			if err := enlistment.applyOption(line, unit_files); err != nil {
				return enlistment, err
			}
			continue
		}

		record, err := newRecord(line, unit_files)
		if err != nil {
			return enlistment, err
		}
		enlistment.records = append(enlistment.records, record)
	}

	if err := scanner.Err(); err != nil {
		return enlistment, err
	}

	return enlistment, nil
}

func NewRecordEnlistmentFromReader(reader io.Reader, unit_files *UnitFiles) (*Enlistment, error) {
	enlistment, err := ScanReaderIntoEnlistment(reader, unit_files)
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

func NewRecordEnlistmentFromFile(filename string, unit_files *UnitFiles) (*Enlistment, error) {
	file, err := os.Open(filename)
	if err != nil {
		return &Enlistment{}, err
	}
	defer file.Close()

	return NewRecordEnlistmentFromReader(file, unit_files)
}

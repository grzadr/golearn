package parsing

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	Sort     bool
	Reversed bool
	Scale    Measure
}

func newOptions() Options {
	return Options{Sort: true,
		Reversed: false,
		Scale:    Measure{}}
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
		records:   make([]Record, 0, 32),
		ref:       nil,
		options: newOptions(),
	}
}

func (e *Enlistment) applyOption(option string) error {
	name, value, _ := strings.Cut(option, " ")

	switch name {
	case "@scale":
		NewU
	}
	if mapper, found := RecordEnlistmentSettingsMapper[name]; found {
		return mapper(e, value)
	}
	return fmt.Errorf("Unknown setting %s", option)
}

func ScanReaderIntoEnlistment(reader io.Reader) (*Enlistment, error) {
	enlistment := NewEnlistment()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, CommentPrefix) {
			continue
		}
		if strings.HasPrefix(line, SettingPrefix) {
			// Parse settings
			if err := enlistment.applyOption(line); err != nil {
				return enlistment, err
			}
			continue
		}

		record, err := NewRecordFromString(line)
		if err != nil {
			return enlistment, err
		}
		enlistment.Records = append(enlistment.Records, record)
	}

	if err := scanner.Err(); err != nil {
		return enlistment, err
	}

	return enlistment, nil
}

func NewRecordEnlistmentFromReader(reader io.Reader) (*Enlistment, error) {
	enlistment, err := ScanReaderIntoEnlistment(reader)
	if err != nil {
		return enlistment, err
	}

	if len(enlistment.Records) == 0 {
		return enlistment, fmt.Errorf("No records found")
	}

	if enlistment.Sorted {
		enlistment.SortRecords()
	}

	enlistment.RefRecord = enlistment.findRefRecord()

	if enlistment.ScaleUnit == nil {
		enlistment.ScaleUnit = enlistment.RefRecord.Unit
	}

	return enlistment, nil
}

func NewRecordEnlistmentFromFile(filename string) (*Enlistment, error) {
	file, err := os.Open(filename)
	if err != nil {
		return &Enlistment{}, err
	}
	defer file.Close()

	return NewRecordEnlistmentFromReader(file)
}

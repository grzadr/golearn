// ReadFileContent reads the content of a Markdown file and returns it as a byte slice.
package parsing

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type RecordUnit interface {
	ConvertToUnit(value float64, unit int64) float64
	ConvertToBaseUnit(value float64) float64
	getUnit() int64
	getType() RecordUnitType
}

type RecordUnitType int

const (
	UnknownUnitType RecordUnitType = iota
	TimeUnitType
)

type TimeUnits int64

const (
	Missing    TimeUnits = 0
	Second     TimeUnits = 1
	Minute     TimeUnits = 60
	Hour       TimeUnits = 3600
	Day        TimeUnits = 86400    // 24 hours
	Week       TimeUnits = 604800   // 7 days
	Month      TimeUnits = 2592000  // 30 days
	Year       TimeUnits = 31536000 // 365 days
	Decade     TimeUnits = 315360000
	Century    TimeUnits = 3153600000
	Millennium TimeUnits = 31536000000
)

var timeUnitMap = map[string]TimeUnits{
	// Second
	"second":  Second,
	"seconds": Second,
	"sec":     Second,
	"s":       Second,

	// Minute
	"minute":  Minute,
	"minutes": Minute,
	"min":     Minute,
	"m":       Minute,

	// Hour
	"hour":  Hour,
	"hours": Hour,
	"hr":    Hour,
	"h":     Hour,

	// Day
	"day":  Day,
	"days": Day,
	"d":    Day,

	// Week
	"week":  Week,
	"weeks": Week,
	"wk":    Week,
	"w":     Week,

	// Month
	"month":  Month,
	"months": Month,
	"mo":     Month,

	// Year
	"year":  Year,
	"years": Year,
	"yr":    Year,
	"y":     Year,

	// Decade
	"decade":  Decade,
	"decades": Decade,

	// Century
	"century":   Century,
	"centuries": Century,

	// Millennium
	"millennium": Millennium,
	"millennia":  Millennium,
}

// ParseTimeUnit converts a string to a TimeUnits constant
func parseTimeUnit(s string) (TimeUnits, bool) {
	unit, found := timeUnitMap[strings.ToLower(strings.TrimSpace(s))]
	return unit, found
}

type TimeUnit struct {
	Type RecordUnitType
	Unit TimeUnits
}

func NewTimeUnit(input string) (RecordUnit, error) {
	unit, found := parseTimeUnit(input)

	if !found {
		return nil, fmt.Errorf("Unknown unit %s", input)
	}

	return &TimeUnit{
		Type: TimeUnitType,
		Unit: unit,
	}, nil
}

func (t *TimeUnit) getUnit() int64 {
	return int64(t.Unit)
}

func (t *TimeUnit) getType() RecordUnitType {
	return t.Type
}

func (t *TimeUnit) ConvertToBaseUnit(value float64) float64 {
	return value * float64(t.Unit)
}

func (t *TimeUnit) ConvertToUnit(value float64, unit int64) float64 {
	return t.ConvertToBaseUnit(value) / float64(unit)
}

func newRecordUnit(unit string) (RecordUnit, error) {
	if unit == "" {
		return nil, nil
	}

	new_functions := []func(input string) (RecordUnit, error){
		NewTimeUnit,
	}

	for _, f := range new_functions {
		recordUnit, err := f(unit)
		if err == nil {
			return recordUnit, nil
		}
	}

	return nil, fmt.Errorf("no matching unit found for %s", unit)
}

type Record struct {
	Label     string
	BaseValue float64
	Unit      RecordUnit
}

// splitInputString splits an input string into a label, value, and unit
// The input string should be in the format "label: value unit"
func splitInputString(input string) (label string, value float64, unit string, err error) {
	label, value_with_unit, found := strings.Cut(strings.TrimSpace(input), ":")

	if !found {
		return "", 0, "", fmt.Errorf("No colon found in input string")
	}

	if label == "" {
		return "", 0, "", fmt.Errorf("No label found in input string")
	}

	value_str, unit, _ := strings.Cut(strings.TrimSpace(value_with_unit), " ")

	value, err = strconv.ParseFloat(value_str, 64)

	if err != nil {
		return "", 0, "", err
	}

	return label, value, unit, nil
}

// NewRecord creates a new Record with the given label, value, and unit string
func NewRecord(label string, value float64, unit_str string) (*Record, error) {
	unit, err := newRecordUnit(unit_str)
	if err != nil {
		return nil, fmt.Errorf("failed to create record: %w", err)
	}
	return &Record{
		Label:     label,
		BaseValue: unit.ConvertToBaseUnit(value),
		Unit:      unit,
	}, nil
}

func NewRecordFromString(input string) (*Record, error) {
	label, value, unit, err := splitInputString(input)

	if err != nil {
		return nil, err
	}

	record, err := NewRecord(label, value, unit)

	if err != nil {
		return nil, err
	}

	return record, nil
}

type RecordEnlistment struct {
	Records   []*Record
	RefRecord *Record
	ScaleUnit RecordUnit
	Sorted    bool
	Reversed  bool
}

var RecordEnlistmentSettingsMapper = map[string]func(*RecordEnlistment, string) error{
	"@scale": func(enlistment *RecordEnlistment, value string) error {
		unit, err := newRecordUnit(value)
		if err != nil {
			return err
		}
		enlistment.ScaleUnit = unit
		return nil
	},
	"@sorted": func(enlistment *RecordEnlistment, value string) error {
		enlistment.Sorted = value == "true"
		return nil
	},
	"@reverse": func(enlistment *RecordEnlistment, value string) error {
		enlistment.Reversed = value == "true"
		return nil
	},
}

func (enlistment *RecordEnlistment) ApplyRecordEnlistmentSetting(setting string) error {
	name, value, _ := strings.Cut(setting, " ")
	if mapper, found := RecordEnlistmentSettingsMapper[name]; found {
		return mapper(enlistment, value)
	}
	return fmt.Errorf("Unknown setting %s", setting)
}

func (enlistment *RecordEnlistment) findRefRecord() *Record {
	if enlistment.RefRecord != nil {
		return enlistment.RefRecord
	}

	if enlistment.Sorted {
		return enlistment.Records[0]
	}

	compareFunc := math.Max

	if enlistment.Reversed {
		compareFunc = math.Min
	}

	refRecord := enlistment.Records[0]

	for _, record := range enlistment.Records[1:] {
		compared := compareFunc(refRecord.BaseValue, record.BaseValue)

		if refRecord.BaseValue != compared {
			refRecord = record
		}
	}

	return refRecord
}

func (re *RecordEnlistment) SortRecords() {
	if len(re.Records) == 0 {
		return
	}

	sort.Slice(re.Records, func(i, j int) bool {
		if re.Reversed {
			return re.Records[i].BaseValue > re.Records[j].BaseValue
		}
		return re.Records[i].BaseValue < re.Records[j].BaseValue
	})

	re.Sorted = true
}

func NewRecordEnlistmentFromReader(reader io.Reader) (RecordEnlistment, error) {
	enlistment := RecordEnlistment{
		Records:   make([]*Record, 0),
		RefRecord: nil,
		ScaleUnit: nil,
		Sorted:    true,
		Reversed:  false,
	}
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Skip comments
			continue
		}
		if strings.HasPrefix(line, "@") {
			// Parse settings
			if err := enlistment.ApplyRecordEnlistmentSetting(line); err != nil {
				return RecordEnlistment{}, err
			}
			continue
		}

		record, err := NewRecordFromString(line)
		if err != nil {
			return RecordEnlistment{}, err
		}
		enlistment.Records = append(enlistment.Records, record)
	}

	if err := scanner.Err(); err != nil {
		return RecordEnlistment{}, err
	}

	if len(enlistment.Records) == 0 {
		return RecordEnlistment{}, fmt.Errorf("No records found")
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

func NewRecordEnlistmentFromFile(filename string) (RecordEnlistment, error) {
	file, err := os.Open(filename)
	if err != nil {
		return RecordEnlistment{}, err
	}
	defer file.Close()

	return NewRecordEnlistmentFromReader(file)
}

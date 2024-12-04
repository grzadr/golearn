package parsing

type Record struct {
	Label string
	Measure
}

func newRecord(query string, unit_files *UnitFiles) (record Record, err error) {
	return Record{}, nil
}

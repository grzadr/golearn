package parsing

type Unit struct {
    Scale float64
    Type string
}

func NewUnit(value float64, unit string) (Measure, error) {

type Measure struct {
    Value float64
    Unit
}


package parsing

type Unit struct {
    Scale float64
    Type string
}

type Measure struct {
    Value float64
    Unit
}

package builder

// Inserter creates insert statements
type Inserter interface {
	Columns(columns ...string) Inserter
	Values(values ...Value) Inserter
	Table(table string) Inserter
	Build() string
}

// Value formats the values to be used in insert clause
type Value struct {
	values []interface{}
}

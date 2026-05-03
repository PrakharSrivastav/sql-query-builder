package builder

// Inserter creates insert statements. Build returns SQL with `?`
// placeholders for values, the args slice in column-sorted order, and
// any identifier-validation error.
type Inserter interface {
	Columns([]string) Inserter
	Values(Value) Inserter
	Table(string) Inserter
	Build() (string, []any, error)
}

// Value formats the values to be used in insert clause
type Value map[string]interface{}

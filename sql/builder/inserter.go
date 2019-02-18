package builder

// Inserter creates insert statements
type Inserter interface {
	Columns([]string) Inserter
	Values(Value) Inserter
	Table(string) Inserter
	Build() string
}

// Value formats the values to be used in insert clause
type Value map[string]interface{}

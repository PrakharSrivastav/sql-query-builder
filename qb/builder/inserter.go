package builder

// Inserter creates insert statements. Build returns SQL with `?`
// placeholders for values, the args slice in column-sorted order, and
// any identifier-validation error.
type Inserter interface {
	Columns([]string) Inserter
	Values(Value) Inserter
	Table(string) Inserter
	// Returning appends a `RETURNING col1, col2, ...` clause. Each
	// argument is validated as a SQL identifier (single `*` is also
	// accepted). Postgres supports RETURNING natively; SQLite 3.35+
	// supports it; MySQL does not.
	Returning(cols ...string) Inserter
	// OnConflictDoNothing emits `ON CONFLICT (target1, ...) DO NOTHING`.
	// Each target column is validated as an identifier.
	OnConflictDoNothing(targets ...string) Inserter
	// OnConflictDoUpdate emits `ON CONFLICT (target1, ...) DO UPDATE
	// SET col=?, ...`. Set keys are validated as identifiers; values
	// are bound via placeholder unless wrapped in an Excluded sentinel,
	// in which case the right-hand side becomes `EXCLUDED.<col>`.
	// Postgres supports this natively; SQLite 3.24+ supports it;
	// MySQL has a different syntax (`ON DUPLICATE KEY UPDATE`).
	OnConflictDoUpdate(targets []string, set map[string]interface{}) Inserter
	Build() (string, []any, error)
}

// Value formats the values to be used in insert clause
type Value map[string]interface{}

// Excluded marks a value in OnConflictDoUpdate's set map as a reference
// to the conflicting row's column (Postgres `EXCLUDED.<Col>`). The Col
// is validated as an identifier.
type Excluded struct {
	Col string
}

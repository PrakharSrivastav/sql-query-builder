package builder

// Updater helps in creating update sql statements. Build returns SQL
// with `?` placeholders, the args slice (SET values first, then
// Condition args), and any identifier-validation error.
type Updater interface {
	Update(string) Updater
	Set(map[string]interface{}) Updater
	Condition(Expression) Updater
	RawCondition(string) Updater
	// Returning appends a `RETURNING col1, col2, ...` clause. Each
	// argument is validated as a SQL identifier (single `*` is also
	// accepted). Postgres supports RETURNING natively; SQLite 3.35+
	// supports it; MySQL does not.
	Returning(cols ...string) Updater
	Build() (string, []any, error)
}

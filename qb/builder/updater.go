package builder

// Updater helps in creating update sql statements. Build returns SQL
// with `?` placeholders, the args slice (SET values first, then
// Condition args), and any identifier-validation error.
type Updater interface {
	Update(string) Updater
	Set(map[string]interface{}) Updater
	Condition(Expression) Updater
	RawCondition(string) Updater
	Build() (string, []any, error)
}

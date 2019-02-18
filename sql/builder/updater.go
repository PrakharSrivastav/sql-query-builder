package builder

// Updater helps in creating update sql statements
type Updater interface {
	Update(string) Updater
	Set(map[string]interface{}) Updater
	Condition(Expression) Updater
	RawCondition(string) Updater
	Build() string
}

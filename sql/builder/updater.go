package builder

// Updater helps in creating update sql statements
type Updater interface {
	Update(table string) Updater
	Set(values map[string]string) Updater
	Expression
	Build() string
}

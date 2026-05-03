package builder

// Creater helps to generate create sql statements. CREATE TABLE has no
// runtime values; the args slice from Build is always empty. Error
// surfaces identifier-validation failures.
type Creater interface {
	SetColumns([]Columns) Creater
	Table(table string) Creater
	Build() (string, []any, error)
}

// Columns help to determine create table statements
type Columns struct {
	Name       string
	Datatype   string
	Constraint string
}

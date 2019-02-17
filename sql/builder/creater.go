package builder

// Creater helps to generate create sql statements
type Creater interface {
	SetColumns(cloumns ...Columns) Creater
	Table(table string) Creater
	Build() string
}

// Columns help to determine create table statements
type Columns struct {
	Name       string
	Datatype   string
	Constraint string
}

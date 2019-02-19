package builder

// Expression interface defines the methods to be used to evaluate a meaningful evaluate clause
type Expression interface {
	Where(clause Clause) Expression
	And(clause Clause) Expression
	Or(clause Clause) Expression
	In(field string, items ...string) Expression
	NotIn(field string, items ...string) Expression
	Express() string
}

// Clause is used to set a where, and, or clause
// example where column1 = "value1" --> Clause{Left:"column1", Operator:"=", Right="value1"}
type Clause struct {
	Left     string
	Operator string
	Right    interface{}
}

package builder

type Expression interface {
	Where(clause Clause) Expression
	And(clause Clause) Expression
	Or(clause Clause) Expression
	In(field string, items ...string) Expression
	NotIn(field string, items ...string) Expression
	Express() string
}

type Clause struct {
	Left     string
	Operator string
	Right    string
}

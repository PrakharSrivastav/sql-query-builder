package builder

// Reader provides a contract to be implemented by all sql generators
type Reader interface {
	Select(...string) Reader
	From(...string) Reader
	FromAlias(...Alias) Reader
	OrderBy(...string) Reader
	Limit(int) Reader
	Offset(int) Reader
	Condition(Expression) Reader
	RawCondition(string) Reader
	InnerJoin(string) Reader
	LeftJoin(string) Reader
	RightJoin(string) Reader
	On(string) Reader
	GroupBy([]string) Reader
	Having(string) Reader
	Build() string
}

// Alias is a struct to provide a table name alias
type Alias struct {
	Name  string
	Alias string
}

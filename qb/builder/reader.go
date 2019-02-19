package builder

// Reader provides a contract to be implemented by all sql generators
type Reader interface {
	Select(columns ...string) Reader
	From(tables ...string) Reader
	FromAlias(alias ...Alias) Reader
	OrderBy(columns ...string) Reader
	Limit(limit int) Reader
	Offset(offset int) Reader
	Condition(condition Expression) Reader
	RawCondition(condition string) Reader
	Build() string
}

// Alias is a struct to provide a table name alias
type Alias struct {
	Name  string
	Alias string
}

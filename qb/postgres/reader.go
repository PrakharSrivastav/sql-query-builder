package postgres

import (
	"github.com/PrakharSrivastav/sql-query-builder/qb/ansi"
	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Reader wraps ansi.Reader so each fluent method returns *Reader (the
// postgres wrapper) rather than the embedded ANSI type. This keeps the
// chain bound to the postgres Build, which performs the placeholder
// rewrite.
type Reader struct {
	inner *ansi.Reader
}

func (r *Reader) Select(s ...string) builder.Reader        { r.inner.Select(s...); return r }
func (r *Reader) From(s ...string) builder.Reader          { r.inner.From(s...); return r }
func (r *Reader) FromAlias(a ...builder.Alias) builder.Reader {
	r.inner.FromAlias(a...)
	return r
}
func (r *Reader) OrderBy(s ...string) builder.Reader { r.inner.OrderBy(s...); return r }
func (r *Reader) Limit(i int) builder.Reader         { r.inner.Limit(i); return r }
func (r *Reader) Offset(i int) builder.Reader        { r.inner.Offset(i); return r }
func (r *Reader) Condition(e builder.Expression) builder.Reader {
	r.inner.Condition(e)
	return r
}
func (r *Reader) RawCondition(s string) builder.Reader { r.inner.RawCondition(s); return r }
func (r *Reader) InnerJoin(t string) builder.Reader    { r.inner.InnerJoin(t); return r }
func (r *Reader) LeftJoin(t string) builder.Reader     { r.inner.LeftJoin(t); return r }
func (r *Reader) RightJoin(t string) builder.Reader    { r.inner.RightJoin(t); return r }
func (r *Reader) On(c string) builder.Reader           { r.inner.On(c); return r }
func (r *Reader) GroupBy(f []string) builder.Reader    { r.inner.GroupBy(f); return r }
func (r *Reader) Having(c string) builder.Reader       { r.inner.Having(c); return r }

func (r *Reader) Build() (string, []any, error) {
	sql, args, err := r.inner.Build()
	return rewritePlaceholders(sql), args, err
}

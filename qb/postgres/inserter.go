package postgres

import (
	"github.com/PrakharSrivastav/sql-query-builder/qb/ansi"
	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Inserter wraps ansi.Inserter to keep the fluent chain bound to the
// postgres Build (which rewrites placeholders).
type Inserter struct {
	inner *ansi.Inserter
}

func (i *Inserter) Columns(s []string) builder.Inserter { i.inner.Columns(s); return i }
func (i *Inserter) Values(v builder.Value) builder.Inserter {
	i.inner.Values(v)
	return i
}
func (i *Inserter) Table(s string) builder.Inserter { i.inner.Table(s); return i }
func (i *Inserter) Returning(cols ...string) builder.Inserter {
	i.inner.Returning(cols...)
	return i
}
func (i *Inserter) OnConflictDoNothing(targets ...string) builder.Inserter {
	i.inner.OnConflictDoNothing(targets...)
	return i
}
func (i *Inserter) OnConflictDoUpdate(targets []string, set map[string]interface{}) builder.Inserter {
	i.inner.OnConflictDoUpdate(targets, set)
	return i
}

func (i *Inserter) Build() (string, []any, error) {
	sql, args, err := i.inner.Build()
	return rewritePlaceholders(sql), args, err
}

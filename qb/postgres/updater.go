package postgres

import (
	"github.com/PrakharSrivastav/sql-query-builder/qb/ansi"
	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Updater wraps ansi.Updater to keep the fluent chain bound to the
// postgres Build (which rewrites placeholders).
type Updater struct {
	inner *ansi.Updater
}

func (u *Updater) Update(table string) builder.Updater { u.inner.Update(table); return u }
func (u *Updater) Set(v map[string]interface{}) builder.Updater {
	u.inner.Set(v)
	return u
}
func (u *Updater) Condition(e builder.Expression) builder.Updater {
	u.inner.Condition(e)
	return u
}
func (u *Updater) RawCondition(s string) builder.Updater { u.inner.RawCondition(s); return u }

func (u *Updater) Build() (string, []any, error) {
	sql, args, err := u.inner.Build()
	return rewritePlaceholders(sql), args, err
}

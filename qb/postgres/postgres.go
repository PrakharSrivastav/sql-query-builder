// Package postgres adapts the ANSI builder to PostgreSQL's $N
// placeholder syntax. Construction and chain methods stay in ANSI form
// (`?` placeholders); the rewrite happens once at Build time, walking
// the final SQL and substituting `$1`, `$2`, ... in order.
package postgres

import (
	"strconv"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/ansi"
	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
	"github.com/PrakharSrivastav/sql-query-builder/qb/core"
)

// NewPostgresBuilder returns a *core.SQL whose Reader/Inserter/Updater
// emit `$N` placeholders. Creater takes no values and reuses the ANSI
// implementation. Expression also stays ANSI: its fragments are merged
// into Reader/Updater Build, where the rewrite happens.
func NewPostgresBuilder() (*core.SQL, error) {
	return &core.SQL{
		Reader:        &Reader{inner: new(ansi.Reader)},
		Inserter:      &Inserter{inner: new(ansi.Inserter)},
		Updater:       &Updater{inner: new(ansi.Updater)},
		Creater:       new(ansi.Creater),
		NewExpression: func() builder.Expression { return new(ansi.Expression) },
	}, nil
}

// rewritePlaceholders substitutes each `?` in sql with `$1`, `$2`, ...
// in left-to-right order.
//
// The builder itself never emits `?` inside SQL string literals
// (values always go through the args slice), so this function does not
// parse strings or comments. Callers using RawCondition in the Postgres
// dialect must therefore avoid literal `?` characters inside the raw
// fragment — including inside string literals and Postgres JSON
// existence operators (`?`, `?|`, `?&`). Use the parameterized
// Condition path or stage the operator outside this builder.
func rewritePlaceholders(sql string) string {
	if strings.IndexByte(sql, '?') < 0 {
		return sql
	}
	var b strings.Builder
	b.Grow(len(sql) + 8)
	n := 1
	for {
		idx := strings.IndexByte(sql, '?')
		if idx < 0 {
			b.WriteString(sql)
			return b.String()
		}
		b.WriteString(sql[:idx])
		b.WriteByte('$')
		b.WriteString(strconv.Itoa(n))
		n++
		sql = sql[idx+1:]
	}
}

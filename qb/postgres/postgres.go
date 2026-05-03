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
// in left-to-right order. It does not parse strings or comments — this
// builder never emits `?` inside string literals because string values
// always go through the args slice.
func rewritePlaceholders(sql string) string {
	if !strings.ContainsRune(sql, '?') {
		return sql
	}
	var b strings.Builder
	b.Grow(len(sql) + 8)
	n := 1
	for _, r := range sql {
		if r == '?' {
			b.WriteByte('$')
			b.WriteString(strconv.Itoa(n))
			n++
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

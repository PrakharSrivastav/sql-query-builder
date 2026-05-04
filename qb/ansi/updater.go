package ansi

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Updater helps in creating parameterized UPDATE statements.
type Updater struct {
	sql       bytes.Buffer
	args      []any
	returning []string
	errs      []error
}

// Build returns the SQL, args (SET values first, then condition args)
// and any identifier-validation error.
func (u *Updater) Build() (string, []any, error) {
	if len(u.returning) > 0 {
		u.sql.WriteString(" RETURNING ")
		u.sql.WriteString(strings.Join(u.returning, seperator))
	}
	u.sql.WriteString(" ;")
	args := append([]any(nil), u.args...)
	return u.sql.String(), args, joinErrors(u.errs)
}

// Returning records columns for a RETURNING clause. Single `*` is
// accepted; everything else is validated as an identifier.
func (u *Updater) Returning(cols ...string) builder.Updater {
	for _, c := range cols {
		if c == "*" {
			continue
		}
		if err := validateIdentifier(c); err != nil {
			u.errs = append(u.errs, err)
		}
	}
	u.returning = append(u.returning, cols...)
	return u
}

// Condition merges an Expression's SQL fragment and args.
func (u *Updater) Condition(e builder.Expression) builder.Updater {
	sql, args, err := e.Express()
	u.sql.WriteString(sql)
	u.args = append(u.args, args...)
	if err != nil {
		u.errs = append(u.errs, err)
	}
	return u
}

// RawCondition appends a caller-supplied where clause verbatim. Caller
// is responsible for safety; use Condition with a Clause for untrusted
// input.
func (u *Updater) RawCondition(s string) builder.Updater {
	if s != "" {
		u.sql.WriteString(strings.Join([]string{" ", s, " "}, ""))
	}
	return u
}

// Set emits SET col1 = ?, col2 = ? and captures values into args in
// sorted column order.
func (u *Updater) Set(values map[string]interface{}) builder.Updater {
	u.sql.WriteString(" SET ")
	columnNames := make([]string, 0, len(values))
	for key := range values {
		columnNames = append(columnNames, key)
	}
	sort.Strings(columnNames)

	setClauses := make([]string, 0, len(columnNames))
	for _, name := range columnNames {
		if err := validateIdentifier(name); err != nil {
			u.errs = append(u.errs, err)
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s=?", name))
		u.args = append(u.args, values[name])
	}
	u.sql.WriteString(strings.Join(setClauses, seperator))
	return u
}

// Update sets the target table.
func (u *Updater) Update(table string) builder.Updater {
	u.sql.Reset()
	u.args = u.args[:0]
	u.returning = nil
	u.errs = nil
	if err := validateIdentifier(table); err != nil {
		u.errs = append(u.errs, err)
	}
	u.sql.WriteString(fmt.Sprintf("UPDATE %s", table))
	return u
}

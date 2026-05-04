package ansi

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

type conflictAction int

const (
	conflictNone conflictAction = iota
	conflictDoNothing
	conflictDoUpdate
)

// Inserter creates a parameterized INSERT statement.
// Generates: INSERT INTO table ( c1, c2 ) values (?, ?),(?, ?)...
type Inserter struct {
	sql       bytes.Buffer
	args      []any
	columns   []string
	returning []string
	errs      []error

	conflictAction  conflictAction
	conflictTargets []string
	conflictSet     map[string]interface{}
	conflictKeys    []string
}

// Build returns the SQL, the args (in row order, sorted by column name)
// and any identifier-validation error. Build is idempotent.
func (i *Inserter) Build() (string, []any, error) {
	sql := strings.TrimSuffix(i.sql.String(), ",")
	args := append([]any(nil), i.args...)

	switch i.conflictAction {
	case conflictDoNothing:
		sql += " ON CONFLICT (" + strings.Join(i.conflictTargets, seperator) + ") DO NOTHING"
	case conflictDoUpdate:
		setClauses := make([]string, 0, len(i.conflictKeys))
		for _, k := range i.conflictKeys {
			if exc, ok := i.conflictSet[k].(builder.Excluded); ok {
				setClauses = append(setClauses, fmt.Sprintf("%s = EXCLUDED.%s", k, exc.Col))
				continue
			}
			setClauses = append(setClauses, k+" = ?")
			args = append(args, i.conflictSet[k])
		}
		sql += " ON CONFLICT (" + strings.Join(i.conflictTargets, seperator) + ") DO UPDATE SET " + strings.Join(setClauses, seperator)
	}

	if len(i.returning) > 0 {
		sql += " RETURNING " + strings.Join(i.returning, seperator)
	}
	sql += ";"
	return sql, args, joinErrors(i.errs)
}

// Returning records columns for a RETURNING clause. Single `*` is
// accepted; everything else is validated as an identifier.
func (i *Inserter) Returning(cols ...string) builder.Inserter {
	for _, c := range cols {
		if c == "*" {
			continue
		}
		if err := validateIdentifier(c); err != nil {
			i.errs = append(i.errs, err)
		}
	}
	i.returning = append(i.returning, cols...)
	return i
}

// OnConflictDoNothing records targets for an ON CONFLICT DO NOTHING
// clause. Each target is validated as an identifier.
func (i *Inserter) OnConflictDoNothing(targets ...string) builder.Inserter {
	for _, t := range targets {
		if err := validateIdentifier(t); err != nil {
			i.errs = append(i.errs, err)
		}
	}
	i.conflictAction = conflictDoNothing
	i.conflictTargets = append(i.conflictTargets[:0], targets...)
	i.conflictSet = nil
	i.conflictKeys = nil
	return i
}

// OnConflictDoUpdate records targets and a set map for an ON CONFLICT
// DO UPDATE clause. Targets and set keys are validated; values are
// bound via placeholders unless wrapped in builder.Excluded.
func (i *Inserter) OnConflictDoUpdate(targets []string, set map[string]interface{}) builder.Inserter {
	for _, t := range targets {
		if err := validateIdentifier(t); err != nil {
			i.errs = append(i.errs, err)
		}
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := validateIdentifier(k); err != nil {
			i.errs = append(i.errs, err)
		}
		if exc, ok := set[k].(builder.Excluded); ok {
			if err := validateIdentifier(exc.Col); err != nil {
				i.errs = append(i.errs, err)
			}
		}
	}
	i.conflictAction = conflictDoUpdate
	i.conflictTargets = append(i.conflictTargets[:0], targets...)
	i.conflictSet = set
	i.conflictKeys = keys
	return i
}

// Columns sets the column list for the insert. Names are sorted so the
// Values map can be looked up deterministically.
func (i *Inserter) Columns(s []string) builder.Inserter {
	cols := append([]string(nil), s...)
	sort.Strings(cols)
	for _, c := range cols {
		if err := validateIdentifier(c); err != nil {
			i.errs = append(i.errs, err)
		}
	}
	i.columns = cols
	i.sql.WriteString(strings.Join(cols, seperator))
	i.sql.WriteString(" ) values ")
	return i
}

// Table sets the destination table name.
func (i *Inserter) Table(s string) builder.Inserter {
	i.sql.Reset()
	i.args = i.args[:0]
	i.columns = nil
	i.returning = nil
	i.errs = nil
	i.conflictAction = conflictNone
	i.conflictTargets = nil
	i.conflictSet = nil
	i.conflictKeys = nil
	if err := validateIdentifier(s); err != nil {
		i.errs = append(i.errs, err)
	}
	i.sql.WriteString(fmt.Sprintf("INSERT INTO %s ( ", s))
	return i
}

// Values appends one row of placeholders, in the order set by Columns,
// and captures the values into args.
func (i *Inserter) Values(v builder.Value) builder.Inserter {
	if len(i.columns) == 0 {
		i.errs = append(i.errs, fmt.Errorf("Values called before Columns"))
		return i
	}
	placeholders := make([]string, 0, len(i.columns))
	for _, col := range i.columns {
		placeholders = append(placeholders, "?")
		i.args = append(i.args, v[col])
	}
	i.sql.WriteString("(")
	i.sql.WriteString(strings.Join(placeholders, seperator))
	i.sql.WriteString("),")
	return i
}

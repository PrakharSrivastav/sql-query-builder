package ansi

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Inserter creates a parameterized INSERT statement.
// Generates: INSERT INTO table ( c1, c2 ) values (?, ?),(?, ?)...
type Inserter struct {
	sql     bytes.Buffer
	args    []any
	columns []string
	errs    []error
}

// Build returns the SQL, the args (in row order, sorted by column name)
// and any identifier-validation error.
func (i *Inserter) Build() (string, []any, error) {
	sql := strings.TrimSuffix(i.sql.String(), ",") + ";"
	args := append([]any(nil), i.args...)
	return sql, args, joinErrors(i.errs)
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
	i.errs = nil
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

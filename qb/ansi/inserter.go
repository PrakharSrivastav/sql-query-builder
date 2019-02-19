package ansi

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Inserter creates a insert sql statement
// INSERT INTO table () values (),(),()
type Inserter struct {
	sql bytes.Buffer
}

// Build yeilds the final sql statement
func (i *Inserter) Build() string {
	sql := i.sql.String()
	sql = strings.TrimSuffix(sql, ",")
	sql = sql + ";"
	return sql
}

// Columns defines the set of columns for which inserts will happen
func (i *Inserter) Columns(s []string) builder.Inserter {
	sort.Strings(s)
	i.sql.WriteString(strings.Join(s, seperator))
	i.sql.WriteString(" ) values ")
	return i
}

// Table sets the table name
func (i *Inserter) Table(s string) builder.Inserter {
	i.sql.Reset()
	i.sql.WriteString(fmt.Sprintf("INSERT INTO %s ( ", s))
	return i
}

// Values evaluate the insert values clause. Should be used multiple times to create a multi-insert clause.
// Check tests for usage
func (i *Inserter) Values(v builder.Value) builder.Inserter {
	fields := make([]string, 0, len(v))
	values := make([]string, 0, len(v))

	for item := range v {
		fields = append(fields, item)
	}
	sort.Strings(fields)
	for _, field := range fields {
		switch v[field].(type) {
		case string:
			values = append(values, fmt.Sprintf("'%s'", v[field].(string)))
		default:
			values = append(values, fmt.Sprintf("%v", v[field]))
		}
	}

	joinedValues := strings.Join(values, seperator)
	i.sql.WriteString(strings.Join([]string{"(", joinedValues, "),"}, ""))
	return i
}

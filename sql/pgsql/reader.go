package pgsql

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"
)

// Reader implements interface to create select clauses
type Reader struct {
	sql bytes.Buffer
}

// Select builds the select clause for sql
func (r *Reader) Select(s ...string) builder.Reader {
	r.sql.Reset()
	r.sql.WriteString("SELECT ")
	r.sql.WriteString(strings.Join(s, seperator))
	return r
}

// From adds the from clause in the sql
func (r *Reader) From(s ...string) builder.Reader {
	r.sql.WriteString(" FROM ")
	r.sql.WriteString(strings.Join(s, seperator))
	return r
}

// FromAlias helps to generate from clause similar to
// select field1, field2 from table1 as t1 , table2 as t2
func (r *Reader) FromAlias(alias ...builder.Alias) builder.Reader {
	r.sql.WriteString(" FROM ")
	str := make([]string, 0, cap(alias))
	for _, a := range alias {
		str = append(str, fmt.Sprintf("%s as %s", a.Name, a.Alias))
	}
	r.sql.WriteString(strings.Join(str, seperator))
	return r
}

// Build compiles the expression and generates a sql equivalent of sql
func (r *Reader) Build() string {
	r.sql.WriteString(" ;")
	return r.sql.String()
}

// Limit adds limit clause to the sql
func (r *Reader) Limit(i int) builder.Reader {
	r.sql.WriteString(fmt.Sprintf(" LIMIT %s", strconv.Itoa(i)))
	return r
}

// Offset adds offset clause to the sql
func (r *Reader) Offset(i int) builder.Reader {
	r.sql.WriteString(fmt.Sprintf(" OFFSET %s", strconv.Itoa(i)))
	return r
}

// OrderBy for the order by clause
func (r *Reader) OrderBy(s ...string) builder.Reader {
	r.sql.WriteString(" ORDER BY ")
	r.sql.WriteString(strings.Join(s, seperator))
	return r
}

// Condition to implement the where clause with Expressions
func (r *Reader) Condition(expression builder.Expression) builder.Reader {
	r.sql.WriteString(expression.Express())
	return r
}

// RawCondition to add where clause in string format
// Assumes that a well formatted where clause is provided.
// The input expression input should start with where
func (r *Reader) RawCondition(expression string) builder.Reader {
	r.sql.WriteString(expression)
	return r
}

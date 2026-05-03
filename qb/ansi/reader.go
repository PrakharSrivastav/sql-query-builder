package ansi

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Reader implements interface to create select clauses
type Reader struct {
	sql  bytes.Buffer
	args []any
	errs []error
}

// Select builds the select clause for sql
func (r *Reader) Select(s ...string) builder.Reader {
	r.sql.Reset()
	r.args = r.args[:0]
	r.errs = nil
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

// Build compiles the SQL, returns it with the captured args slice and
// any identifier-validation error.
func (r *Reader) Build() (string, []any, error) {
	r.sql.WriteString(" ;")
	args := append([]any(nil), r.args...)
	return r.sql.String(), args, joinErrors(r.errs)
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

// Condition merges an Expression's SQL fragment and args into the reader.
func (r *Reader) Condition(expression builder.Expression) builder.Reader {
	sql, args, err := expression.Express()
	r.sql.WriteString(sql)
	r.args = append(r.args, args...)
	if err != nil {
		r.errs = append(r.errs, err)
	}
	return r
}

// RawCondition appends a caller-supplied where clause verbatim. The
// caller is responsible for safety; use Condition with a Clause for
// untrusted input.
func (r *Reader) RawCondition(expression string) builder.Reader {
	r.sql.WriteString(expression)
	return r
}

// InnerJoin creates an inner join clause
func (r *Reader) InnerJoin(table string) builder.Reader {
	r.sql.WriteString(" INNER JOIN ")
	r.sql.WriteString(table)
	return r
}

// LeftJoin creates a Left Join clause
func (r *Reader) LeftJoin(table string) builder.Reader {
	r.sql.WriteString(" LEFT JOIN ")
	r.sql.WriteString(table)
	return r
}

// RightJoin creates a Right Join clause
func (r *Reader) RightJoin(table string) builder.Reader {
	r.sql.WriteString(" RIGHT JOIN ")
	r.sql.WriteString(table)
	return r
}

// On creates an on clause
func (r *Reader) On(condition string) builder.Reader {
	r.sql.WriteString(" ON ")
	r.sql.WriteString(condition)
	return r
}

// Having creates a having clause
func (r *Reader) Having(condition string) builder.Reader {
	r.sql.WriteString(" HAVING ")
	r.sql.WriteString(condition)
	return r
}

// GroupBy creates a group by clause on the input fields
func (r *Reader) GroupBy(fields []string) builder.Reader {
	r.sql.WriteString(" GROUP BY ")
	r.sql.WriteString(strings.Join(fields, seperator))
	return r
}

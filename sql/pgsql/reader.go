package pgsql

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"
)

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

func (r *Reader) FromAlias(alias ...builder.Alias) builder.Reader {
	r.sql.WriteString(" FROM ")
	str := make([]string, 0, cap(alias))
	for _, a := range alias {
		str = append(str, fmt.Sprintf("%s as %s", a.Name, a.Alias))
	}
	r.sql.WriteString(strings.Join(str, seperator))
	return r
}

func (r *Reader) Build() string {
	r.sql.WriteString(" ;")
	return r.sql.String()
}

func (r *Reader) Limit(i int) builder.Reader {
	r.sql.WriteString(fmt.Sprintf(" LIMIT %s", strconv.Itoa(i)))
	return r
}
func (r *Reader) Offset(i int) builder.Reader {
	r.sql.WriteString(fmt.Sprintf(" OFFSET %s", strconv.Itoa(i)))
	return r
}

func (r *Reader) OrderBy(s ...string) builder.Reader {
	r.sql.WriteString(" ORDER BY ")
	r.sql.WriteString(strings.Join(s, seperator))
	return r
}

func (r *Reader) Condition(expression builder.Expression) builder.Reader {
	r.sql.WriteString(expression.Express())
	return r
}

func (r *Reader) RawCondition(expression string) builder.Reader {
	r.sql.WriteString(expression)
	return r
}

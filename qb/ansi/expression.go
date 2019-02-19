package ansi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Expression is used to evaluate basic where clause statements
type Expression struct {
	b bytes.Buffer
}

// And evaluates an AND clause in sql
func (e *Expression) And(c builder.Clause) builder.Expression {
	switch c.Right.(type) {
	case string:
		e.b.WriteString(fmt.Sprintf("%s %s %s '%s' %s", " AND (", c.Left, c.Operator, c.Right, ")"))
	default:
		e.b.WriteString(fmt.Sprintf("%s %s %s %v %s", " AND (", c.Left, c.Operator, c.Right, ")"))
	}
	return e
}

// Express yeilds a complex where clause
func (e *Expression) Express() string {
	return e.b.String()
}

// In creates an In clause
func (e *Expression) In(field string, items ...string) builder.Expression {
	e.b.WriteString(fmt.Sprintf(" %s ( %s IN [%s]", "AND", field, strings.Join(items, seperator)))
	return e
}

// NotIn evaluates a not in clause
func (e *Expression) NotIn(field string, items ...string) builder.Expression {
	e.b.WriteString(fmt.Sprintf(" %s ( %s NOT IN [%s]", "AND", field, strings.Join(items, seperator)))
	return e
}

// Or evaluates an Or sql clause
func (e *Expression) Or(c builder.Clause) builder.Expression {
	switch c.Right.(type) {
	case string:
		e.b.WriteString(fmt.Sprintf("%s %s %s '%s' %s", " OR (", c.Left, c.Operator, c.Right, ")"))
	default:
		e.b.WriteString(fmt.Sprintf("%s %s %s %v %s", " OR (", c.Left, c.Operator, c.Right, ")"))
	}
	return e
}

// Where evaluates a where clause
func (e *Expression) Where(c builder.Clause) builder.Expression {
	e.b.Reset()
	switch c.Right.(type) {
	case string:
		e.b.WriteString(fmt.Sprintf("%s %s %s '%s' %s", " WHERE (", c.Left, c.Operator, c.Right, " )"))
	default:
		e.b.WriteString(fmt.Sprintf("%s %s %s %v %s", " WHERE (", c.Left, c.Operator, c.Right, " )"))
	}

	return e
}

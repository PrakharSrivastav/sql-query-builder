package ansi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Expression is used to evaluate basic where clause statements
type Expression struct {
	b    bytes.Buffer
	args []any
	errs []error
}

// And evaluates an AND clause in sql
func (e *Expression) And(c builder.Clause) builder.Expression {
	e.appendClause(" AND", c)
	return e
}

// Or evaluates an Or sql clause
func (e *Expression) Or(c builder.Clause) builder.Expression {
	e.appendClause(" OR", c)
	return e
}

// Where evaluates a where clause
func (e *Expression) Where(c builder.Clause) builder.Expression {
	e.b.Reset()
	e.args = e.args[:0]
	e.errs = nil
	e.appendClause(" WHERE", c)
	return e
}

// In creates an IN clause joined to the existing expression with AND
func (e *Expression) In(field string, items ...any) builder.Expression {
	e.appendInClause(field, "IN", items)
	return e
}

// NotIn creates a NOT IN clause joined to the existing expression with AND
func (e *Expression) NotIn(field string, items ...any) builder.Expression {
	e.appendInClause(field, "NOT IN", items)
	return e
}

// Express yields the where SQL fragment, captured args, and any
// identifier-validation error.
func (e *Expression) Express() (string, []any, error) {
	args := append([]any(nil), e.args...)
	return e.b.String(), args, joinErrors(e.errs)
}

func (e *Expression) appendClause(keyword string, c builder.Clause) {
	if err := validateIdentifier(c.Left); err != nil {
		e.errs = append(e.errs, err)
		return
	}
	e.b.WriteString(fmt.Sprintf("%s (%s %s ?)", keyword, c.Left, c.Operator))
	e.args = append(e.args, c.Right)
}

func (e *Expression) appendInClause(field, keyword string, items []any) {
	if err := validateIdentifier(field); err != nil {
		e.errs = append(e.errs, err)
		return
	}
	if len(items) == 0 {
		e.errs = append(e.errs, fmt.Errorf("%s clause for %q requires at least one item", keyword, field))
		return
	}
	placeholders := strings.Repeat("?, ", len(items))
	placeholders = strings.TrimSuffix(placeholders, ", ")
	e.b.WriteString(fmt.Sprintf(" AND (%s %s (%s))", field, keyword, placeholders))
	e.args = append(e.args, items...)
}

func joinErrors(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}
	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}
	return fmt.Errorf("%s", strings.Join(msgs, "; "))
}

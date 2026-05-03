package ansi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Creater helps in creating a CREATE TABLE command.
type Creater struct {
	sql  bytes.Buffer
	errs []error
}

// Build yields the final SQL. CREATE TABLE has no values, so args is
// always empty; error captures identifier-validation failures.
func (c *Creater) Build() (string, []any, error) {
	return c.sql.String(), nil, joinErrors(c.errs)
}

// SetColumns defines column declarations. Column Name is validated as
// an identifier; Datatype and Constraint are caller-trusted SQL
// fragments (they may contain parens, length specifiers, defaults, etc).
func (c *Creater) SetColumns(c1 []builder.Columns) builder.Creater {
	columnDefs := make([]string, 0, len(c1))
	for _, item := range c1 {
		if err := validateIdentifier(item.Name); err != nil {
			c.errs = append(c.errs, err)
			continue
		}
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s %s", item.Name, item.Datatype, item.Constraint))
	}
	columns := strings.Join(columnDefs, seperator)
	c.sql.WriteString(fmt.Sprintf("%s %s %s", "(", columns, ");"))
	return c
}

// Table sets the table name.
func (c *Creater) Table(s string) builder.Creater {
	c.sql.Reset()
	c.errs = nil
	if err := validateIdentifier(s); err != nil {
		c.errs = append(c.errs, err)
	}
	c.sql.WriteString(fmt.Sprintf("CREATE TABLE %s ", s))
	return c
}

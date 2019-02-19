package ansi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
)

// Creater helps in creating a CREATE TABLE command
type Creater struct {
	sql bytes.Buffer
}

// Build yeilds the final sql statement
func (c *Creater) Build() string {
	return c.sql.String()
}

// SetColumns defines the column definition
func (c *Creater) SetColumns(c1 []builder.Columns) builder.Creater {

	columnDefs := make([]string, 0, len(c1))

	for _, item := range c1 {
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s %s", item.Name, item.Datatype, item.Constraint))
	}

	columns := strings.Join(columnDefs, seperator)
	c.sql.WriteString(fmt.Sprintf("%s %s %s", "(", columns, ");"))

	return c
}

// Table sets the table name
func (c *Creater) Table(s string) builder.Creater {
	c.sql.Reset()
	c.sql.WriteString(fmt.Sprintf("CREATE TABLE %s ", s))
	return c
}

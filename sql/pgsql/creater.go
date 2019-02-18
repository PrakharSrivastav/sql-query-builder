package pgsql

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"
)

type Creater struct {
	sql bytes.Buffer
}

func (c *Creater) Build() string {
	return c.sql.String()
}

func (c *Creater) SetColumns(c1 []builder.Columns) builder.Creater {

	columnDefs := make([]string, 0, len(c1))

	for _, item := range c1 {
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s %s", item.Name, item.Datatype, item.Constraint))
	}

	columns := strings.Join(columnDefs, seperator)
	c.sql.WriteString(fmt.Sprintf("%s %s %s", "(", columns, ");"))

	return c
}

func (c *Creater) Table(s string) builder.Creater {
	c.sql.Reset()
	c.sql.WriteString(fmt.Sprintf("CREATE TABLE %s ", s))
	return c
}

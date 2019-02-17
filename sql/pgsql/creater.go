package pgsql

import (
	"errors"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"
)

type Creater struct{}

func (c *Creater) Build() string {
	panic(errors.New("*Creater.Build not implemented"))
}

func (c *Creater) SetColumns(c1 ...builder.Columns) builder.Creater {
	panic(errors.New("*Creater.SetColumns not implemented"))
}

func (c *Creater) Table(s string) builder.Creater {
	panic(errors.New("*Creater.Table not implemented"))
}

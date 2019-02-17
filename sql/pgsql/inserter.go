package pgsql

import (
	"errors"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"
)

type Inserter struct{}

func (i *Inserter) Build() string {
	panic(errors.New("*Inserter.Build not implemented"))
}

func (i *Inserter) Columns(s ...string) builder.Inserter {
	panic(errors.New("*Inserter.Columns not implemented"))
}

func (i *Inserter) Table(s string) builder.Inserter {
	panic(errors.New("*Inserter.Table not implemented"))
}

func (i *Inserter) Values(v ...builder.Value) builder.Inserter {
	panic(errors.New("*Inserter.Values not implemented"))
}

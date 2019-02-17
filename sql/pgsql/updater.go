package pgsql

import (
	"errors"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"
)

type Updater struct{}

func (u *Updater) And(c builder.Clause) builder.Expression {
	panic(errors.New("*Updater.And not implemented"))
}

func (u *Updater) Build() string {
	panic(errors.New("*Updater.Build not implemented"))
}

func (u *Updater) Express() string {
	panic(errors.New("*Updater.Express not implemented"))
}

func (u *Updater) In(field string, items ...string) builder.Expression {
	panic(errors.New("*Updater.In not implemented"))
}

func (u *Updater) NotIn(field string, items ...string) builder.Expression {
	panic(errors.New("*Updater.NotIn not implemented"))
}

func (u *Updater) Or(c builder.Clause) builder.Expression {
	panic(errors.New("*Updater.Or not implemented"))
}

func (u *Updater) Set(m map[string]string) builder.Updater {
	panic(errors.New("*Updater.Set not implemented"))
}

func (u *Updater) Update(s string) builder.Updater {
	panic(errors.New("*Updater.Update not implemented"))
}

func (u *Updater) Where(c builder.Clause) builder.Expression {
	panic(errors.New("*Updater.Where not implemented"))
}

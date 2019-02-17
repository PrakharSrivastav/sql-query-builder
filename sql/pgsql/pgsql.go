package pgsql

import (
	"github.com/PrakharSrivastav/sql-query-builder/sql/core"
)

const (
	seperator = ", "
)

func NewPgSQLBuilder() (*core.SQL, error) {
	return &core.SQL{
		Creater:  new(Creater),
		Inserter: new(Inserter),
		Reader:   new(Reader),
		Updater:  new(Updater),
	}, nil
}

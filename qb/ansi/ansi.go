package ansi

import (
	"github.com/PrakharSrivastav/sql-query-builder/qb/core"
)

const (
	seperator = ", "
)

// NewANSIBuilder creates a builder that conforms to ANSI
func NewANSIBuilder() (*core.SQL, error) {
	return &core.SQL{
		Creater:  new(Creater),
		Inserter: new(Inserter),
		Reader:   new(Reader),
		Updater:  new(Updater),
	}, nil
}

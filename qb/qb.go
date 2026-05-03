/*
Package SQL helps to generate sql queries in different dialects.
This package can be best used with the scenarios where the structure of the domains models is unknown beforehand.
*/
package qb

import (
	"fmt"
	"sync"

	"github.com/PrakharSrivastav/sql-query-builder/qb/ansi"
	"github.com/PrakharSrivastav/sql-query-builder/qb/core"
	"github.com/PrakharSrivastav/sql-query-builder/qb/postgres"
)

var once sync.Once
var sql *core.SQL

// NewQueryBuilder takes in a dialect and returns QueryBuilder for a specific dialect
func NewQueryBuilder(driver int) (*core.SQL, error) {
	return dbFactory(driver)
}

// NewSingletonQueryBuilder returns a singleton querybuilder instance.
// Once the dialect is chosen, it can not be modified to another dialect.
// Prefer this if your application only connects to once database type
func NewSingletonQueryBuilder(driver int) (*core.SQL, error) {
	var err error
	once.Do(func() {
		sql, err = dbFactory(driver)
	})
	return sql, err
}

func dbFactory(driver int) (*core.SQL, error) {
	switch driver {
	case core.ANSI, core.MYSQL, core.SQLITE:
		// MySQL and SQLite both accept ANSI's `?` placeholder and the
		// validated identifier shape; ANSI output is portable to them.
		return ansi.NewANSIBuilder()
	case core.PGSQL:
		return postgres.NewPostgresBuilder()
	default:
		return nil, fmt.Errorf("qb: unsupported dialect %d (expected one of core.ANSI, core.PGSQL, core.MYSQL, core.SQLITE)", driver)
	}
}

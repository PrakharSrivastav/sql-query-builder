package core

import "github.com/PrakharSrivastav/sql-query-builder/sql/builder"

const (
	_ = iota
	// PGSQL adheres to Postgres dialect
	PGSQL
	// MYSQL adheres to MySql dialect
	MYSQL
	// SQLITE adheres to sqlite dialect
	SQLITE
	// MONGO adheres to MongoDB dialect
	MONGO
)

// SQL is wrapper for different driver implementations
type SQL struct {
	Reader   builder.Reader
	Inserter builder.Inserter
	Updater  builder.Updater
	Creater  builder.Creater
}

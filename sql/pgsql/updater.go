package pgsql

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"
)

// Updater helps in creating update sql queries
type Updater struct {
	sql bytes.Buffer
}

// Build returns the compiled update statement
func (u *Updater) Build() string {
	u.sql.WriteString(" ;")
	return u.sql.String()
}

// Condition accepts input of type buider.Expressiong to evaluate a where clause
func (u *Updater) Condition(e builder.Expression) builder.Updater {
	u.sql.WriteString(e.Express())
	return u
}

// RawCondition to add where clause in string format
// Assumes that a well formatted where clause is provided.
// The input expression input should start with where
func (u *Updater) RawCondition(s string) builder.Updater {
	if s != "" {
		u.sql.WriteString(strings.Join([]string{" ", s, " "}, ""))
	}
	return u
}

// Set sets the columns values for update
// Eg Set column1 = 'value1' , column2 = 'value2' , column3 = 3.21
func (u *Updater) Set(values map[string]interface{}) builder.Updater {
	u.sql.WriteString(" SET ")
	setClause := make([]string, 0, len(values))
	for key, value := range values {
		switch value.(type) {
		case string:
			setClause = append(setClause, fmt.Sprintf("%s = '%s'", key, value))
		default:
			setClause = append(setClause, fmt.Sprintf("%s = %s", key, value))
		}
	}
	u.sql.WriteString(strings.Join(setClause, seperator))
	return u
}

// Update clause takes the table name to prepare the correct update clause
// UPDATE table <table>
func (u *Updater) Update(table string) builder.Updater {
	u.sql.Reset()
	u.sql.WriteString("UPDATE table ")
	u.sql.WriteString(table)
	return u
}

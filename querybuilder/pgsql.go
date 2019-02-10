package querybuilder

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

/*
PgsqlQB generates sql query in postgres dialect
*/
type PgsqlQB struct {
	table     string
	idPrimary bool
}

// Create generates a create table statement
func (p *PgsqlQB) Create(m map[string]interface{}) (string, error) {
	if p.table == "" {
		return "", errors.New("table name required")
	}
	var bf bytes.Buffer
	bf.WriteString("CREATE TABLE ")
	bf.WriteString(p.table)
	bf.WriteString(" (")

	// This is for the consistency of the test cases. The maps in go are not ordered,
	// therefore we write the keys to an array, order it and then use it so that it always gives the consistent result
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, k := range keys {
		bf.WriteString(fmt.Sprintf("%s %s ,", k, getDatatype(m[k].(string))))
	}

	if p.idPrimary {
		bf.WriteString("id UUID primary key);")
	} else {
		bf.WriteString("id UUID);")
	}
	return bf.String(), nil
}

func (p *PgsqlQB) Get(columns []string, where map[string]interface{}, limit int, offset int) (string, error) {
	if p.table == "" {
		return "", errors.New("No table name provided")
	}
	if len(columns) == 0 {
		return "", errors.New("No columns provided")
	}

	var bf bytes.Buffer
	bf.WriteString("SELECT ")
	bf.WriteString(strings.Join(columns, ","))
	bf.WriteString(" FROM ")
	bf.WriteString(p.table)
	keys := make([]string, 0, len(where))
	for k := range where {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if len(keys) > 0 {
		bf.WriteString(" WHERE ")
	}
	whereClause := make([]string, 0, len(where))
	for _, key := range keys {
		if reflect.TypeOf(where[key]).String() == "string" {
			whereClause = append(whereClause, fmt.Sprintf("%s = '%s'", key, where[key]))
			continue
		}
		whereClause = append(whereClause, fmt.Sprintf("%s = %v", key, where[key]))
	}
	bf.WriteString(strings.Join(whereClause, " AND "))

	if limit != 0 {
		bf.WriteString(fmt.Sprintf(" LIMIT %d ", limit))
	}
	if offset != 0 {
		bf.WriteString(fmt.Sprintf(" OFFSET %d ", offset))
	}
	bf.WriteString(";")
	return bf.String(), nil
}

func (p *PgsqlQB) Insert(columns []string, data []map[string]interface{}) (string, error) {
	if b.table == "" {
		return "", errors.New("No table name provided")
	}
	if len(columns) == 0 {
		return "", errors.New("No columns provided for update")
	}
	if len(data) == 0 {
		return "", errors.New("No data provided for insert")
	}
	// prepare insert
	// prepare values
	sort.Strings(columns)
	fmt.Println(columns)
	var bf bytes.Buffer
	bf.WriteString("INSERT INTO ")
	bf.WriteString(b.table)
	bf.WriteString(" (")
	bf.WriteString(strings.Join(columns, ","))
	bf.WriteString(") VALUES ")
	values := []string{}
	for _, row := range data {
		var tempBuffer bytes.Buffer
		tempBuffer.WriteString("(")
		value := make([]string, 0, len(columns))
		for _, column := range columns {
			if reflect.TypeOf(row[column]).String() == "string" {
				value = append(value, fmt.Sprintf("'%s'", row[column]))
				continue
			}
			value = append(value, fmt.Sprintf("%v", row[column]))
		}
		tempBuffer.WriteString(strings.Join(value, ","))
		tempBuffer.WriteString(")")
		values = append(values, tempBuffer.String())
		tempBuffer.Reset()
	}
	bf.WriteString(strings.Join(values, ","))
	bf.WriteString(";")
	return bf.String(), nil
}

func (p *PgsqlQB) Update(columns, where map[string]interface{}) (string, error) {
	if p.table == "" {
		return "", errors.New("No table name provided")
	}
	if len(columns) == 0 {
		return "", errors.New("No columns provided for update")
	}
	var bf bytes.Buffer
	bf.WriteString("UPDATE ")
	bf.WriteString(p.table)
	bf.WriteString(" SET ")

	// Set columns
	keys := make([]string, 0, len(columns))
	for key := range columns {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	cols := make([]string, 0, len(columns))
	for _, key := range keys {
		if reflect.TypeOf(columns[key]).String() == "string" {
			cols = append(cols, fmt.Sprintf("%s='%s'", key, columns[key]))
			continue
		}
		cols = append(cols, fmt.Sprintf("%s=%v", key, columns[key]))
	}
	bf.WriteString(strings.Join(cols, ","))

	whereLen := len(where)
	if whereLen > 0 {
		bf.WriteString(" WHERE ")
		// set where clause
		whereKeys := make([]string, 0, whereLen)
		for key := range where {
			whereKeys = append(whereKeys, key)
		}
		sort.Strings(whereKeys)

		whereClause := make([]string, 0, whereLen)
		for _, key := range whereKeys {
			if reflect.TypeOf(where[key]).String() == "string" {
				whereClause = append(whereClause, fmt.Sprintf("%s='%s'", key, where[key]))
				continue
			}
			whereClause = append(whereClause, fmt.Sprintf("%s=%v", key, where[key]))
		}
		bf.WriteString(strings.Join(whereClause, " AND "))
	}

	bf.WriteString(" ;")

	return bf.String(), nil
}

func getDatatype(from string) string {
	switch from {
	case "float":
		return "float"
	case "keyword":
		return "varchar(200) "
	default:
		return "text"
	}
}

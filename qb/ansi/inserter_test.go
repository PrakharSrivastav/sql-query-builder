package ansi

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"

	"github.com/stretchr/testify/assert"
)

func TestSimpleInsert(t *testing.T) {
	t.Parallel()
	i := new(Inserter)

	sql, args, err := i.Table("xyz").Columns([]string{"field1", "field2"}).Values(builder.Value{
		"field1": 123,
		"field2": "123",
	}).Build()

	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO xyz ( field1, field2 ) values (?, ?);", sql)
	assert.Equal(t, []any{123, "123"}, args)
}

func TestMultiInsert(t *testing.T) {
	t.Parallel()
	i := new(Inserter)

	tableName := "xyz"
	columns := []string{"field1", "field2"}

	data := []map[string]interface{}{
		{"field1": 345, "field2": "asdf"},
		{"field1": 123, "field2": "qwer"},
		{"field1": 234, "field2": "zxcv"},
	}
	stmt := i.Table(tableName).Columns(columns)
	for _, item := range data {
		stmt.Values(item)
	}
	sql, args, err := stmt.Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO xyz ( field1, field2 ) values (?, ?),(?, ?),(?, ?);", sql)
	assert.Equal(t, []any{345, "asdf", 123, "qwer", 234, "zxcv"}, args)
}

func TestInsertValuesBeforeColumnsErrors(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	_, _, err := i.Table("xyz").Values(builder.Value{"field1": 1}).Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Values called before Columns")
}

func TestMultiInsertFloat(t *testing.T) {
	t.Parallel()
	i := new(Inserter)

	data := []map[string]interface{}{
		{"field1": 34.5, "field2": "asdf"},
		{"field1": 12.3, "field2": "qwer"},
	}
	stmt := i.Table("xyz").Columns([]string{"field1", "field2"})
	for _, item := range data {
		stmt.Values(item)
	}
	sql, args, err := stmt.Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO xyz ( field1, field2 ) values (?, ?),(?, ?);", sql)
	assert.Equal(t, []any{34.5, "asdf", 12.3, "qwer"}, args)
}

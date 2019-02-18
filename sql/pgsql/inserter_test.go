package pgsql

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"

	"github.com/stretchr/testify/assert"
)

func TestSimpleInsert(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	assert.NotNil(t, i)

	sql := i.Table("xyz").Columns([]string{"field1", "field2"}).Values(builder.Value{
		"field1": 123,
		"field2": "123",
	}).Build()

	assert.Equal(t, sql, "INSERT INTO xyz ( field1, field2 ) values (123, '123');")
}

func TestMultiInsert(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	assert.NotNil(t, i)

	tableName := "xyz"
	columns := []string{"field1", "field2"}

	data := []map[string]interface{}{
		{"field1": 345, "field2": "asdf"},
		{"field1": 123, "field2": "qwer"},
		{"field1": 234, "field2": "zxcv"},
		{"field1": 456, "field2": "fgjh"},
		{"field1": 567, "field2": "ghjk"},
		{"field1": 678, "field2": "kjgg"},
	}
	sql := i.Table(tableName).Columns(columns)
	for _, item := range data {
		sql.Values(item)
	}
	builtSQL := sql.Build()
	assert.NotNil(t, builtSQL)
	assert.Equal(t, "INSERT INTO xyz ( field1, field2 ) values (345, 'asdf'),(123, 'qwer'),(234, 'zxcv'),(456, 'fgjh'),(567, 'ghjk'),(678, 'kjgg');", builtSQL)
}

func TestMultiInsertFloat(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	assert.NotNil(t, i)

	tableName := "xyz"
	columns := []string{"field1", "field2"}

	data := []map[string]interface{}{
		{"field1": 34.5, "field2": "asdf"},
		{"field1": 12.3, "field2": "qwer"},
		{"field1": 23.4, "field2": "zxcv"},
		{"field1": 45.6, "field2": "fgjh"},
		{"field1": 56.7, "field2": "ghjk"},
		{"field1": 67.8, "field2": "kjgg"},
	}
	sql := i.Table(tableName).Columns(columns)
	for _, item := range data {
		sql.Values(item)
	}
	builtSQL := sql.Build()
	assert.NotNil(t, builtSQL)
	assert.Equal(t, "INSERT INTO xyz ( field1, field2 ) values (34.5, 'asdf'),(12.3, 'qwer'),(23.4, 'zxcv'),(45.6, 'fgjh'),(56.7, 'ghjk'),(67.8, 'kjgg');", builtSQL)
}

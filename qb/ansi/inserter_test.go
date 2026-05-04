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

func TestInsertReturning(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	sql, args, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		Returning("id", "created_at").Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO users ( name ) values (?) RETURNING id, created_at;", sql)
	assert.Equal(t, []any{"alice"}, args)
}

func TestInsertReturningStar(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	sql, _, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		Returning("*").Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO users ( name ) values (?) RETURNING *;", sql)
}

func TestInsertReturningBadIdentifierErrors(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	_, _, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		Returning("id; DROP TABLE users;--").Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

func TestInsertOnConflictDoNothing(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	sql, args, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoNothing("id").Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO users ( name ) values (?) ON CONFLICT (id) DO NOTHING;", sql)
	assert.Equal(t, []any{"alice"}, args)
}

func TestInsertOnConflictDoUpdate(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	sql, args, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{
			"name":       "alice-updated",
			"updated_at": 12345,
		}).Build()
	assert.NoError(t, err)
	// Set keys are sorted.
	assert.Equal(t, "INSERT INTO users ( name ) values (?) ON CONFLICT (id) DO UPDATE SET name = ?, updated_at = ?;", sql)
	assert.Equal(t, []any{"alice", "alice-updated", 12345}, args)
}

func TestInsertOnConflictDoUpdate_WithExcluded(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	sql, args, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{
			"name":       builder.Excluded{Col: "name"},
			"updated_at": 12345,
		}).Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO users ( name ) values (?) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, updated_at = ?;", sql)
	assert.Equal(t, []any{"alice", 12345}, args)
}

func TestInsertOnConflictDoUpdate_WithReturning(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	sql, args, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{
			"name": builder.Excluded{Col: "name"},
		}).
		Returning("id").Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO users ( name ) values (?) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name RETURNING id;", sql)
	assert.Equal(t, []any{"alice"}, args)
}

func TestInsertOnConflict_BadTargetIdentifierErrors(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	_, _, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoNothing("id; DROP TABLE users;--").Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

func TestInsertOnConflictDoUpdate_BadSetKeyErrors(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	_, _, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{
			"name; DROP": "x",
		}).Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

func TestInsertOnConflictDoUpdate_BadExcludedColErrors(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	_, _, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{
			"name": builder.Excluded{Col: "name; DROP"},
		}).Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

func TestInsertBuildIsIdempotent(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	stmt := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{"name": "x"})

	sql1, args1, err1 := stmt.Build()
	sql2, args2, err2 := stmt.Build()
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, sql1, sql2)
	assert.Equal(t, args1, args2)
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

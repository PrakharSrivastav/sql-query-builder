package ansi

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"

	"github.com/stretchr/testify/assert"
)

func TestSimpleUpdateClause(t *testing.T) {
	t.Parallel()
	u := new(Updater)

	sql, args, err := u.Update("xyz").Build()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "UPDATE xyz ;", sql)

	columns := map[string]interface{}{
		"field1": "value1",
		"field2": 123,
		"field3": 321.123,
	}
	sql, args, err = u.Update("xyz").Set(columns).Build()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE xyz SET field1=?, field2=?, field3=? ;", sql)
	assert.Equal(t, []any{"value1", 123, 321.123}, args)
}

func TestSUpdateWithRawCondition(t *testing.T) {
	t.Parallel()
	u := new(Updater)

	columns := map[string]interface{}{
		"field1": "value1",
		"field2": 123,
		"field3": 321.123,
	}
	sql, args, err := u.Update("xyz").Set(columns).RawCondition("WHERE field1='another value'").Build()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE xyz SET field1=?, field2=?, field3=? WHERE field1='another value'  ;", sql)
	assert.Equal(t, []any{"value1", 123, 321.123}, args)
}

func TestUpdateReturning(t *testing.T) {
	t.Parallel()
	u := new(Updater)
	sql, args, err := u.Update("users").
		Set(map[string]interface{}{"name": "alice"}).
		Returning("id", "updated_at").Build()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE users SET name=? RETURNING id, updated_at ;", sql)
	assert.Equal(t, []any{"alice"}, args)
}

func TestUpdateReturningBadIdentifierErrors(t *testing.T) {
	t.Parallel()
	u := new(Updater)
	_, _, err := u.Update("users").
		Set(map[string]interface{}{"name": "alice"}).
		Returning("id; DROP TABLE users;--").Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

func TestUpdateWithCondition(t *testing.T) {
	t.Parallel()
	u := new(Updater)
	expression := new(Expression)

	frag, exprArgs, err := expression.
		Where(builder.Clause{Left: "field1", Operator: "=", Right: "value1"}).
		And(builder.Clause{Left: "field2", Operator: ">", Right: 12.2}).
		Express()
	assert.NoError(t, err)
	assert.Equal(t, " WHERE (field1 = ?) AND (field2 > ?)", frag)
	assert.Equal(t, []any{"value1", 12.2}, exprArgs)

	sqlExpression := expression.
		Where(builder.Clause{Left: "field1", Operator: "=", Right: "value1"}).
		And(builder.Clause{Left: "field2", Operator: ">", Right: 12.2})

	columns := map[string]interface{}{
		"field1": "value1",
		"field2": 123,
		"field3": 321.123,
	}
	sql, args, err := u.Update("xyz").Set(columns).Condition(sqlExpression).Build()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE xyz SET field1=?, field2=?, field3=? WHERE (field1 = ?) AND (field2 > ?) ;", sql)
	assert.Equal(t, []any{"value1", 123, 321.123, "value1", 12.2}, args)
}

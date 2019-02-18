package pgsql

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/sql/builder"

	"github.com/stretchr/testify/assert"
)

func TestSimpleUpdateClause(t *testing.T) {
	t.Parallel()
	u := new(Updater)

	assert.NotNil(t, u)
	sql := u.Update("xyz").Build()
	assert.NotNil(t, sql)
	assert.Equal(t, "UPDATE table xyz ;", sql)

	columns := make(map[string]interface{})
	columns["field1"] = "value1"
	columns["field2"] = 123
	columns["field3"] = 321.123
	sql = u.Update("xyz").Set(columns).Build()

	assert.NotNil(t, sql)
	assert.Equal(t, "UPDATE table xyz SET field1='value1', field2=123, field3=321.123 ;", sql)
}

func TestSUpdateWithRawCondition(t *testing.T) {
	t.Parallel()
	u := new(Updater)

	columns := make(map[string]interface{})
	columns["field1"] = "value1"
	columns["field2"] = 123
	columns["field3"] = 321.123
	sql := u.Update("xyz").Set(columns).RawCondition("WHERE field1='another value'").Build()

	assert.NotNil(t, sql)
	assert.Equal(t, "UPDATE table xyz SET field1='value1', field2=123, field3=321.123 WHERE field1='another value'  ;", sql)
}

func TestUpdateWithCondition(t *testing.T) {
	t.Parallel()
	u := new(Updater)
	assert.NotNil(t, u)

	expression := new(Expression)
	assert.NotNil(t, expression)

	sql := expression.
		Where(builder.Clause{Left: "field1", Operator: "=", Right: "value1"}).
		And(builder.Clause{Left: "field2", Operator: ">", Right: 12.2}).
		Express()

	assert.NotEmpty(t, sql)
	assert.Equal(t, " WHERE ( field1 = 'value1'  ) AND ( field2 > 12.2 )", sql)

	sqlExpression := expression.
		Where(builder.Clause{Left: "field1", Operator: "=", Right: "value1"}).
		And(builder.Clause{Left: "field2", Operator: ">", Right: 12.2})

	columns := make(map[string]interface{})
	columns["field1"] = "value1"
	columns["field2"] = 123
	columns["field3"] = 321.123
	updateSQL := u.Update("xyz").Set(columns).Condition(sqlExpression).Build()

	assert.NotNil(t, updateSQL)
	assert.NotEmpty(t, updateSQL)

	assert.Equal(t, "UPDATE table xyz SET field1='value1', field2=123, field3=321.123 WHERE ( field1 = 'value1'  ) AND ( field2 > 12.2 ) ;", updateSQL)
}

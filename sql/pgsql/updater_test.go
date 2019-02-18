package pgsql

import (
	"testing"

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
	assert.Equal(t, "UPDATE table xyz set field1='value1', field2=123, field3=321.123 ;", sql)
}

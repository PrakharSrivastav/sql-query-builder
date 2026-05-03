package ansi

import (
	"strings"
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"

	"github.com/stretchr/testify/assert"
)

// TestInjection_ValueInWhereIsParameterized proves that an attacker
// payload supplied as the right-hand side of a Where clause never
// reaches the SQL string — it goes into args verbatim.
func TestInjection_ValueInWhereIsParameterized(t *testing.T) {
	t.Parallel()
	payload := "anything' OR '1'='1"
	r := new(Reader)
	expr := new(Expression)
	expr.Where(builder.Clause{Left: "name", Operator: "=", Right: payload})
	sql, args, err := r.Select("id").From("users").Condition(expr).Build()

	assert.NoError(t, err)
	assert.NotContains(t, sql, "OR")
	assert.NotContains(t, sql, "'")
	assert.Contains(t, sql, "?")
	assert.Equal(t, []any{payload}, args)
}

// TestInjection_ValueInSetIsParameterized proves UPDATE ... SET values
// are bound, not inlined.
func TestInjection_ValueInSetIsParameterized(t *testing.T) {
	t.Parallel()
	payload := "x'; DROP TABLE users;--"
	u := new(Updater)
	sql, args, err := u.Update("users").Set(map[string]interface{}{"name": payload}).Build()

	assert.NoError(t, err)
	assert.NotContains(t, sql, "DROP")
	assert.NotContains(t, sql, "'")
	assert.Equal(t, "UPDATE users SET name=? ;", sql)
	assert.Equal(t, []any{payload}, args)
}

// TestInjection_ValueInInsertIsParameterized proves INSERT values are
// bound, not inlined.
func TestInjection_ValueInInsertIsParameterized(t *testing.T) {
	t.Parallel()
	payload := "'); DELETE FROM users;--"
	i := new(Inserter)
	sql, args, err := i.Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": payload}).Build()

	assert.NoError(t, err)
	assert.NotContains(t, sql, "DELETE")
	assert.NotContains(t, sql, "'")
	assert.Equal(t, "INSERT INTO users ( name ) values (?);", sql)
	assert.Equal(t, []any{payload}, args)
}

// TestInjection_ValueInInClauseIsParameterized proves IN list items are
// bound, not inlined.
func TestInjection_ValueInInClauseIsParameterized(t *testing.T) {
	t.Parallel()
	payloads := []any{"a", "b' OR 1=1--"}
	expr := new(Expression)
	sql, args, err := expr.
		Where(builder.Clause{Left: "id", Operator: "=", Right: 1}).
		In("name", payloads...).Express()

	assert.NoError(t, err)
	assert.NotContains(t, sql, "OR")
	assert.NotContains(t, sql, "'")
	assert.Equal(t, " WHERE (id = ?) AND (name IN (?, ?))", sql)
	assert.Equal(t, []any{1, "a", "b' OR 1=1--"}, args)
}

// TestInjection_BadIdentifierInWhereLeftRejected proves that a column
// name slot — which gets inlined into SQL — refuses anything not a
// plain identifier.
func TestInjection_BadIdentifierInWhereLeftRejected(t *testing.T) {
	t.Parallel()
	expr := new(Expression)
	expr.Where(builder.Clause{Left: "name; DROP TABLE users;--", Operator: "=", Right: "x"})
	_, _, err := expr.Express()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

// TestInjection_BadIdentifierInTableNameRejected proves the table name
// slot rejects injection attempts.
func TestInjection_BadIdentifierInTableNameRejected(t *testing.T) {
	t.Parallel()
	u := new(Updater)
	_, _, err := u.Update("users; DROP TABLE users;--").
		Set(map[string]interface{}{"name": "x"}).Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

// TestInjection_BadIdentifierInSetKeyRejected proves the SET column
// slot rejects injection attempts.
func TestInjection_BadIdentifierInSetKeyRejected(t *testing.T) {
	t.Parallel()
	u := new(Updater)
	_, _, err := u.Update("users").
		Set(map[string]interface{}{"name=1; DROP TABLE users;--": "x"}).Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

// TestInjection_BadIdentifierInInsertColumnRejected proves the Insert
// column slot rejects injection attempts.
func TestInjection_BadIdentifierInInsertColumnRejected(t *testing.T) {
	t.Parallel()
	i := new(Inserter)
	_, _, err := i.Table("users").
		Columns([]string{"name); DROP TABLE users;--"}).Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sql identifier")
}

// TestInjection_DottedIdentifierAccepted proves that schema-qualified
// names still pass validation.
func TestInjection_DottedIdentifierAccepted(t *testing.T) {
	t.Parallel()
	expr := new(Expression)
	sql, args, err := expr.
		Where(builder.Clause{Left: "users.name", Operator: "=", Right: "alice"}).Express()
	assert.NoError(t, err)
	assert.True(t, strings.Contains(sql, "users.name = ?"))
	assert.Equal(t, []any{"alice"}, args)
}

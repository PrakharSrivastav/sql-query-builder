package postgres

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"

	"github.com/stretchr/testify/assert"
)

func TestPostgresReader_RewritesPlaceholders(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "name", Operator: "=", Right: "alice"}).
		And(builder.Clause{Left: "age", Operator: ">", Right: 30}).
		In("status", "active", "pending")

	sql, args, err := sqlb.Reader.
		Select("id", "name").
		From("users").
		Condition(expr).
		Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT id, name FROM users WHERE (name = $1) AND (age > $2) AND (status IN ($3, $4)) ;", sql)
	assert.Equal(t, []any{"alice", 30, "active", "pending"}, args)
}

func TestPostgresUpdater_RewritesPlaceholdersAcrossSetAndCondition(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "id", Operator: "=", Right: 42})

	sql, args, err := sqlb.Updater.
		Update("users").
		Set(map[string]interface{}{"name": "alice", "age": 31}).
		Condition(expr).
		Build()
	assert.NoError(t, err)
	// Set keys are sorted: age first, then name; condition arg follows.
	assert.Equal(t, "UPDATE users SET age=$1, name=$2 WHERE (id = $3) ;", sql)
	assert.Equal(t, []any{31, "alice", 42}, args)
}

func TestPostgresInserter_RewritesPlaceholdersAcrossRows(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	stmt := sqlb.Inserter.Table("users").Columns([]string{"name", "age"})
	stmt.Values(builder.Value{"name": "a", "age": 1})
	stmt.Values(builder.Value{"name": "b", "age": 2})

	sql, args, err := stmt.Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO users ( age, name ) values ($1, $2),($3, $4);", sql)
	assert.Equal(t, []any{1, "a", 2, "b"}, args)
}

func TestPostgresCreater_NoPlaceholders(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	sql, args, err := sqlb.Creater.Table("users").SetColumns([]builder.Columns{
		{Name: "id", Datatype: "uuid", Constraint: "PRIMARY KEY"},
	}).Build()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "CREATE TABLE users ( id uuid PRIMARY KEY );", sql)
}

func TestPostgresInjection_ValueIsParameterized(t *testing.T) {
	t.Parallel()
	payload := "x' OR 1=1--"
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "name", Operator: "=", Right: payload})
	sql, args, err := sqlb.Reader.Select("id").From("users").Condition(expr).Build()
	assert.NoError(t, err)
	assert.NotContains(t, sql, "OR")
	assert.NotContains(t, sql, "'")
	assert.Equal(t, "SELECT id FROM users WHERE (name = $1) ;", sql)
	assert.Equal(t, []any{payload}, args)
}

// TestPostgresReader_AllFluentMethods exercises every Reader delegate
// (FromAlias, OrderBy, Limit, Offset, RawCondition, joins, On, GroupBy,
// Having) so the wrappers don't decay silently.
func TestPostgresReader_AllFluentMethods(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	sql, _, err := sqlb.Reader.
		Select("a.f1", "b.f2").
		FromAlias(builder.Alias{Name: "t1", Alias: "a"}).
		InnerJoin("t2 as b").On("a.id = b.id").
		LeftJoin("t3 as c").On("a.id = c.id").
		RightJoin("t4 as d").On("a.id = d.id").
		RawCondition(" WHERE a.f1 IS NOT NULL").
		GroupBy([]string{"a.f1"}).
		Having("count(*) > 1").
		OrderBy("a.f1").
		Limit(10).
		Offset(5).
		Build()
	assert.NoError(t, err)
	assert.Contains(t, sql, "INNER JOIN t2 as b ON a.id = b.id")
	assert.Contains(t, sql, "LEFT JOIN t3 as c")
	assert.Contains(t, sql, "RIGHT JOIN t4 as d")
	assert.Contains(t, sql, "GROUP BY a.f1")
	assert.Contains(t, sql, "HAVING count(*) > 1")
	assert.Contains(t, sql, "ORDER BY a.f1")
	assert.Contains(t, sql, "LIMIT 10")
	assert.Contains(t, sql, "OFFSET 5")
}

// TestPostgresUpdater_RawCondition exercises the postgres Updater's
// RawCondition delegate.
func TestPostgresInserter_Returning(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	sql, args, err := sqlb.Inserter.
		Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		Returning("id", "created_at").Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO users ( name ) values ($1) RETURNING id, created_at;", sql)
	assert.Equal(t, []any{"alice"}, args)
}

func TestPostgresInserter_OnConflictDoUpdate_Excluded(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	sql, args, err := sqlb.Inserter.
		Table("users").
		Columns([]string{"name"}).
		Values(builder.Value{"name": "alice"}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{
			"name":       builder.Excluded{Col: "name"},
			"updated_at": 12345,
		}).
		Returning("id").Build()
	assert.NoError(t, err)
	assert.Equal(t, "INSERT INTO users ( name ) values ($1) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, updated_at = $2 RETURNING id;", sql)
	assert.Equal(t, []any{"alice", 12345}, args)
}

func TestPostgresUpdater_Returning(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "id", Operator: "=", Right: 42})
	sql, args, err := sqlb.Updater.
		Update("users").
		Set(map[string]interface{}{"name": "alice"}).
		Condition(expr).
		Returning("id").Build()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE users SET name=$1 WHERE (id = $2) RETURNING id ;", sql)
	assert.Equal(t, []any{"alice", 42}, args)
}

func TestPostgresUpdater_RawCondition(t *testing.T) {
	t.Parallel()
	sqlb, err := NewPostgresBuilder()
	assert.NoError(t, err)

	sql, args, err := sqlb.Updater.
		Update("users").
		Set(map[string]interface{}{"name": "alice"}).
		RawCondition("WHERE id = 1").
		Build()
	assert.NoError(t, err)
	assert.Equal(t, "UPDATE users SET name=$1 WHERE id = 1  ;", sql)
	assert.Equal(t, []any{"alice"}, args)
}

func TestRewritePlaceholders_Standalone(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "no markers", rewritePlaceholders("no markers"))
	assert.Equal(t, "$1 $2 $3", rewritePlaceholders("? ? ?"))
	assert.Equal(t, "a=$1 AND b=$2", rewritePlaceholders("a=? AND b=?"))
	assert.Equal(t, "", rewritePlaceholders(""))
	assert.Equal(t, "$1", rewritePlaceholders("?"))
	// Adjacent and trailing placeholders.
	assert.Equal(t, "$1$2", rewritePlaceholders("??"))
	assert.Equal(t, "x=$1", rewritePlaceholders("x=?"))
}

// TestRewritePlaceholders_RawConditionCollision documents (and pins)
// the known limitation: any `?` inside a RawCondition gets rewritten,
// including inside string literals and Postgres JSON `?` operators.
// Use the parameterized Condition path to avoid this.
func TestRewritePlaceholders_RawConditionCollision(t *testing.T) {
	t.Parallel()
	// `?` inside a string literal is rewritten — broken SQL.
	assert.Equal(t, "WHERE name = 'who$1'", rewritePlaceholders("WHERE name = 'who?'"))
	// JSON existence operator `?` is rewritten — broken SQL.
	assert.Equal(t, "WHERE data $1 'key'", rewritePlaceholders("WHERE data ? 'key'"))
}

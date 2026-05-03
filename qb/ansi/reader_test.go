package ansi

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	t.Parallel()
	c := new(Reader)
	assert.NotNil(t, c)
}

func TestSimpleSelect(t *testing.T) {
	t.Parallel()
	c := new(Reader)

	sql, args, err := c.Select("field1", "field2").Build()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "SELECT field1, field2 ;", sql)

	sql, args, err = c.Select("field1", "field2").From("xyz").Build()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "SELECT field1, field2 FROM xyz ;", sql)

	sql, _, err = c.Select("field1", "field2").From("xyz", "abc").Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT field1, field2 FROM xyz, abc ;", sql)

	sql, _, err = c.Select("a.field1", "b.field2").FromAlias(
		builder.Alias{Name: "table1", Alias: "a"},
		builder.Alias{Name: "table2", Alias: "b"},
	).Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a, table2 as b ;", sql)
}

func TestExpressionBuilder(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	r.Select("a.field1", "b.field2").FromAlias(
		builder.Alias{Name: "table1", Alias: "a"},
		builder.Alias{Name: "table2", Alias: "b"},
	)

	whr := new(Expression)
	whr.Where(builder.Clause{Left: "field1", Operator: "=", Right: "abc"}).
		And(builder.Clause{Left: "field2", Operator: "=", Right: 456})
	r.Condition(whr)

	sql, args, err := r.Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a, table2 as b WHERE (field1 = ?) AND (field2 = ?) ;", sql)
	assert.Equal(t, []any{"abc", 456}, args)
}

func TestRawWhereClause(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	r.Select("a.field1", "b.field2").FromAlias(
		builder.Alias{Name: "table1", Alias: "a"},
		builder.Alias{Name: "table2", Alias: "b"},
	).RawCondition(" WHERE field1 = 'abc'")

	sql, args, err := r.Build()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a, table2 as b WHERE field1 = 'abc' ;", sql)
}

func TestLimitAndOffset(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql, _, err := r.Select("a.field1", "b.field2").From("xyz").Limit(10).Offset(20).Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM xyz LIMIT 10 OFFSET 20 ;", sql)
}

func TestOrderBy(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql, _, err := r.Select("field1", "field2").From("xyz").OrderBy("field1").Limit(10).Offset(20).Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT field1, field2 FROM xyz ORDER BY field1 LIMIT 10 OFFSET 20 ;", sql)
}

func TestLeftJoin(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql, _, err := r.
		Select("a.field1", "b.field2").
		FromAlias(builder.Alias{Name: "table1", Alias: "a"}).
		LeftJoin("table2 as b").
		On("a.field3 = b.field2").Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a LEFT JOIN table2 as b ON a.field3 = b.field2 ;", sql)
}

func TestRightJoin(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql, _, err := r.
		Select("a.field1", "b.field2").
		FromAlias(builder.Alias{Name: "table1", Alias: "a"}).
		RightJoin("table2 as b").
		On("a.field3 = b.field2").Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a RIGHT JOIN table2 as b ON a.field3 = b.field2 ;", sql)
}

func TestInnerJoin(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql, _, err := r.
		Select("a.field1", "b.field2").
		FromAlias(builder.Alias{Name: "table1", Alias: "a"}).
		InnerJoin("table2 as b").
		On("a.field3 = b.field2").Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a INNER JOIN table2 as b ON a.field3 = b.field2 ;", sql)
}

func TestInClause(t *testing.T) {
	t.Parallel()
	whr := new(Expression)
	sql, args, err := whr.Where(builder.Clause{Left: "field1", Operator: "=", Right: "abc"}).
		In("field2", "a", "b", "c").Express()
	assert.NoError(t, err)
	assert.Equal(t, " WHERE (field1 = ?) AND (field2 IN (?, ?, ?))", sql)
	assert.Equal(t, []any{"abc", "a", "b", "c"}, args)
}

func TestOrClause(t *testing.T) {
	t.Parallel()
	whr := new(Expression)
	sql, args, err := whr.
		Where(builder.Clause{Left: "field1", Operator: "=", Right: "abc"}).
		Or(builder.Clause{Left: "field2", Operator: "=", Right: 42}).Express()
	assert.NoError(t, err)
	assert.Equal(t, " WHERE (field1 = ?) OR (field2 = ?)", sql)
	assert.Equal(t, []any{"abc", 42}, args)
}

func TestEmptyInClauseErrors(t *testing.T) {
	t.Parallel()
	whr := new(Expression)
	_, _, err := whr.Where(builder.Clause{Left: "f", Operator: "=", Right: 1}).
		In("field2").Express()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires at least one item")
}

func TestMultipleErrorsJoined(t *testing.T) {
	t.Parallel()
	whr := new(Expression)
	_, _, err := whr.
		Where(builder.Clause{Left: "bad name", Operator: "=", Right: 1}).
		And(builder.Clause{Left: "also bad", Operator: "=", Right: 2}).Express()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ";")
	assert.Contains(t, err.Error(), "bad name")
	assert.Contains(t, err.Error(), "also bad")
}

func TestNotInClause(t *testing.T) {
	t.Parallel()
	whr := new(Expression)
	sql, args, err := whr.Where(builder.Clause{Left: "field1", Operator: "=", Right: "abc"}).
		NotIn("field2", 1, 2, 3).Express()
	assert.NoError(t, err)
	assert.Equal(t, " WHERE (field1 = ?) AND (field2 NOT IN (?, ?, ?))", sql)
	assert.Equal(t, []any{"abc", 1, 2, 3}, args)
}

func TestGroupBy(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	fields := []string{"field1", "field2"}
	sql, _, err := r.Select("field1", "field2", "field3").From("xyz").GroupBy(fields).Having("field1 > 500").Build()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT field1, field2, field3 FROM xyz GROUP BY field1, field2 HAVING field1 > 500 ;", sql)
}

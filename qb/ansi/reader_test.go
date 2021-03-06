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
	sql := c.Select("field1", "field2").Build()
	assert.NotNil(t, sql)
	assert.Equal(t, "SELECT field1, field2 ;", sql)

	sql = c.Select("field1", "field2").From("xyz").Build()
	assert.Equal(t, "SELECT field1, field2 FROM xyz ;", sql)

	sql = c.Select("field1", "field2").From("xyz", "abc").Build()
	assert.Equal(t, "SELECT field1, field2 FROM xyz, abc ;", sql)

	sql = c.Select("a.field1", "b.field2").FromAlias(
		builder.Alias{Name: "table1", Alias: "a"},
		builder.Alias{Name: "table2", Alias: "b"},
	).Build()

	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a, table2 as b ;", sql)
}

func TestExpressionBuilder(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql := r.Select("a.field1", "b.field2").FromAlias(
		builder.Alias{Name: "table1", Alias: "a"},
		builder.Alias{Name: "table2", Alias: "b"},
	)

	whr := new(Expression)
	whr.Where(builder.Clause{Left: "field1", Operator: "=", Right: "abc"}).
		And(builder.Clause{Left: "field2", Operator: "=", Right: 456})
	sql.Condition(whr)
	// ex := new(builder.Expression)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a, table2 as b WHERE ( field1 = 'abc'  ) AND ( field2 = 456 ) ;", sql.Build())
}

func TestRawWhereClause(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql := r.Select("a.field1", "b.field2").FromAlias(
		builder.Alias{Name: "table1", Alias: "a"},
		builder.Alias{Name: "table2", Alias: "b"},
	)

	whr := new(Expression)
	clause := whr.Where(builder.Clause{Left: "field1", Operator: "=", Right: "abc"}).
		And(builder.Clause{Left: "field2", Operator: "=", Right: 456}).Express()
	sql.RawCondition(clause)
	// ex := new(builder.Expression)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a, table2 as b WHERE ( field1 = 'abc'  ) AND ( field2 = 456 ) ;", sql.Build())
}

func TestLimitAndOffset(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql := r.Select("a.field1", "b.field2").From("xyz").Limit(10).Offset(20).Build()

	assert.Equal(t, "SELECT a.field1, b.field2 FROM xyz LIMIT 10 OFFSET 20 ;", sql)
}

func TestOrderBy(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql := r.Select("field1", "field2").From("xyz").OrderBy("field1").Limit(10).Offset(20).Build()

	assert.Equal(t, "SELECT field1, field2 FROM xyz ORDER BY field1 LIMIT 10 OFFSET 20 ;", sql)
}

func TestLeftJoin(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql := r.
		Select("a.field1", "b.field2").
		FromAlias(builder.Alias{Name: "table1", Alias: "a"}).
		LeftJoin("table2 as b").
		On("a.field3 = b.field2").Build()

	assert.NotNil(t, sql)
	assert.NotEmpty(t, sql)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a LEFT JOIN table2 as b ON a.field3 = b.field2 ;", sql)
}

func TestRightJoin(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql := r.
		Select("a.field1", "b.field2").
		FromAlias(builder.Alias{Name: "table1", Alias: "a"}).
		RightJoin("table2 as b").
		On("a.field3 = b.field2").Build()

	assert.NotNil(t, sql)
	assert.NotEmpty(t, sql)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a RIGHT JOIN table2 as b ON a.field3 = b.field2 ;", sql)
}

func TestInnerJoin(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	sql := r.
		Select("a.field1", "b.field2").
		FromAlias(builder.Alias{Name: "table1", Alias: "a"}).
		InnerJoin("table2 as b").
		On("a.field3 = b.field2").Build()

	assert.NotNil(t, sql)
	assert.NotEmpty(t, sql)
	assert.Equal(t, "SELECT a.field1, b.field2 FROM table1 as a INNER JOIN table2 as b ON a.field3 = b.field2 ;", sql)
}

func TestGroupBy(t *testing.T) {
	t.Parallel()
	r := new(Reader)
	fields := []string{"field1", "field2"}
	sql := r.Select("field1", "field2", "field3").From("xyz").GroupBy(fields).Having("field1 > 500").Build()
	assert.NotNil(t, sql)
	assert.NotEmpty(t, sql)
	assert.Equal(t, "SELECT field1, field2, field3 FROM xyz GROUP BY field1, field2 HAVING field1 > 500 ;", sql)
}

package ansi

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"

	"github.com/stretchr/testify/assert"
)

// TestQuotedIdentifier_FlowsThroughToSQL pins that delimited identifiers
// reach the emitted SQL verbatim — Postgres / ANSI-mode users with
// reserved-word column or table names rely on this.
func TestQuotedIdentifier_FlowsThroughToSQL(t *testing.T) {
	t.Parallel()

	expr := new(Expression)
	frag, args, err := expr.Where(builder.Clause{
		Left: `"User"`, Operator: "=", Right: 1,
	}).Express()
	assert.NoError(t, err)
	assert.Equal(t, ` WHERE ("User" = ?)`, frag)
	assert.Equal(t, []any{1}, args)

	i := new(Inserter)
	sql, _, err := i.Table(`"order"`).
		Columns([]string{`"User"`}).
		Values(builder.Value{`"User"`: 1}).Build()
	assert.NoError(t, err)
	assert.Equal(t, `INSERT INTO "order" ( "User" ) values (?);`, sql)
}

func TestValidateIdentifier(t *testing.T) {
	t.Parallel()

	accepted := []string{
		"field",
		"_field",
		"field1",
		"table.field",
		"public.users",
		// Delimited identifiers (SQL:1999 / Postgres style).
		`"User"`,
		`"order"`,           // reserved word becomes a legal name when delimited
		`"with""quote"`,     // embedded `"` doubled as `""`
		`"public"."User"`,
		`public."User"`,
		`"User".id`,
	}
	for _, name := range accepted {
		assert.NoError(t, validateIdentifier(name), "expected accepted: %q", name)
	}

	rejected := []string{
		"",
		"1field",            // can't start with digit
		"field-name",        // hyphen not allowed
		"field name",        // space not allowed
		"field;DROP",        // injection attempt
		"a.b.c",             // more than one qualifier dot
		`""`,                // empty delimited identifier
		`"unterminated`,     // unbalanced quote
		"`backtick`",        // MySQL-style not supported
		`"with"quote"`,      // single embedded quote (must be doubled)
	}
	for _, name := range rejected {
		assert.Error(t, validateIdentifier(name), "expected rejected: %q", name)
	}
}

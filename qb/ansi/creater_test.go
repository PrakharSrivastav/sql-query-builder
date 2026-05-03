package ansi

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"

	"github.com/stretchr/testify/assert"
)

func TestCreaterWithPKey(t *testing.T) {
	t.Parallel()
	c := new(Creater)

	columns := []builder.Columns{
		{Name: "id", Datatype: "uuid", Constraint: "PRIMARY KEY"},
		{Name: "name", Datatype: "varchar(200)", Constraint: ""},
		{Name: "description", Datatype: "text", Constraint: ""},
	}
	sql, args, err := c.Table("xyz").SetColumns(columns).Build()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "CREATE TABLE xyz ( id uuid PRIMARY KEY, name varchar(200) , description text  );", sql)
}

func TestCreaterWithoutPKey(t *testing.T) {
	t.Parallel()
	c := new(Creater)

	columns := []builder.Columns{
		{Name: "PersonID", Datatype: "uuid"},
		{Name: "LastName", Datatype: "varchar(200)"},
		{Name: "FirstName", Datatype: "varchar(200)"},
		{Name: "Address", Datatype: "varchar(200)"},
		{Name: "City", Datatype: "varchar(200)"},
	}
	sql, args, err := c.Table("xyz").SetColumns(columns).Build()
	assert.NoError(t, err)
	assert.Empty(t, args)
	assert.Equal(t, "CREATE TABLE xyz ( PersonID uuid , LastName varchar(200) , FirstName varchar(200) , Address varchar(200) , City varchar(200)  );", sql)
}

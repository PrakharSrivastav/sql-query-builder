package ansi

import (
	"testing"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"

	"github.com/stretchr/testify/assert"
)

func TestCreaterWithPKey(t *testing.T) {
	t.Parallel()

	c := new(Creater)

	assert.NotNil(t, c)

	columns := []builder.Columns{
		builder.Columns{Name: "id", Datatype: "uuid", Constraint: "PRIMARY KEY"},
		builder.Columns{Name: "name", Datatype: "varchar(200)", Constraint: ""},
		builder.Columns{Name: "description", Datatype: "text", Constraint: ""},
	}
	sql := c.Table("xyz").SetColumns(columns).Build()

	assert.NotNil(t, sql)

	assert.Equal(t, "CREATE TABLE xyz ( id uuid PRIMARY KEY, name varchar(200) , description text  );", sql)
}

func TestCreaterWithoutPKey(t *testing.T) {
	t.Parallel()

	c := new(Creater)

	assert.NotNil(t, c)

	columns := []builder.Columns{
		builder.Columns{Name: "PersonID", Datatype: "uuid"},
		builder.Columns{Name: "LastName", Datatype: "varchar(200)"},
		builder.Columns{Name: "FirstName", Datatype: "varchar(200)"},
		builder.Columns{Name: "Address", Datatype: "varchar(200)"},
		builder.Columns{Name: "City", Datatype: "varchar(200)"},
	}
	sql := c.Table("xyz").SetColumns(columns).Build()

	assert.NotNil(t, sql)

	assert.Equal(t, "CREATE TABLE xyz ( PersonID uuid , LastName varchar(200) , FirstName varchar(200) , Address varchar(200) , City varchar(200)  );", sql)
}

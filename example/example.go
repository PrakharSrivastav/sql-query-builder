package main

import (
	"fmt"

	"github.com/PrakharSrivastav/sql-query-builder/qb"
	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
	"github.com/PrakharSrivastav/sql-query-builder/qb/core"
)

func main() {
	qbuilder, err := qb.NewSingletonQueryBuilder(core.ANSI)
	if err != nil {
		panic(err)
	}

	// SELECT with parameterized WHERE.
	expr := qbuilder.NewExpression().
		Where(builder.Clause{Left: "field1", Operator: "=", Right: "123"}).
		And(builder.Clause{Left: "field2", Operator: ">", Right: 37})
	sql, args, err := qbuilder.Reader.
		Select("field1", "field2").
		From("xyz").
		Condition(expr).
		Build()
	if err != nil {
		panic(err)
	}
	// SELECT field1, field2 FROM xyz WHERE (field1 = ?) AND (field2 > ?) ;
	// args: ["123", 37]
	fmt.Println(sql)
	fmt.Println(args)
	// Caller passes both to the driver:
	//   db.Query(sql, args...)

	// UPDATE with parameterized SET and WHERE.
	whereExpr := qbuilder.NewExpression().
		Where(builder.Clause{Left: "field1", Operator: "=", Right: "some value"})
	sql, args, err = qbuilder.Updater.
		Update("xyz").
		Set(map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		}).
		Condition(whereExpr).
		Build()
	if err != nil {
		panic(err)
	}
	// UPDATE xyz SET field1=?, field2=? WHERE (field1 = ?) ;
	// args: ["value1", "value2", "some value"]
	fmt.Println(sql)
	fmt.Println(args)
}

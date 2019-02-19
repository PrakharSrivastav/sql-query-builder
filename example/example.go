package main

import (
	"fmt"

	"github.com/PrakharSrivastav/sql-query-builder/qb"
	"github.com/PrakharSrivastav/sql-query-builder/qb/core"
)

func main() {
	builder, err := qb.NewSingletonQueryBuilder(core.ANSI)
	if err != nil {
		panic(err)
	}

	// Write a select statement

	sql := builder.Reader.
		Select("field1", "field2").
		From("xyz").
		RawCondition(" Where field1 = '123' AND field2>37").
		Build()

	// SELECT field1, field2 FROM xyz Where field1 = '123' AND field2>37 ;
	fmt.Println(sql)

	columns := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}
	sql = builder.Updater.
		Update("xyz").
		Set(columns).
		RawCondition(" WHERE field1 = 'some value'").
		Build()
	// UPDATE table xyz SET field1='value1', field2='value2'  WHERE field1 = 'some value'  ;
	fmt.Println(sql)
}

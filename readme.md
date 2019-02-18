# sql-query-builder

A lot of go packages do a wonderful job in creating sql queries from structs. This approach is applicable when you know the structure of the data beforehand. In one of the project which I am working on, the use case demands to prepare and populate the tables on the fly from the data which is mostly map of interfaces, slices (basically unknown structures).

The reason for not having a concrete structure is mostly project driven. We ingest a lot of streaming data from various sources, where each datasource maintains its own data structure/schema. Using ORM like GORM (or packages from other packages) makes more sense if structure for the incoming dataset is known beforehand. 

The intention for this package is to act as a query builder that can be used to create simple DDL and DML scripts. If you have your data avaialable as slice, map or similar golang type, it should be fairly easy to use this package.

This package provides Fluent Query Builder that can be used like

```
builder, err := qb.NewSingletonQueryBuilder(core.PGSQL)
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
```

Special mention to the [Beego](https://beego.me/docs/mvc/model/querybuilder.md) framework. The Reader interface in this project is highly inspired by Beego's Query Builder interface. However, this project introduces other interfaces (Inserter, Creater, Updater) to compliment more db scripting scenrios.

There are other interesting projects that touch the same idea notably, but they provide a lot more than I need.
- [loukoum](https://github.com/ulule/loukoum)
- [squirrel](https://github.com/Masterminds/squirrel)
- [goqu](https://github.com/Masterminds/squirrel) 

## Supported databases
- ANSI SQL : Right now my focus is to complete the maximum coverage for standard ANSI SQL clauses. Support for other dialects will be added eventually.

**Note** : Contributions and PRs are welcome.

## Supported operations
Checkout the test cases and examples to under usage better.
# sql-query-builder

A lot of go packages do a wonderful job in creating sql queries from structs. This approach is applicable when you know the structure of the data beforehand. In one of the project which I am working on, the use case demands to prepare and populate the tables on the fly.

The packages like sql, sqlx gorm provide a lot of functionality out of the box if you are aware of the underlying structure of the domain model. You can create structs for the known model type and then scan the values easily. However, in my case, I am unaware of the domain model beforehand. 

This project tries to simplify this problem. It assumes that the table schema is provided as the map[string]interface which represents (column name, datatype) combination and prepares a query for you which you can execute using your db library of choice.

There are other interesting projects that touch the same idea notably [loukoum](https://github.com/ulule/loukoum), [squirrel](https://github.com/Masterminds/squirrel), [goqu](https://github.com/Masterminds/squirrel) but they provide a lot more than I needed.

## Supported databases
- Postgres

## Supported operations
- CreateTable : Send a map[string]interface{} to generate create table sql.
- Get : Create a select clause. Since `select` is a keyword in Go, I use Get instead.
- Update : create an update sql clause.
- Insert: create an insert statement. Same api for single vs multi-insert

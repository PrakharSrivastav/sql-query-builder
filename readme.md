# sql-query-builder

A lot of go packages do a wonderful job in creating sql queries from structs. This approach is applicable when you know the structure of the data beforehand. In one of the project which I am working on, the use case demands to prepare and populate the tables on the fly from the data which is mostly map of interfaces, slices (basically unknown structures).

The reason for not having a concrete structure is mostly project driven. We ingest a lot of streaming data from various sources, where each datasource maintains its own data structure/schema. Using ORM like GORM (or packages from other packages) makes more sense if structure for the incoming dataset is known beforehand. 

The intention for this package is to act as a query builder that can be used to create simple DDL and DML scripts. If you have your data avaialable as slice, map or similar golang type, it should be fairly easy to use this package.

This package provides a Fluent Query Builder that emits parameterized SQL. `Build()` returns `(sql string, args []any, err error)` — pass them straight to `db.Query` / `db.Exec`. Values are bound through `?` placeholders; identifiers (table and column names) are validated.

```go
qbuilder, err := qb.NewSingletonQueryBuilder(core.ANSI)
if err != nil {
    panic(err)
}

// SELECT with parameterized WHERE
expr := qbuilder.NewExpression().
    Where(builder.Clause{Left: "field1", Operator: "=", Right: "123"}).
    And(builder.Clause{Left: "field2", Operator: ">", Right: 37})

sql, args, err := qbuilder.Reader.
    Select("field1", "field2").
    From("xyz").
    Condition(expr).
    Build()
// sql:  SELECT field1, field2 FROM xyz WHERE (field1 = ?) AND (field2 > ?) ;
// args: ["123", 37]
rows, err := db.Query(sql, args...)

// UPDATE with parameterized SET and WHERE
whereExpr := qbuilder.NewExpression().
    Where(builder.Clause{Left: "field1", Operator: "=", Right: "some value"})

sql, args, err = qbuilder.Updater.
    Update("xyz").
    Set(map[string]interface{}{"field1": "value1", "field2": "value2"}).
    Condition(whereExpr).
    Build()
// sql:  UPDATE xyz SET field1=?, field2=? WHERE (field1 = ?) ;
// args: ["value1", "value2", "some value"]
_, err = db.Exec(sql, args...)
```

### Safety boundaries

- **Values** — `Clause.Right`, `Set` map values, `Insert.Values`, `In`/`NotIn` items: bound via `?`. Safe for untrusted input.
- **Validated identifiers** — `Update` table, `Insert.Table`, `Insert.Columns`, `Set` keys, `Where` / `And` / `Or` / `In` / `NotIn` left-hand side, `Creater.Table`, `Creater.SetColumns[i].Name`: must match `[A-Za-z_][A-Za-z0-9_]*` (optionally one `.` qualifier). Violations surface from `Build()` as an error.
- **Caller-trusted free-form** — `Select` args, `From`, `OrderBy`, `GroupBy`, `Having`, `On`, joins, `FromAlias` names, `RawCondition`, `Creater.Columns.Datatype` and `.Constraint`: written verbatim. Do not pass untrusted input here.

Special mention to the [Beego](https://beego.me/docs/mvc/model/querybuilder.md) framework. The Reader interface in this project is highly inspired by Beego's Query Builder interface. However, this project introduces other interfaces (Inserter, Creater, Updater) to compliment more db scripting scenrios.

There are other interesting projects that touch the same idea notably, but they provide a lot more than I need.
- [loukoum](https://github.com/ulule/loukoum)
- [squirrel](https://github.com/Masterminds/squirrel)
- [goqu](https://github.com/Masterminds/squirrel) 

## Supported dialects

| Dialect | Constant | Placeholders | Notes |
|---|---|---|---|
| ANSI | `core.ANSI` | `?` | Default. Portable to most drivers. |
| PostgreSQL | `core.PGSQL` | `$1, $2, ...` | Rewrites placeholders at `Build()`. **Otherwise identical to ANSI** — see scope below. |
| MySQL | `core.MYSQL` | `?` | Falls through to ANSI. |
| SQLite | `core.SQLITE` | `?` | Falls through to ANSI. |

### What's actually in scope

This package generates **simple, portable DDL/DML**: `SELECT` (with WHERE/JOIN/GROUP BY/HAVING/ORDER BY/LIMIT/OFFSET), `INSERT`, `UPDATE`, `CREATE TABLE`. Values are parameterized; identifiers are validated. That's it.

The Postgres dialect today differs from ANSI **only in placeholder syntax**. The following Postgres-specific features are **not implemented** — use `RawCondition` as an escape hatch, or another library if you need them as first-class:

- `RETURNING` clause
- `ON CONFLICT (...) DO UPDATE / DO NOTHING` (upsert)
- `ILIKE`, array operators (`@>`, `&&`, `ANY`, `ALL`)
- Identifier quoting for reserved words (the validator rejects them outright)
- Multi-level schema qualification (only `schema.table` — not `db.schema.table`)
- `WITH` (CTE), window functions, `EXCLUDED.*`
- Cast syntax `value::type`
- `CREATE TABLE IF NOT EXISTS`, partitioning, generated columns

Postgres-specific **column types** (`uuid`, `jsonb`, `timestamptz`, ...) work because `Creater.Columns.Datatype` is a free-form string the caller supplies verbatim. Postgres-specific **operators or values** in WHERE clauses must be passed through `RawCondition`, which is caller-trusted.

**Note** : Contributions and PRs are welcome.
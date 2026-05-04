//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/PrakharSrivastav/sql-query-builder/qb/builder"
	"github.com/PrakharSrivastav/sql-query-builder/qb/core"
)

// startDB returns a connected *sql.DB and a postgres query builder.
//
// If TEST_POSTGRES_DSN is set it connects directly. Otherwise it spins up
// a postgres:16-alpine container. wait.ForLog is used instead of wait.ForSQL
// because ForLog reads the container log stream without calling the Docker
// inspect API — which deadlocks on Docker Desktop ≥4.67 during container
// startup. MappedPort (called by ConnectionString) is safe to call after
// Run() returns because Docker Desktop releases its internal lock by then.
func startDB(t *testing.T) (*sql.DB, *core.SQL) {
	t.Helper()

	var dsn string
	if v := os.Getenv("TEST_POSTGRES_DSN"); v != "" {
		dsn = v
	} else {
		testcontainers.SkipIfProviderIsNotHealthy(t)

		ctx := context.Background()
		ctr, err := tcpostgres.Run(ctx,
			"postgres:16-alpine",
			tcpostgres.WithDatabase("testdb"),
			tcpostgres.WithUsername("user"),
			tcpostgres.WithPassword("password"),
			// ForLog reads the log stream (no inspect API call). Two
			// occurrences: postgres logs this once on init, then again on
			// the real startup after the init cycle shuts it down.
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(60*time.Second),
			),
		)
		// Register cleanup BEFORE error check (skill best practice).
		testcontainers.CleanupContainer(t, ctr)
		require.NoError(t, err)

		var cerr error
		dsn, cerr = ctr.ConnectionString(ctx, "sslmode=disable")
		require.NoError(t, cerr)
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	require.NoError(t, db.PingContext(ctx))

	sqlb, err := NewPostgresBuilder()
	require.NoError(t, err)

	return db, sqlb
}

func createUsersTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			name       TEXT NOT NULL,
			age        INT,
			status     TEXT DEFAULT 'active',
			updated_at BIGINT,
			created_at BIGINT DEFAULT extract(epoch from now())
		)
	`)
	require.NoError(t, err)
}

// TestIntegration_CreaterDDL verifies CREATE TABLE reaches the DB.
func TestIntegration_CreaterDDL(t *testing.T) {
	db, sqlb := startDB(t)

	q, args, err := sqlb.Creater.Table("products").SetColumns([]builder.Columns{
		{Name: "id", Datatype: "SERIAL", Constraint: "PRIMARY KEY"},
		{Name: "name", Datatype: "TEXT", Constraint: "NOT NULL"},
		{Name: "price", Datatype: "NUMERIC(10,2)"},
	}).Build()
	require.NoError(t, err)
	assert.Empty(t, args)

	_, err = db.Exec(q)
	require.NoError(t, err)
}

// TestIntegration_InserterSingle inserts one row and checks row count.
func TestIntegration_InserterSingle(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)

	q, args, err := sqlb.Inserter.
		Table("users").
		Columns([]string{"name", "age"}).
		Values(builder.Value{"name": "alice", "age": 30}).
		Build()
	require.NoError(t, err)

	res, err := db.Exec(q, args...)
	require.NoError(t, err)
	n, _ := res.RowsAffected()
	assert.Equal(t, int64(1), n)
}

// TestIntegration_InserterMultiRow inserts several rows in one statement.
func TestIntegration_InserterMultiRow(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)

	stmt := sqlb.Inserter.Table("users").Columns([]string{"name", "age"})
	stmt.Values(builder.Value{"name": "bob", "age": 25})
	stmt.Values(builder.Value{"name": "carol", "age": 35})
	stmt.Values(builder.Value{"name": "dave", "age": 40})

	q, args, err := stmt.Build()
	require.NoError(t, err)

	res, err := db.Exec(q, args...)
	require.NoError(t, err)
	n, _ := res.RowsAffected()
	assert.Equal(t, int64(3), n)
}

// TestIntegration_InserterReturning checks RETURNING delivers the generated id.
func TestIntegration_InserterReturning(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)

	q, args, err := sqlb.Inserter.
		Table("users").
		Columns([]string{"name", "age"}).
		Values(builder.Value{"name": "eve", "age": 28}).
		Returning("id", "name").
		Build()
	require.NoError(t, err)

	var id int
	var name string
	err = db.QueryRow(q, args...).Scan(&id, &name)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
	assert.Equal(t, "eve", name)
}

// TestIntegration_OnConflictDoNothing proves duplicate inserts don't error.
func TestIntegration_OnConflictDoNothing(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (id, name, age) VALUES (99, 'fixed', 0)") //nolint:errcheck

	q, args, err := sqlb.Inserter.
		Table("users").
		Columns([]string{"id", "name", "age"}).
		Values(builder.Value{"id": 99, "name": "dup", "age": 0}).
		OnConflictDoNothing("id").
		Build()
	require.NoError(t, err)

	_, err = db.Exec(q, args...)
	require.NoError(t, err)

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM users WHERE id = 99").Scan(&count))
	assert.Equal(t, 1, count)
}

// TestIntegration_OnConflictDoUpdate_Excluded verifies EXCLUDED.<col> upsert.
func TestIntegration_OnConflictDoUpdate_Excluded(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (id, name, age) VALUES (1, 'original', 10)") //nolint:errcheck

	q, args, err := sqlb.Inserter.
		Table("users").
		Columns([]string{"id", "name", "age"}).
		Values(builder.Value{"id": 1, "name": "updated", "age": 99}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{
			"name": builder.Excluded{Col: "name"},
			"age":  builder.Excluded{Col: "age"},
		}).
		Returning("name", "age").
		Build()
	require.NoError(t, err)

	var name string
	var age int
	err = db.QueryRow(q, args...).Scan(&name, &age)
	require.NoError(t, err)
	assert.Equal(t, "updated", name)
	assert.Equal(t, 99, age)
}

// TestIntegration_ReaderSelectWhereIn exercises SELECT + WHERE + IN.
func TestIntegration_ReaderSelectWhereIn(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (name, age, status) VALUES ('alice', 30, 'active'), ('bob', 25, 'inactive'), ('carol', 35, 'active')") //nolint:errcheck

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "age", Operator: ">", Right: 24}).
		In("status", "active")

	q, args, err := sqlb.Reader.
		Select("name", "age").
		From("users").
		Condition(expr).
		OrderBy("name").
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		var age int
		require.NoError(t, rows.Scan(&name, &age))
		names = append(names, name)
	}
	require.NoError(t, rows.Err())
	assert.Equal(t, []string{"alice", "carol"}, names)
}

// TestIntegration_ReaderLimitOffset verifies pagination arguments work end-to-end.
func TestIntegration_ReaderLimitOffset(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	for i := range 10 {
		db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", fmt.Sprintf("user%02d", i), i) //nolint:errcheck
	}

	q, args, err := sqlb.Reader.
		Select("name").
		From("users").
		OrderBy("name").
		Limit(3).
		Offset(2).
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	require.NoError(t, rows.Err())
	assert.Equal(t, 3, count)
}

// TestIntegration_UpdaterSetAndReturning verifies UPDATE with RETURNING.
func TestIntegration_UpdaterSetAndReturning(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (name, age) VALUES ('frank', 50)") //nolint:errcheck

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "name", Operator: "=", Right: "frank"})

	q, args, err := sqlb.Updater.
		Update("users").
		Set(map[string]interface{}{"age": 51}).
		Condition(expr).
		Returning("age").
		Build()
	require.NoError(t, err)

	var age int
	err = db.QueryRow(q, args...).Scan(&age)
	require.NoError(t, err)
	assert.Equal(t, 51, age)
}

// TestIntegration_QuotedIdentifier verifies delimited identifiers reach the DB unchanged.
func TestIntegration_QuotedIdentifier(t *testing.T) {
	db, sqlb := startDB(t)

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "Users" ( "Name" TEXT, "Age" INT )`)
	require.NoError(t, err)

	q, args, err := sqlb.Inserter.
		Table(`"Users"`).
		Columns([]string{`"Name"`, `"Age"`}).
		Values(builder.Value{`"Name"`: "grace", `"Age"`: 22}).
		Returning(`"Name"`).
		Build()
	require.NoError(t, err)

	var name string
	err = db.QueryRow(q, args...).Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "grace", name)
}

// TestIntegration_InjectionValueIsParameterized proves a SQL injection payload
// in a WHERE value is never executed — the row simply isn't found.
func TestIntegration_InjectionValueIsParameterized(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (name, age) VALUES ('safe', 1)") //nolint:errcheck

	payload := "safe' OR '1'='1"
	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "name", Operator: "=", Right: payload})

	q, args, err := sqlb.Reader.
		Select("name").
		From("users").
		Condition(expr).
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	assert.False(t, rows.Next(), "injection payload must not match any rows")
	require.NoError(t, rows.Err())
}

func createPostsTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id      SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			title   TEXT NOT NULL
		)
	`)
	require.NoError(t, err)
}

// TestIntegration_ReaderOrClause verifies WHERE ... OR ... returns the right rows.
func TestIntegration_ReaderOrClause(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (name, age) VALUES ('alice', 18), ('bob', 30), ('carol', 50)") //nolint:errcheck

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "age", Operator: "<", Right: 20}).
		Or(builder.Clause{Left: "age", Operator: ">", Right: 45})

	q, args, err := sqlb.Reader.
		Select("name").
		From("users").
		Condition(expr).
		OrderBy("name").
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		names = append(names, name)
	}
	require.NoError(t, rows.Err())
	assert.Equal(t, []string{"alice", "carol"}, names)
}

// TestIntegration_ReaderNotIn verifies WHERE col NOT IN (...) excludes the right rows.
func TestIntegration_ReaderNotIn(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (name, status) VALUES ('alice', 'active'), ('bob', 'inactive'), ('carol', 'banned')") //nolint:errcheck

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "id", Operator: ">", Right: 0}).
		NotIn("status", "inactive", "banned")

	q, args, err := sqlb.Reader.
		Select("name").
		From("users").
		Condition(expr).
		OrderBy("name").
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		names = append(names, name)
	}
	require.NoError(t, rows.Err())
	assert.Equal(t, []string{"alice"}, names)
}

// TestIntegration_ReaderInnerJoin verifies INNER JOIN returns only matched rows.
func TestIntegration_ReaderInnerJoin(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	createPostsTable(t, db)

	var userID int
	require.NoError(t, db.QueryRow("INSERT INTO users (name) VALUES ('alice') RETURNING id").Scan(&userID))
	db.Exec("INSERT INTO users (name) VALUES ('bob')") //nolint:errcheck — bob has no post
	db.Exec(fmt.Sprintf("INSERT INTO posts (user_id, title) VALUES (%d, 'hello')", userID)) //nolint:errcheck

	q, args, err := sqlb.Reader.
		Select("u.name", "p.title").
		FromAlias(builder.Alias{Name: "users", Alias: "u"}).
		InnerJoin("posts p").On("u.id = p.user_id").
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	require.NoError(t, rows.Err())
	assert.Equal(t, 1, count, "only alice has a post; bob must be excluded by INNER JOIN")
}

// TestIntegration_ReaderLeftJoin verifies LEFT JOIN includes unmatched left rows.
func TestIntegration_ReaderLeftJoin(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	createPostsTable(t, db)

	var aliceID int
	require.NoError(t, db.QueryRow("INSERT INTO users (name) VALUES ('alice') RETURNING id").Scan(&aliceID))
	db.Exec("INSERT INTO users (name) VALUES ('bob')") //nolint:errcheck — bob has no post
	db.Exec(fmt.Sprintf("INSERT INTO posts (user_id, title) VALUES (%d, 'hello')", aliceID)) //nolint:errcheck

	q, args, err := sqlb.Reader.
		Select("u.name", "p.title").
		FromAlias(builder.Alias{Name: "users", Alias: "u"}).
		LeftJoin("posts p").On("u.id = p.user_id").
		OrderBy("u.name").
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	require.NoError(t, rows.Err())
	assert.Equal(t, 2, count, "LEFT JOIN must include bob even though he has no post")
}

// TestIntegration_ReaderRightJoin verifies RIGHT JOIN includes unmatched right rows.
func TestIntegration_ReaderRightJoin(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	createPostsTable(t, db)

	var aliceID int
	require.NoError(t, db.QueryRow("INSERT INTO users (name) VALUES ('alice') RETURNING id").Scan(&aliceID))
	db.Exec(fmt.Sprintf("INSERT INTO posts (user_id, title) VALUES (%d, 'matched'), (999, 'orphan')", aliceID)) //nolint:errcheck

	q, args, err := sqlb.Reader.
		Select("u.name", "p.title").
		FromAlias(builder.Alias{Name: "users", Alias: "u"}).
		RightJoin("posts p").On("u.id = p.user_id").
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}
	require.NoError(t, rows.Err())
	assert.Equal(t, 2, count, "RIGHT JOIN must include the orphan post (no matching user)")
}

// TestIntegration_ReaderGroupByHaving verifies GROUP BY + HAVING aggregates correctly.
func TestIntegration_ReaderGroupByHaving(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (name, status) VALUES ('alice', 'active'), ('bob', 'active'), ('carol', 'inactive')") //nolint:errcheck

	q, args, err := sqlb.Reader.
		Select("status").
		From("users").
		GroupBy([]string{"status"}).
		Having("COUNT(*) > 1").
		Build()
	require.NoError(t, err)

	rows, err := db.Query(q, args...)
	require.NoError(t, err)
	defer rows.Close()

	var statuses []string
	for rows.Next() {
		var s string
		require.NoError(t, rows.Scan(&s))
		statuses = append(statuses, s)
	}
	require.NoError(t, rows.Err())
	assert.Equal(t, []string{"active"}, statuses, "only 'active' has COUNT > 1")
}

// TestIntegration_InserterReturningStar verifies RETURNING * delivers all columns.
func TestIntegration_InserterReturningStar(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)

	q, args, err := sqlb.Inserter.
		Table("users").
		Columns([]string{"name", "age"}).
		Values(builder.Value{"name": "zara", "age": 25}).
		Returning("*").
		Build()
	require.NoError(t, err)

	// users: id, name, age, status, updated_at, created_at
	var id, age int
	var name, status string
	var updatedAt *int64
	var createdAt int64
	err = db.QueryRow(q, args...).Scan(&id, &name, &age, &status, &updatedAt, &createdAt)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
	assert.Equal(t, "zara", name)
	assert.Equal(t, 25, age)
}

// TestIntegration_OnConflictDoUpdate_PlainValues verifies DO UPDATE SET with literal values.
func TestIntegration_OnConflictDoUpdate_PlainValues(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (id, name, age) VALUES (5, 'orig', 10)") //nolint:errcheck

	q, args, err := sqlb.Inserter.
		Table("users").
		Columns([]string{"id", "name", "age"}).
		Values(builder.Value{"id": 5, "name": "dup", "age": 0}).
		OnConflictDoUpdate([]string{"id"}, map[string]interface{}{
			"name": "updated",
			"age":  99,
		}).
		Returning("name", "age").
		Build()
	require.NoError(t, err)

	var name string
	var age int
	err = db.QueryRow(q, args...).Scan(&name, &age)
	require.NoError(t, err)
	assert.Equal(t, "updated", name)
	assert.Equal(t, 99, age)
}

// TestIntegration_UpdaterMultiColumnSet verifies UPDATE SET with multiple columns.
func TestIntegration_UpdaterMultiColumnSet(t *testing.T) {
	db, sqlb := startDB(t)
	createUsersTable(t, db)
	db.Exec("INSERT INTO users (name, age, status) VALUES ('frank', 30, 'active')") //nolint:errcheck

	expr := sqlb.NewExpression().
		Where(builder.Clause{Left: "name", Operator: "=", Right: "frank"})

	q, args, err := sqlb.Updater.
		Update("users").
		Set(map[string]interface{}{"name": "franklin", "age": 31, "status": "vip"}).
		Condition(expr).
		Returning("name", "age", "status").
		Build()
	require.NoError(t, err)

	var name, status string
	var age int
	err = db.QueryRow(q, args...).Scan(&name, &age, &status)
	require.NoError(t, err)
	assert.Equal(t, "franklin", name)
	assert.Equal(t, 31, age)
	assert.Equal(t, "vip", status)
}

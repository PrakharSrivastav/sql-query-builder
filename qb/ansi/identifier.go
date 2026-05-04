package ansi

import (
	"fmt"
	"regexp"
)

// identRe matches a SQL identifier in either bare or delimited form,
// optionally qualified by one schema/table prefix.
//
// Bare identifier:      [A-Za-z_][A-Za-z0-9_]*
// Delimited identifier: "..."  (any chars; embedded `"` doubled as `""`)
//
// Examples accepted: `field`, `table.field`, `"User"`, `"public"."User"`,
// `"with""quote"`, `public."User"`. This follows the SQL:1999 delimited
// identifier rule, supported by Postgres, SQLite, and ANSI-mode MySQL.
// MySQL backticks are not accepted.
//
// Inside a delimited identifier the content is treated as the literal
// column/table name by the database engine — characters like `;`, space,
// or reserved words are part of the name, not SQL syntax — so allowing
// them is safe.
var identRe = regexp.MustCompile(`^(?:[A-Za-z_][A-Za-z0-9_]*|"(?:[^"]|"")+")(?:\.(?:[A-Za-z_][A-Za-z0-9_]*|"(?:[^"]|"")+"))?$`)

// validateIdentifier returns an error if name does not match the
// permitted identifier shape. Used on column/table names that get
// inlined into SQL (Where left-hand side, Set keys, Insert columns,
// table names) to block injection through identifier slots.
func validateIdentifier(name string) error {
	if !identRe.MatchString(name) {
		return fmt.Errorf("invalid sql identifier: %q", name)
	}
	return nil
}

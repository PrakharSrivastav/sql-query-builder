package ansi

import (
	"fmt"
	"regexp"
)

// identRe matches a SQL identifier, optionally qualified by one schema
// or table prefix: `field`, `table.field`. Anything that needs spaces,
// quotes, function calls, or aliases must go through the explicit
// raw/free-form fields (Select args, OrderBy, Having, joins, etc.).
var identRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)?$`)

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

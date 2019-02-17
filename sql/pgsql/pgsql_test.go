package pgsql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPgSQLBuilder(t *testing.T) {
	t.Parallel()
	sql, err := NewPgSQLBuilder()
	assert.Nil(t, err)
	assert.NotNil(t, sql)
}

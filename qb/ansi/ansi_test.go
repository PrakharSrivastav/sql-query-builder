package ansi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPgSQLBuilder(t *testing.T) {
	t.Parallel()
	sql, err := NewANSIBuilder()
	assert.Nil(t, err)
	assert.NotNil(t, sql)
}

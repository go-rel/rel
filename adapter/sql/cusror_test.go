package sql

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCursor_NopScanner(t *testing.T) {
	assert.Equal(t, &sql.RawBytes{}, (&Cursor{}).NopScanner())
}

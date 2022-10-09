package rel_test

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestSQLQuery(t *testing.T) {
	assert.Equal(t, rel.SQLQuery{
		Statement: "SELECT 1;",
	}, rel.SQL("SELECT 1;"))

	assert.Equal(t, rel.SQLQuery{
		Statement: "SELECT * FROM `users` WHERE id=?;",
		Values:    []any{1},
	}, rel.SQL("SELECT * FROM `users` WHERE id=?;", 1))
}

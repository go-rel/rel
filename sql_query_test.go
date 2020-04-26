package rel_test

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

func TestSQLQuery(t *testing.T) {
	assert.Equal(t, rel.SQLQuery{
		Statement: "SELECT 1",
	}, rel.SQL("SELECT 1"))

	assert.Equal(t, rel.SQLQuery{
		Statement: "SELECT * FROM `users` WHERE id=?",
		Values:    []interface{}{1},
	}, rel.SQL("SELECT * FROM `users` WHERE id=?", 1))
}

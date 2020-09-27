package rel_test

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	assert.Equal(t, rel.GroupQuery{
		Fields: []string{"status"},
	}, rel.NewGroup("status"))
}

func TestGroup_Having(t *testing.T) {
	q := rel.GroupQuery{
		Fields: []string{"status"},
		Filter: rel.Ne("status", "expired"),
	}

	assert.Equal(t, q, rel.NewGroup("status").Having(rel.Ne("status", "expired")))
	assert.Equal(t, q, rel.NewGroup("status").Where(rel.Ne("status", "expired")))
}

func TestGroup_OrHaving(t *testing.T) {
	q := rel.GroupQuery{
		Fields: []string{"status"},
		Filter: rel.Ne("status", "expired").OrNotNil("deleted_at"),
	}

	assert.Equal(t, q, rel.NewGroup("status").Having(rel.Ne("status", "expired")).OrHaving(rel.NotNil("deleted_at")))
	assert.Equal(t, q, rel.NewGroup("status").Where(rel.Ne("status", "expired")).OrWhere(rel.NotNil("deleted_at")))
}

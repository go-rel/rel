package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	assert.Equal(t, grimoire.GroupQuery{
		Fields: []string{"status"},
	}, grimoire.NewGroup("status"))
}

func TestGroup_Having(t *testing.T) {
	q := grimoire.GroupQuery{
		Fields: []string{"status"},
		Filter: grimoire.Ne("status", "expired"),
	}

	assert.Equal(t, q, grimoire.NewGroup("status").Having(grimoire.Ne("status", "expired")))
	assert.Equal(t, q, grimoire.NewGroup("status").Where(grimoire.Ne("status", "expired")))
}

func TestGroup_OrHaving(t *testing.T) {
	q := grimoire.GroupQuery{
		Fields: []string{"status"},
		Filter: grimoire.Ne("status", "expired").OrNotNil("deleted_at"),
	}

	assert.Equal(t, q, grimoire.NewGroup("status").Having(grimoire.Ne("status", "expired")).OrHaving(grimoire.NotNil("deleted_at")))
	assert.Equal(t, q, grimoire.NewGroup("status").Where(grimoire.Ne("status", "expired")).OrWhere(grimoire.NotNil("deleted_at")))
}

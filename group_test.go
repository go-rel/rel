package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	assert.Equal(t, grimoire.GroupClause{
		Fields: []string{"status"},
	}, grimoire.NewGroup("status"))
}

func TestGroup_Having(t *testing.T) {
	q := grimoire.GroupClause{
		Fields: []string{"status"},
		Filter: grimoire.FilterNe("status", "expired"),
	}

	assert.Equal(t, q, grimoire.NewGroup("status").Having(grimoire.FilterNe("status", "expired")))
	assert.Equal(t, q, grimoire.NewGroup("status").Where(grimoire.FilterNe("status", "expired")))
}

func TestGroup_OrHaving(t *testing.T) {
	q := grimoire.GroupClause{
		Fields: []string{"status"},
		Filter: grimoire.FilterNe("status", "expired").OrNotNil("deleted_at"),
	}

	assert.Equal(t, q, grimoire.NewGroup("status").Having(grimoire.FilterNe("status", "expired")).OrHaving(grimoire.FilterNotNil("deleted_at")))
	assert.Equal(t, q, grimoire.NewGroup("status").Where(grimoire.FilterNe("status", "expired")).OrWhere(grimoire.FilterNotNil("deleted_at")))
}

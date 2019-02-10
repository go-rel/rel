package query_test

import (
	"testing"

	"github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	assert.Equal(t, query.GroupClause{
		Fields: []string{"status"},
	}, query.NewGroup("status"))
}

func TestGroup_Having(t *testing.T) {
	q := query.GroupClause{
		Fields: []string{"status"},
		Filter: query.FilterNe("status", "expired"),
	}

	assert.Equal(t, q, query.NewGroup("status").Having(query.FilterNe("status", "expired")))
	assert.Equal(t, q, query.NewGroup("status").Where(query.FilterNe("status", "expired")))
}

func TestGroup_OrHaving(t *testing.T) {
	q := query.GroupClause{
		Fields: []string{"status"},
		Filter: query.FilterNe("status", "expired").OrNotNil("deleted_at"),
	}

	assert.Equal(t, q, query.NewGroup("status").Having(query.FilterNe("status", "expired")).OrHaving(query.FilterNotNil("deleted_at")))
	assert.Equal(t, q, query.NewGroup("status").Where(query.FilterNe("status", "expired")).OrWhere(query.FilterNotNil("deleted_at")))
}

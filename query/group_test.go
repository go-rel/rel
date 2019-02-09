package query_test

import (
	"testing"

	"github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	assert.Equal(t, query.GroupClause{
		Fields: []string{"status"},
	}, query.Group("status"))
}

func TestGroup_Where(t *testing.T) {
	assert.Equal(t, query.GroupClause{
		Fields: []string{"status"},
		Filter: query.FilterNe("status", "expired"),
	}, query.Group("status").Where(query.FilterNe("status", "expired")))
}

package query_test

import (
	"testing"

	"github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	assert.Equal(t, query.SelectClause{
		OnlyDistinct: false,
		Fields:       []string{"id", "name"},
	}, query.Select("id", "name"))
}

func TestSelectDistinct(t *testing.T) {
	assert.Equal(t, query.SelectClause{
		OnlyDistinct: true,
		Fields:       []string{"id", "name"},
	}, query.Select("id", "name").Distinct())
}

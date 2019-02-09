package query_test

import (
	"testing"

	"github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

func TestSortClause_Asc(t *testing.T) {
	assert.True(t, query.SortAsc("score").Asc())
}

func TestSortClause_Desc(t *testing.T) {
	assert.True(t, query.SortDesc("score").Desc())
}

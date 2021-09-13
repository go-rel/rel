package rel_test

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestSortQuery_Asc(t *testing.T) {
	assert.True(t, rel.SortAsc("score").Asc())
}

func TestSortQuery_Desc(t *testing.T) {
	assert.True(t, rel.SortDesc("score").Desc())
}

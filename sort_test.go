package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

func TestSortClause_Asc(t *testing.T) {
	assert.True(t, grimoire.NewSortAsc("score").Asc())
}

func TestSortClause_Desc(t *testing.T) {
	assert.True(t, grimoire.NewSortDesc("score").Desc())
}

package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	assert.Equal(t, grimoire.SelectClause{
		OnlyDistinct: false,
		Fields:       []string{"id", "name"},
	}, grimoire.NewSelect("id", "name"))
}

func TestSelect_Distinct(t *testing.T) {
	assert.Equal(t, grimoire.SelectClause{
		OnlyDistinct: true,
		Fields:       []string{"id", "name"},
	}, grimoire.NewSelect("id", "name").Distinct())
}

package rel_test

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	assert.Equal(t, rel.SelectQuery{
		OnlyDistinct: false,
		Fields:       []string{"id", "name"},
	}, rel.NewSelect("id", "name"))
}

func TestSelect_Distinct(t *testing.T) {
	assert.Equal(t, rel.SelectQuery{
		OnlyDistinct: true,
		Fields:       []string{"id", "name"},
	}, rel.NewSelect("id", "name").Distinct())
}

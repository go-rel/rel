package c

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsc(t *testing.T) {
	asc := Asc("id")

	assert.Equal(t, asc, Order{
		Field: "id",
		Order: 1,
	})
	assert.True(t, asc.Asc())
}

func TestDesc(t *testing.T) {
	desc := Desc("id")

	assert.Equal(t, desc, Order{
		Field: "id",
		Order: -1,
	})
	assert.True(t, desc.Desc())
}

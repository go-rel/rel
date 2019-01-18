package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeset(t *testing.T) {
	ch := Changeset{}
	assert.Nil(t, ch.Errors())
	assert.Nil(t, ch.Error())
	assert.Nil(t, ch.Changes())
	assert.Nil(t, ch.Values())
	assert.Nil(t, ch.Types())
	assert.Nil(t, ch.Constraints())
}

func TestChangeset_Get(t *testing.T) {
	ch := Changeset{
		changes: map[string]interface{}{
			"a": 2,
		},
	}

	assert.Equal(t, 2, ch.Get("a"))
	assert.Equal(t, nil, ch.Get("b"))
	assert.Equal(t, 1, len(ch.changes))
}

func TestChangeset_Fetch(t *testing.T) {
	ch := Changeset{
		changes: map[string]interface{}{
			"a": 1,
		},
		values: map[string]interface{}{
			"b": 2,
		},
	}

	assert.Equal(t, 1, ch.Fetch("a"))
	assert.Equal(t, 2, ch.Fetch("b"))
	assert.Equal(t, nil, ch.Fetch("c"))
	assert.Equal(t, 1, len(ch.changes))
	assert.Equal(t, 1, len(ch.values))
}

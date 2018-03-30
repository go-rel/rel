package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeset(t *testing.T) {
	ch := Changeset{}
	assert.Nil(t, ch.Changes())
	assert.Nil(t, ch.Values())
}

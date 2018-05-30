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

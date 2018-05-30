package changeset

import (
	"testing"

	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

func TestUniqueConstraint(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.Constraints())

	UniqueConstraint(ch, "field1")
	assert.Equal(t, 1, len(ch.Constraints()))
	assert.Equal(t, errors.UniqueConstraint, ch.Constraints()[0].Kind)
}

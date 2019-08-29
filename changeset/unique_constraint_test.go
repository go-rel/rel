package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueConstraint(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.changers)

	UniqueConstraint(ch, "field1")
	assert.Equal(t, 1, len(ch.changers))
}

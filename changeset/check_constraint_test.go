package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckConstraint(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.changers)

	CheckConstraint(ch, "field1")
	assert.Equal(t, 1, len(ch.changers))
}

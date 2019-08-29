package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForeignKeyConstraint(t *testing.T) {
	ch := &Changeset{}
	assert.Nil(t, ch.changers)

	ForeignKeyConstraint(ch, "field1")
	assert.Equal(t, 1, len(ch.changers))
}

package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangesetChanges(t *testing.T) {
	ch := Changeset{}
	assert.Nil(t, ch.Changes())
}

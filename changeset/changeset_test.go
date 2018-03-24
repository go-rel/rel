package changeset

import (
	"testing"
)

func TestChangesetChanges(t *testing.T) {
	ch := Changeset{}

	if ch.Changes() != nil {
		t.Error("Expected nil, but got", ch.Changes())
	}
}

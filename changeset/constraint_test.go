package changeset

import (
	"testing"

	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

func TestConstraint_GetError(t *testing.T) {
	ch := &Changeset{}
	UniqueConstraint(ch, "slug")
	ForeignKeyConstraint(ch, "user_id", Name("user_id_ibfk1"), Exact(true))
	CheckConstraint(ch, "state")

	tests := []struct {
		Error    errors.Error
		Expected errors.Error
	}{
		{
			Error:    errors.New("unique constraint", "slug_unique_index", errors.UniqueConstraint),
			Expected: errors.New("slug has already been taken", "slug", errors.UniqueConstraint),
		},
		{
			Error:    errors.New("foreign key costraint", "user_id_ibfk1", errors.ForeignKeyConstraint),
			Expected: errors.New("does not exist", "user_id", errors.ForeignKeyConstraint),
		},
		{
			Error:    errors.New("check costraint", "state_check", errors.CheckConstraint),
			Expected: errors.New("state is invalid", "state", errors.CheckConstraint),
		},
		{
			Error:    errors.New("unexpected unique constraint", "other_unique_index", errors.UniqueConstraint),
			Expected: errors.NewUnexpected("unexpected unique constraint"),
		},
		{
			Error:    errors.New("unexpected foreign key costraint", "other_id_ibfk1", errors.ForeignKeyConstraint),
			Expected: errors.NewUnexpected("unexpected foreign key costraint"),
		},
		{
			Error:    errors.New("unexpected", "", errors.Unexpected),
			Expected: errors.New("unexpected", "", errors.Unexpected),
		},
		{
			Error:    errors.New("changeset", "", errors.Changeset),
			Expected: errors.New("changeset", "", errors.Changeset),
		},
		{
			Error:    errors.New("not found", "", errors.NotFound),
			Expected: errors.New("not found", "", errors.NotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Error.Error(), func(t *testing.T) {
			assert.Equal(t, tt.Expected, ch.Constraints().GetError(tt.Error))
		})
	}
}

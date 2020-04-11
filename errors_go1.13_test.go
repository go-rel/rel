// +build go1.13

package rel

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstraintError_ErrorsIs(t *testing.T) {
	tests := []struct {
		err    error
		target error
		equal  bool
	}{
		{
			err:    ConstraintError{Type: CheckConstraint},
			target: ErrCheckConstraint,
			equal:  true,
		},
		{
			err:    ConstraintError{Type: UniqueConstraint, Key: "username"},
			target: ErrUniqueConstraint,
			equal:  true,
		},
		{
			err:    ErrUniqueConstraint,
			target: ConstraintError{Type: UniqueConstraint, Key: "username"},
			equal:  true,
		},
		{
			err:    ConstraintError{Type: NotNullConstraint, Key: "username"},
			target: ConstraintError{Type: NotNullConstraint, Key: "email"},
			equal:  false,
		},
		{
			err:    ConstraintError{Type: ForeignKeyConstraint, Key: "book_id"},
			target: ErrNotFound,
			equal:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.err.Error(), func(t *testing.T) {
			assert.Equal(t, test.equal, errors.Is(test.err, test.target))
		})
	}
}

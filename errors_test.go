package rel

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoResultError(t *testing.T) {
	assert.Equal(t, "Record not found", NotFoundError{}.Error())
}

func TestNotFoundError_Is(t *testing.T) {
	tests := []struct {
		err    NotFoundError
		target error
		equal  bool
	}{
		{
			err:    NotFoundError{},
			target: NotFoundError{},
			equal:  true,
		},
		{
			err:    NotFoundError{},
			target: sql.ErrNoRows,
			equal:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.err.Error(), func(t *testing.T) {
			assert.Equal(t, test.equal, test.err.Is(test.target))
		})
	}
}

func TestConstraintType(t *testing.T) {
	assert.Equal(t, "CheckConstraint", CheckConstraint.String())
	assert.Equal(t, "NotNullConstraint", NotNullConstraint.String())
	assert.Equal(t, "UniqueConstraint", UniqueConstraint.String())
	assert.Equal(t, "PrimaryKeyConstraint", PrimaryKeyConstraint.String())
	assert.Equal(t, "ForeignKeyConstraint", ForeignKeyConstraint.String())
	assert.Equal(t, "", ConstraintType(100).String())
}

func TestConstraintError(t *testing.T) {
	err := ConstraintError{Key: "field", Type: UniqueConstraint, Err: errors.New("not unique")}
	assert.NotNil(t, err.Unwrap())
	assert.Equal(t, "UniqueConstraintError: not unique", err.Error())

	err = ConstraintError{Key: "field", Type: UniqueConstraint}
	assert.Nil(t, err.Unwrap())
	assert.Equal(t, "UniqueConstraintError", err.Error())
}

func TestConstraintError_Is(t *testing.T) {
	tests := []struct {
		err    ConstraintError
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
			assert.Equal(t, test.equal, test.err.Is(test.target))
		})
	}
}

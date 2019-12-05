package rel

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoResultError(t *testing.T) {
	assert.Equal(t, "No result found", NoResultError{}.Error())
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

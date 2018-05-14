package errors

import (
	e "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorInterface(t *testing.T) {
	assert.Equal(t, "error", error(Error{Message: "error"}).Error())
}

func TestNew(t *testing.T) {
	assert.Equal(t, Error{
		Message: "error",
		Field:   "error",
		Code:    1000,
	}, New("error", "error", 1000))
}

func TestUnexpectedError(t *testing.T) {
	err := UnexpectedError("error")

	assert.Equal(t, "error", err.Error())
	assert.Equal(t, "", err.Field)
	assert.True(t, err.UnexpectedError())
}

func TestNotFoundError(t *testing.T) {
	err := NotFoundError("error")

	assert.Equal(t, "error", err.Error())
	assert.Equal(t, "", err.Field)
	assert.True(t, err.NotFoundError())
}

func TestChangesetError(t *testing.T) {
	err := ChangesetError("error", "field")

	assert.Equal(t, "error", err.Error())
	assert.Equal(t, "field", err.Field)
	assert.True(t, err.ChangesetError())
}

func TestUniqueConstraintError(t *testing.T) {
	err := UniqueConstraintError("error", "field")

	assert.Equal(t, "error", err.Error())
	assert.Equal(t, "field", err.Field)
	assert.True(t, err.UniqueConstraintError())
}

func TestWrap(t *testing.T) {
	assert.Equal(t, nil, Wrap(nil))
	assert.Equal(t, Error{}, Wrap(Error{}))
	assert.Equal(t, UnexpectedError("error"), Wrap(e.New("error")))
}

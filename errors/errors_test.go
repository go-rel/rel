package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	assert.Equal(t, "error", error(Error{Message: "error"}).Error())
}

func TestNew(t *testing.T) {
	err := New("error", "error", Unexpected)
	assert.Equal(t, Error{
		Message: "error",
		Field:   "error",
		kind:    Unexpected,
	}, err)

	assert.Equal(t, Unexpected, err.Kind())
}

func TestNewWithCode(t *testing.T) {
	err := NewWithCode("error", "error", 1000, Unexpected)
	assert.Equal(t, Error{
		Message: "error",
		Field:   "error",
		Code:    1000,
		kind:    Unexpected,
	}, err)

	assert.Equal(t, Unexpected, err.Kind())
}

func TestNewUnexpected(t *testing.T) {
	err := NewUnexpected("error")
	assert.Equal(t, Error{
		Message: "error",
		kind:    Unexpected,
	}, err)

	assert.Equal(t, Unexpected, err.Kind())
}

package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflectInternalType(t *testing.T) {
	assert.Equal(t,
		reflect.TypeOf(User{}),
		reflectInternalType(User{}),
	)

	assert.Equal(t,
		reflect.TypeOf(User{}),
		reflectInternalType(&User{}),
	)

	assert.Equal(t,
		reflect.TypeOf(User{}),
		reflectInternalType([]User{}),
	)

	assert.Equal(t,
		reflect.TypeOf(User{}),
		reflectInternalType(&[]User{}),
	)

	assert.Equal(t,
		reflect.TypeOf(User{}),
		reflectInternalType([0]User{}),
	)

	assert.Equal(t,
		reflect.TypeOf(User{}),
		reflectInternalType(&[0]User{}),
	)
}

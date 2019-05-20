package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflectInternalType(t *testing.T) {
	type User struct{}

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

func TestReflectTypePtr(t *testing.T) {
	type User struct{}

	assert.Equal(t, reflect.TypeOf(User{}), reflectTypePtr(User{}))
	assert.Equal(t, reflect.TypeOf(User{}), reflectTypePtr(&User{}))

	assert.Panics(t, func() {
		reflectTypePtr("not struct")
	})
}

func TestReflectValuePtr(t *testing.T) {
	type User struct{}

	assert.Equal(t, User{}, reflectValuePtr(User{}).Interface())
	assert.Equal(t, User{}, reflectValuePtr(&User{}).Interface())

	assert.Panics(t, func() {
		reflectValuePtr("not struct")
	})
}

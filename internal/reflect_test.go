package internal

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

func TestReflectInternalStruct(t *testing.T) {
	type User struct{}

	assert.Equal(t, User{}, reflectInternalStruct(User{}).Interface())
	assert.Equal(t, User{}, reflectInternalStruct(&User{}).Interface())

	assert.Panics(t, func() {
		reflectInternalStruct("not struct")
	})
}

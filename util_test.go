package rel

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	assert.Panics(t, func() {
		must(errors.New("error"))
	})
}

func TestMustTrue(t *testing.T) {
	assert.Panics(t, func() {
		mustTrue(true, "error")
	})
}

func TestIsZero(t *testing.T) {
	tests := []any{
		nil,
		false,
		"",
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		uintptr(0),
		float32(0),
		float64(0),
		time.Time{},
		struct{}{},
	}

	for i := range tests {
		t.Run("IsZero", func(t *testing.T) {
			assert.True(t, isZero(tests[i]))
		})
	}
}

func TestIsDeepZero(t *testing.T) {
	v := struct {
		A bool
		B int
		C uint
		D float32
		E complex64
		F [1]int
		G *string
		H string
		J struct {
			JA int
		}
	}{}

	assert.True(t, isDeepZero(reflect.ValueOf(v), 1))
}

func TestIsDeepZero_notZeroArray(t *testing.T) {
	v := [1]bool{true}

	assert.False(t, isDeepZero(reflect.ValueOf(v), 1))
}

func TestIsDeepZero_notZeroStruct(t *testing.T) {
	v := struct {
		A bool
	}{true}

	assert.False(t, isDeepZero(reflect.ValueOf(v), 1))
}

func TestIsDeepZero_reflectInvalid(t *testing.T) {
	assert.True(t, isDeepZero(reflect.Value{}, 1))
}

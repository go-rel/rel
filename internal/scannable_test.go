package internal

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScannable(t *testing.T) {
	vint := 0
	vstring := ""
	vbool := false
	vtime := time.Now()
	vstruct := struct{}{}

	assert.True(t, Scannable(reflect.TypeOf(vint)))
	assert.True(t, Scannable(reflect.TypeOf(vstring)))
	assert.True(t, Scannable(reflect.TypeOf(vbool)))
	assert.True(t, Scannable(reflect.TypeOf(vtime)))
	assert.False(t, Scannable(reflect.TypeOf(vstruct)))

	assert.True(t, Scannable(reflect.TypeOf(&vint)))
	assert.True(t, Scannable(reflect.TypeOf(&vstring)))
	assert.True(t, Scannable(reflect.TypeOf(&vbool)))
	assert.True(t, Scannable(reflect.TypeOf(&vtime)))
	assert.False(t, Scannable(reflect.TypeOf(&vstruct)))
}

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue_Scan(t *testing.T) {
	a := 10
	v := Value(&a).(value)

	v.Scan(5)
	assert.Equal(t, *v.dest.(*int), 5)
}

func TestValue(t *testing.T) {
	a := 10
	v := Value(&a)
	assert.Equal(t, value{dest: &a}, v)
}

type customScanner int

func (*customScanner) Scan(interface{}) error {
	return nil
}

func TestValue_scanner(t *testing.T) {
	a := customScanner(10)
	v := Value(&a)
	assert.Equal(t, &a, v)
}

func TestValue_ptr(t *testing.T) {
	var a *int
	v := Value(&a)
	assert.Equal(t, &a, v)
}

func TestValue_notPtr(t *testing.T) {
	assert.Panics(t, func() {
		Value(0)
	})
}

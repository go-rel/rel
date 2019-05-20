package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullable_Scan(t *testing.T) {
	a := 10
	v := Nullable(&a).(nullable)

	v.Scan(5)
	assert.Equal(t, *v.dest.(*int), 5)
}

func TestNullable(t *testing.T) {
	a := 10
	v := Nullable(&a)
	assert.Equal(t, nullable{dest: &a}, v)
}

type customScanner int

func (*customScanner) Scan(interface{}) error {
	return nil
}

func TestNullable_nullable(t *testing.T) {
	a := customScanner(10)
	v := Nullable(&a)
	assert.Equal(t, &a, v)
}

func TestNullable_ptr(t *testing.T) {
	var a *int
	v := Nullable(&a)
	assert.Equal(t, &a, v)
}

func TestNullable_notPtr(t *testing.T) {
	assert.Panics(t, func() {
		Nullable(0)
	})
}

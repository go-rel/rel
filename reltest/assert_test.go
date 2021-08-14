package reltest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssert_default(t *testing.T) {
	var (
		a = &Assert{}
	)

	assert.False(t, a.assert())
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert())
	assert.True(t, a.call(context.TODO()))
}

func TestAssert_once(t *testing.T) {
	var (
		a = &Assert{}
	)

	a.Once()

	assert.False(t, a.assert())
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert())
	assert.False(t, a.call(context.TODO()))
}

func TestAssert_times(t *testing.T) {
	var (
		a = &Assert{}
	)

	a.Times(2)

	assert.False(t, a.assert())
	assert.True(t, a.call(context.TODO()))
	assert.False(t, a.assert())
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert())
	assert.False(t, a.call(context.TODO()))
}

func TestAssert_maybe(t *testing.T) {
	var (
		a = &Assert{}
	)

	a.Maybe()

	assert.True(t, a.assert())
	assert.True(t, a.call(context.TODO()))
	assert.True(t, a.assert())
}

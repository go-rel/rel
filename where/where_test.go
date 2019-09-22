package where

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterAndOp,
	}, And())
}

func TestOr(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
	}, Or())
}

func TestNot(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterNotOp,
	}, Not())
}

func TestEq(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterEqOp,
		Field: "field",
		Value: "value",
	}, Eq("field", "value"))
}

func TestNe(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNeOp,
		Field: "field",
		Value: "value",
	}, Ne("field", "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterLtOp,
		Field: "field",
		Value: 10,
	}, Lt("field", 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterLteOp,
		Field: "field",
		Value: 10,
	}, Lte("field", 10))
}

func TestFilter_Gt(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterGtOp,
		Field: "field",
		Value: 10,
	}, Gt("field", 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterGteOp,
		Field: "field",
		Value: 10,
	}, Gte("field", 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNilOp,
		Field: "field",
	}, Nil("field"))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNotNilOp,
		Field: "field",
	}, NotNil("field"))
}

func TestIn(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterInOp,
		Field: "field",
		Value: []interface{}{"value1", "value2"},
	}, In("field", "value1", "value2"))
}

func TestNin(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNinOp,
		Field: "field",
		Value: []interface{}{"value1", "value2"},
	}, Nin("field", "value1", "value2"))
}

func TestLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterLikeOp,
		Field: "field",
		Value: "%expr%",
	}, Like("field", "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNotLikeOp,
		Field: "field",
		Value: "%expr%",
	}, NotLike("field", "%expr%"))
}

func TestFragment(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterFragmentOp,
		Field: "expr",
		Value: []interface{}{"value"},
	}, Fragment("expr", "value"))
}

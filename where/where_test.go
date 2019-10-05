package where

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterAndOp,
	}, And())
}

func TestOr(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
	}, Or())
}

func TestNot(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterNotOp,
	}, Not())
}

func TestEq(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterEqOp,
		Field: "field",
		Value: "value",
	}, Eq("field", "value"))
}

func TestNe(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNeOp,
		Field: "field",
		Value: "value",
	}, Ne("field", "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterLtOp,
		Field: "field",
		Value: 10,
	}, Lt("field", 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterLteOp,
		Field: "field",
		Value: 10,
	}, Lte("field", 10))
}

func TestFilter_Gt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterGtOp,
		Field: "field",
		Value: 10,
	}, Gt("field", 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterGteOp,
		Field: "field",
		Value: 10,
	}, Gte("field", 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNilOp,
		Field: "field",
	}, Nil("field"))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNotNilOp,
		Field: "field",
	}, NotNil("field"))
}

func TestIn(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterInOp,
		Field: "field",
		Value: []interface{}{"value1", "value2"},
	}, In("field", "value1", "value2"))
}

func TestNin(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNinOp,
		Field: "field",
		Value: []interface{}{"value1", "value2"},
	}, Nin("field", "value1", "value2"))
}

func TestLike(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterLikeOp,
		Field: "field",
		Value: "%expr%",
	}, Like("field", "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNotLikeOp,
		Field: "field",
		Value: "%expr%",
	}, NotLike("field", "%expr%"))
}

func TestFragment(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterFragmentOp,
		Field: "expr",
		Value: []interface{}{"value"},
	}, Fragment("expr", "value"))
}

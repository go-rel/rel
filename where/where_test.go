package where

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.AndOp,
	}, And())
}

func TestOr(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
	}, Or())
}

func TestNot(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.NotOp,
	}, Not())
}

func TestEq(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.EqOp,
		Field:  "field",
		Values: []interface{}{"value"},
	}, Eq("field", "value"))
}

func TestNe(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.NeOp,
		Field:  "field",
		Values: []interface{}{"value"},
	}, Ne("field", "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.LtOp,
		Field:  "field",
		Values: []interface{}{10},
	}, Lt("field", 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.LteOp,
		Field:  "field",
		Values: []interface{}{10},
	}, Lte("field", 10))
}

func TestFilter_Gt(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.GtOp,
		Field:  "field",
		Values: []interface{}{10},
	}, Gt("field", 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.GteOp,
		Field:  "field",
		Values: []interface{}{10},
	}, Gte("field", 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:  grimoire.NilOp,
		Field: "field",
	}, Nil("field"))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:  grimoire.NotNilOp,
		Field: "field",
	}, NotNil("field"))
}

func TestIn(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.InOp,
		Field:  "field",
		Values: []interface{}{"value1", "value2"},
	}, In("field", "value1", "value2"))
}

func TestNin(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.NinOp,
		Field:  "field",
		Values: []interface{}{"value1", "value2"},
	}, Nin("field", "value1", "value2"))
}

func TestLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.LikeOp,
		Field:  "field",
		Values: []interface{}{"%expr%"},
	}, Like("field", "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.NotLikeOp,
		Field:  "field",
		Values: []interface{}{"%expr%"},
	}, NotLike("field", "%expr%"))
}

func TestFragment(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.FragmentOp,
		Field:  "expr",
		Values: []interface{}{"value"},
	}, Fragment("expr", "value"))
}

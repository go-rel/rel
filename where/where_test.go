package where

import (
	"testing"

	"github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type: query.AndOp,
	}, And())
}

func TestOr(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type: query.OrOp,
	}, Or())
}

func TestNot(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type: query.NotOp,
	}, Not())
}

func TestEq(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.EqOp,
		Column: "field",
		Values: []interface{}{"value"},
	}, Eq("field", "value"))
}

func TestNe(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.NeOp,
		Column: "field",
		Values: []interface{}{"value"},
	}, Ne("field", "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.LtOp,
		Column: "field",
		Values: []interface{}{10},
	}, Lt("field", 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.LteOp,
		Column: "field",
		Values: []interface{}{10},
	}, Lte("field", 10))
}

func TestFilter_Gt(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.GtOp,
		Column: "field",
		Values: []interface{}{10},
	}, Gt("field", 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.GteOp,
		Column: "field",
		Values: []interface{}{10},
	}, Gte("field", 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.NilOp,
		Column: "field",
	}, Nil("field"))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.NotNilOp,
		Column: "field",
	}, NotNil("field"))
}

func TestIn(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.InOp,
		Column: "field",
		Values: []interface{}{"value1", "value2"},
	}, In("field", "value1", "value2"))
}

func TestNin(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.NinOp,
		Column: "field",
		Values: []interface{}{"value1", "value2"},
	}, Nin("field", "value1", "value2"))
}

func TestLike(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.LikeOp,
		Column: "field",
		Values: []interface{}{"%expr%"},
	}, Like("field", "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.NotLikeOp,
		Column: "field",
		Values: []interface{}{"%expr%"},
	}, NotLike("field", "%expr%"))
}

func TestFragment(t *testing.T) {
	assert.Equal(t, query.Filter{
		Type:   query.FragmentOp,
		Column: "expr",
		Values: []interface{}{"value"},
	}, Fragment("expr", "value"))
}

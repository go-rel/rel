package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var result Filter

func BenchmarkFilter_chain1(b *testing.B) {
	var f Filter
	for n := 0; n < b.N; n++ {
		f = FilterEq("id", 1)
	}
	result = f
}

func BenchmarkFilter_chain2(b *testing.B) {
	var f Filter
	for n := 0; n < b.N; n++ {
		f = FilterEq("id", 1).AndNe("name", "foo")
	}
	result = f
}

func BenchmarkFilter_chain3(b *testing.B) {
	var f Filter
	for n := 0; n < b.N; n++ {
		f = FilterEq("id", 1).AndNe("name", "foo").AndGt("score", 80)
	}
	result = f
}

func BenchmarkFilter_chain4(b *testing.B) {
	var f Filter
	for n := 0; n < b.N; n++ {
		f = FilterEq("id", 1).AndNe("name", "foo").AndGt("score", 80).AndLt("avg", 10)
	}
	result = f
}

func BenchmarkFilter_slice1(b *testing.B) {
	var f Filter
	for n := 0; n < b.N; n++ {
		f = FilterAnd(FilterEq("id", 1))
	}
	result = f
}

func BenchmarkFilter_slice2(b *testing.B) {
	var f Filter
	for n := 0; n < b.N; n++ {
		f = FilterAnd(FilterEq("id", 1), FilterNe("name", "foo"))
	}
	result = f
}

func BenchmarkFilter_slice3(b *testing.B) {
	var f Filter
	for n := 0; n < b.N; n++ {
		f = FilterAnd(FilterEq("id", 1), FilterNe("name", "foo"), FilterGt("score", 80))
	}
	result = f
}

func BenchmarkFilter_slice4(b *testing.B) {
	var f Filter
	for n := 0; n < b.N; n++ {
		f = FilterAnd(FilterEq("id", 1), FilterNe("name", "foo"), FilterGt("score", 80), FilterLt("avg", 10))
	}
	result = f
}

var filter1 = FilterEq("id", 1)
var filter2 = FilterNe("name", "foo")
var filter3 = FilterGt("score", 80)
var filter4 = FilterLt("avg", 10)

func TestFilter_None(t *testing.T) {
	assert.True(t, Filter{}.None())
	assert.True(t, FilterAnd().None())
	assert.True(t, FilterNot().None())

	assert.False(t, FilterAnd(filter1).None())
	assert.False(t, FilterAnd(filter1, filter2).None())
	assert.False(t, filter1.None())
}

func TestFilter_And(t *testing.T) {
	tests := []struct {
		Case      string
		Operation Filter
		Result    Filter
	}{
		{
			`Filter{}.And()`,
			Filter{}.And(),
			FilterAnd(),
		},
		{
			`Filter{}.And(filter1)`,
			Filter{}.And(filter1),
			filter1,
		},
		{
			`Filter{}.And(filter1).And()`,
			Filter{}.And(filter1).And(),
			filter1,
		},
		{
			`Filter{}.And(filter1, filter2)`,
			Filter{}.And(filter1, filter2),
			FilterAnd(filter1, filter2),
		},
		{
			`Filter{}.And(filter1, filter2).And()`,
			Filter{}.And(filter1, filter2).And(),
			FilterAnd(filter1, filter2),
		},
		{
			`Filter{}.And(filter1, filter2, filter3)`,
			Filter{}.And(filter1, filter2, filter3),
			FilterAnd(filter1, filter2, filter3),
		},
		{
			`Filter{}.And(filter1, filter2, filter3).And()`,
			Filter{}.And(filter1, filter2, filter3).And(),
			FilterAnd(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2)`,
			filter1.And(filter2),
			FilterAnd(filter1, filter2),
		},
		{
			`filter1.And(filter2).And()`,
			filter1.And(filter2).And(),
			FilterAnd(filter1, filter2),
		},
		{
			`filter1.And(filter2).And(filter3)`,
			filter1.And(filter2).And(filter3),
			FilterAnd(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2).And(filter3).And()`,
			filter1.And(filter2).And(filter3).And(),
			FilterAnd(filter1, filter2, filter3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestFilter_Or(t *testing.T) {
	tests := []struct {
		Case      string
		Operation Filter
		Result    Filter
	}{
		{
			`Filter{}.Or()`,
			Filter{}.Or(),
			FilterOr(),
		},
		{
			`Filter{}.Or(filter1)`,
			Filter{}.Or(filter1),
			filter1,
		},
		{
			`Filter{}.Or(filter1).Or()`,
			Filter{}.Or(filter1).Or(),
			filter1,
		},
		{
			`Filter{}.Or(filter1, filter2)`,
			Filter{}.Or(filter1, filter2),
			FilterOr(filter1, filter2),
		},
		{
			`Filter{}.Or(filter1, filter2).Or()`,
			Filter{}.Or(filter1, filter2).Or(),
			FilterOr(filter1, filter2),
		},
		{
			`Filter{}.Or(filter1, filter2, filter3)`,
			Filter{}.Or(filter1, filter2, filter3),
			FilterOr(filter1, filter2, filter3),
		},
		{
			`Filter{}.Or(filter1, filter2, filter3).Or()`,
			Filter{}.Or(filter1, filter2, filter3).Or(),
			FilterOr(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2)`,
			filter1.Or(filter2),
			FilterOr(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or()`,
			filter1.Or(filter2).Or(),
			FilterOr(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or(filter3)`,
			filter1.Or(filter2).Or(filter3),
			FilterOr(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2).Or(filter3).Or()`,
			filter1.Or(filter2).Or(filter3).Or(),
			FilterOr(filter1, filter2, filter3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestFilterAnd(t *testing.T) {
	tests := []struct {
		Case      string
		Operation Filter
		Result    Filter
	}{
		{
			`FilterAnd()`,
			FilterAnd(),
			Filter{Type: AndOp},
		},
		{
			`FilterAnd(filter1)`,
			FilterAnd(filter1),
			filter1,
		},
		{
			`FilterAnd(filter1, filter2)`,
			FilterAnd(filter1, filter2),
			Filter{
				Type:  AndOp,
				Inner: []Filter{filter1, filter2},
			},
		},
		{
			`FilterAnd(filter1, FilterOr(filter2, filter3))`,
			FilterAnd(filter1, FilterOr(filter2, filter3)),
			Filter{
				Type: AndOp,
				Inner: []Filter{
					filter1,
					{
						Type:  OrOp,
						Inner: []Filter{filter2, filter3},
					},
				},
			},
		},
		{
			`FilterAnd(FilterOr(filter1, filter2), filter3)`,
			FilterAnd(FilterOr(filter1, filter2), filter3),
			Filter{
				Type: AndOp,
				Inner: []Filter{
					{
						Type:  OrOp,
						Inner: []Filter{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`FilterAnd(FilterOr(filter1, filter2), FilterOr(filter3, filter4))`,
			FilterAnd(FilterOr(filter1, filter2), FilterOr(filter3, filter4)),
			Filter{
				Type: AndOp,
				Inner: []Filter{
					{
						Type:  OrOp,
						Inner: []Filter{filter1, filter2},
					},
					{
						Type:  OrOp,
						Inner: []Filter{filter3, filter4},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestFilterOr(t *testing.T) {
	tests := []struct {
		Case      string
		Operation Filter
		Result    Filter
	}{
		{
			`FilterOr()`,
			FilterOr(),
			Filter{Type: OrOp},
		},
		{
			`FilterOr(filter1)`,
			FilterOr(filter1),
			filter1,
		},
		{
			`FilterOr(filter1, filter2)`,
			FilterOr(filter1, filter2),
			Filter{
				Type:  OrOp,
				Inner: []Filter{filter1, filter2},
			},
		},
		{
			`FilterOr(filter1, FilterAnd(filter2, filter3))`,
			FilterOr(filter1, FilterAnd(filter2, filter3)),
			Filter{
				Type: OrOp,
				Inner: []Filter{
					filter1,
					{
						Type:  AndOp,
						Inner: []Filter{filter2, filter3},
					},
				},
			},
		},
		{
			`FilterOr(FilterAnd(filter1, filter2), filter3)`,
			FilterOr(FilterAnd(filter1, filter2), filter3),
			Filter{
				Type: OrOp,
				Inner: []Filter{
					{
						Type:  AndOp,
						Inner: []Filter{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`FilterOr(FilterAnd(filter1, filter2), FilterAnd(filter3, filter4))`,
			FilterOr(FilterAnd(filter1, filter2), FilterAnd(filter3, filter4)),
			Filter{
				Type: OrOp,
				Inner: []Filter{
					{
						Type:  AndOp,
						Inner: []Filter{filter1, filter2},
					},
					{
						Type:  AndOp,
						Inner: []Filter{filter3, filter4},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestFilter_Not(t *testing.T) {
	tests := []struct {
		Case     string
		Input    FilterOp
		Expected FilterOp
	}{
		{
			`Not Eq`,
			EqOp,
			NeOp,
		},
		{
			`Not Lt`,
			LtOp,
			GteOp,
		},
		{
			`Not Lte`,
			LteOp,
			GtOp,
		},
		{
			`Not Gt`,
			GtOp,
			LteOp,
		},
		{
			`Not Gte`,
			GteOp,
			LtOp,
		},
		{
			`Not Nil`,
			NilOp,
			NotNilOp,
		},
		{
			`Not In`,
			InOp,
			NinOp,
		},
		{
			`Not Like`,
			LikeOp,
			NotLikeOp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, FilterNot(Filter{Type: tt.Input}).Type)
		})
	}
}

func TestFilter_AndEq(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   EqOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, Filter{}.AndEq("field", "value"))
}

func TestFilter_AndNe(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   NeOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, Filter{}.AndNe("field", "value"))
}

func TestFilter_AndLt(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   LtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, Filter{}.AndLt("field", 10))
}

func TestFilter_AndLte(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   LteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, Filter{}.AndLte("field", 10))
}

func TestFilter_AndFilter_Gt(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   GtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, Filter{}.AndGt("field", 10))
}

func TestFilter_AndGte(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   GteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, Filter{}.AndGte("field", 10))
}

func TestFilter_AndNil(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:  NilOp,
				Field: "field",
			},
		},
	}, Filter{}.AndNil("field"))
}

func TestFilter_AndNotNil(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:  NotNilOp,
				Field: "field",
			},
		},
	}, Filter{}.AndNotNil("field"))
}

func TestFilter_AndIn(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   InOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, Filter{}.AndIn("field", "value1", "value2"))
}

func TestFilter_AndNin(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   NinOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, Filter{}.AndNin("field", "value1", "value2"))
}

func TestFilter_AndLike(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   LikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, Filter{}.AndLike("field", "%expr%"))
}

func TestFilter_AndNotLike(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   NotLikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, Filter{}.AndNotLike("field", "%expr%"))
}

func TestFilter_AndFragment(t *testing.T) {
	assert.Equal(t, Filter{
		Inner: []Filter{
			{
				Type:   FragmentOp,
				Field:  "expr",
				Values: []interface{}{"value"},
			},
		},
	}, Filter{}.AndFragment("expr", "value"))
}

func TestFilter_OrEq(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   EqOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, Filter{}.OrEq("field", "value"))
}

func TestFilter_OrNe(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   NeOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, Filter{}.OrNe("field", "value"))
}

func TestFilter_OrLt(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   LtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, Filter{}.OrLt("field", 10))
}

func TestFilter_OrLte(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   LteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, Filter{}.OrLte("field", 10))
}

func TestFilter_OrFilter_Gt(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   GtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, Filter{}.OrGt("field", 10))
}

func TestFilter_OrGte(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   GteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, Filter{}.OrGte("field", 10))
}

func TestFilter_OrNil(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:  NilOp,
				Field: "field",
			},
		},
	}, Filter{}.OrNil("field"))
}

func TestFilter_OrNotNil(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:  NotNilOp,
				Field: "field",
			},
		},
	}, Filter{}.OrNotNil("field"))
}

func TestFilter_OrIn(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   InOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, Filter{}.OrIn("field", "value1", "value2"))
}

func TestFilter_OrNin(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   NinOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, Filter{}.OrNin("field", "value1", "value2"))
}

func TestFilter_OrLike(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   LikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, Filter{}.OrLike("field", "%expr%"))
}

func TestFilter_OrNotLike(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   NotLikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, Filter{}.OrNotLike("field", "%expr%"))
}

func TestFilter_OrFragment(t *testing.T) {
	assert.Equal(t, Filter{
		Type: OrOp,
		Inner: []Filter{
			{
				Type:   FragmentOp,
				Field:  "expr",
				Values: []interface{}{"value"},
			},
		},
	}, Filter{}.OrFragment("expr", "value"))
}

func TestFilterEq(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   EqOp,
		Field:  "field",
		Values: []interface{}{"value"},
	}, FilterEq("field", "value"))
}

func TestFilterNe(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   NeOp,
		Field:  "field",
		Values: []interface{}{"value"},
	}, FilterNe("field", "value"))
}

func TestFilterLt(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   LtOp,
		Field:  "field",
		Values: []interface{}{10},
	}, FilterLt("field", 10))
}

func TestFilterLte(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   LteOp,
		Field:  "field",
		Values: []interface{}{10},
	}, FilterLte("field", 10))
}

func TestFilterFilter_Gt(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   GtOp,
		Field:  "field",
		Values: []interface{}{10},
	}, FilterGt("field", 10))
}

func TestFilterGte(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   GteOp,
		Field:  "field",
		Values: []interface{}{10},
	}, FilterGte("field", 10))
}

func TestFilterNil(t *testing.T) {
	assert.Equal(t, Filter{
		Type:  NilOp,
		Field: "field",
	}, FilterNil("field"))
}

func TestFilterNotNil(t *testing.T) {
	assert.Equal(t, Filter{
		Type:  NotNilOp,
		Field: "field",
	}, FilterNotNil("field"))
}

func TestFilterIn(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   InOp,
		Field:  "field",
		Values: []interface{}{"value1", "value2"},
	}, FilterIn("field", "value1", "value2"))
}

func TestFilterNin(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   NinOp,
		Field:  "field",
		Values: []interface{}{"value1", "value2"},
	}, FilterNin("field", "value1", "value2"))
}

func TestFilterLike(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   LikeOp,
		Field:  "field",
		Values: []interface{}{"%expr%"},
	}, FilterLike("field", "%expr%"))
}

func TestFilterNotLike(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   NotLikeOp,
		Field:  "field",
		Values: []interface{}{"%expr%"},
	}, FilterNotLike("field", "%expr%"))
}

func TestFilterFragment(t *testing.T) {
	assert.Equal(t, Filter{
		Type:   FragmentOp,
		Field:  "expr",
		Values: []interface{}{"value"},
	}, FilterFragment("expr", "value"))
}

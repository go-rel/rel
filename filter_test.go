package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

var result grimoire.FilterClause

func BenchmarkFilterClause_chain1(b *testing.B) {
	var f grimoire.FilterClause
	for n := 0; n < b.N; n++ {
		f = grimoire.FilterEq("id", 1)
	}
	result = f
}

func BenchmarkFilterClause_chain2(b *testing.B) {
	var f grimoire.FilterClause
	for n := 0; n < b.N; n++ {
		f = grimoire.FilterEq("id", 1).AndNe("name", "foo")
	}
	result = f
}

func BenchmarkFilterClause_chain3(b *testing.B) {
	var f grimoire.FilterClause
	for n := 0; n < b.N; n++ {
		f = grimoire.FilterEq("id", 1).AndNe("name", "foo").AndGt("score", 80)
	}
	result = f
}

func BenchmarkFilterClause_chain4(b *testing.B) {
	var f grimoire.FilterClause
	for n := 0; n < b.N; n++ {
		f = grimoire.FilterEq("id", 1).AndNe("name", "foo").AndGt("score", 80).AndLt("avg", 10)
	}
	result = f
}

func BenchmarkFilterClause_slice1(b *testing.B) {
	var f grimoire.FilterClause
	for n := 0; n < b.N; n++ {
		f = grimoire.FilterAnd(grimoire.FilterEq("id", 1))
	}
	result = f
}

func BenchmarkFilterClause_slice2(b *testing.B) {
	var f grimoire.FilterClause
	for n := 0; n < b.N; n++ {
		f = grimoire.FilterAnd(grimoire.FilterEq("id", 1), grimoire.FilterNe("name", "foo"))
	}
	result = f
}

func BenchmarkFilterClause_slice3(b *testing.B) {
	var f grimoire.FilterClause
	for n := 0; n < b.N; n++ {
		f = grimoire.FilterAnd(grimoire.FilterEq("id", 1), grimoire.FilterNe("name", "foo"), grimoire.FilterGt("score", 80))
	}
	result = f
}

func BenchmarkFilterClause_slice4(b *testing.B) {
	var f grimoire.FilterClause
	for n := 0; n < b.N; n++ {
		f = grimoire.FilterAnd(grimoire.FilterEq("id", 1), grimoire.FilterNe("name", "foo"), grimoire.FilterGt("score", 80), grimoire.FilterLt("avg", 10))
	}
	result = f
}

var filter1 = grimoire.FilterEq("id", 1)
var filter2 = grimoire.FilterNe("name", "foo")
var filter3 = grimoire.FilterGt("score", 80)
var filter4 = grimoire.FilterLt("avg", 10)

func TestFilterClause_None(t *testing.T) {
	assert.True(t, grimoire.FilterClause{}.None())
	assert.True(t, grimoire.FilterAnd().None())
	assert.True(t, grimoire.FilterNot().None())

	assert.False(t, grimoire.FilterAnd(filter1).None())
	assert.False(t, grimoire.FilterAnd(filter1, filter2).None())
	assert.False(t, filter1.None())
}

func TestFilterClause_And(t *testing.T) {
	tests := []struct {
		Case      string
		Operation grimoire.FilterClause
		Result    grimoire.FilterClause
	}{
		{
			`grimoire.FilterClause{}.And()`,
			grimoire.FilterClause{}.And(),
			grimoire.FilterAnd(),
		},
		{
			`grimoire.FilterClause{}.And(filter1)`,
			grimoire.FilterClause{}.And(filter1),
			filter1,
		},
		{
			`grimoire.FilterClause{}.And(filter1).And()`,
			grimoire.FilterClause{}.And(filter1).And(),
			filter1,
		},
		{
			`grimoire.FilterClause{}.And(filter1, filter2)`,
			grimoire.FilterClause{}.And(filter1, filter2),
			grimoire.FilterAnd(filter1, filter2),
		},
		{
			`grimoire.FilterClause{}.And(filter1, filter2).And()`,
			grimoire.FilterClause{}.And(filter1, filter2).And(),
			grimoire.FilterAnd(filter1, filter2),
		},
		{
			`grimoire.FilterClause{}.And(filter1, filter2, filter3)`,
			grimoire.FilterClause{}.And(filter1, filter2, filter3),
			grimoire.FilterAnd(filter1, filter2, filter3),
		},
		{
			`grimoire.FilterClause{}.And(filter1, filter2, filter3).And()`,
			grimoire.FilterClause{}.And(filter1, filter2, filter3).And(),
			grimoire.FilterAnd(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2)`,
			filter1.And(filter2),
			grimoire.FilterAnd(filter1, filter2),
		},
		{
			`filter1.And(filter2).And()`,
			filter1.And(filter2).And(),
			grimoire.FilterAnd(filter1, filter2),
		},
		{
			`filter1.And(filter2).And(filter3)`,
			filter1.And(filter2).And(filter3),
			grimoire.FilterAnd(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2).And(filter3).And()`,
			filter1.And(filter2).And(filter3).And(),
			grimoire.FilterAnd(filter1, filter2, filter3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestFilterClause_Or(t *testing.T) {
	tests := []struct {
		Case      string
		Operation grimoire.FilterClause
		Result    grimoire.FilterClause
	}{
		{
			`grimoire.FilterClause{}.Or()`,
			grimoire.FilterClause{}.Or(),
			grimoire.FilterOr(),
		},
		{
			`grimoire.FilterClause{}.Or(filter1)`,
			grimoire.FilterClause{}.Or(filter1),
			filter1,
		},
		{
			`grimoire.FilterClause{}.Or(filter1).Or()`,
			grimoire.FilterClause{}.Or(filter1).Or(),
			filter1,
		},
		{
			`grimoire.FilterClause{}.Or(filter1, filter2)`,
			grimoire.FilterClause{}.Or(filter1, filter2),
			grimoire.FilterOr(filter1, filter2),
		},
		{
			`grimoire.FilterClause{}.Or(filter1, filter2).Or()`,
			grimoire.FilterClause{}.Or(filter1, filter2).Or(),
			grimoire.FilterOr(filter1, filter2),
		},
		{
			`grimoire.FilterClause{}.Or(filter1, filter2, filter3)`,
			grimoire.FilterClause{}.Or(filter1, filter2, filter3),
			grimoire.FilterOr(filter1, filter2, filter3),
		},
		{
			`grimoire.FilterClause{}.Or(filter1, filter2, filter3).Or()`,
			grimoire.FilterClause{}.Or(filter1, filter2, filter3).Or(),
			grimoire.FilterOr(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2)`,
			filter1.Or(filter2),
			grimoire.FilterOr(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or()`,
			filter1.Or(filter2).Or(),
			grimoire.FilterOr(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or(filter3)`,
			filter1.Or(filter2).Or(filter3),
			grimoire.FilterOr(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2).Or(filter3).Or()`,
			filter1.Or(filter2).Or(filter3).Or(),
			grimoire.FilterOr(filter1, filter2, filter3),
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
		Operation grimoire.FilterClause
		Result    grimoire.FilterClause
	}{
		{
			`grimoire.FilterAnd()`,
			grimoire.FilterAnd(),
			grimoire.FilterClause{Type: grimoire.AndOp},
		},
		{
			`grimoire.FilterAnd(filter1)`,
			grimoire.FilterAnd(filter1),
			filter1,
		},
		{
			`grimoire.FilterAnd(filter1, filter2)`,
			grimoire.FilterAnd(filter1, filter2),
			grimoire.FilterClause{
				Type:  grimoire.AndOp,
				Inner: []grimoire.FilterClause{filter1, filter2},
			},
		},
		{
			`grimoire.FilterAnd(filter1, grimoire.FilterOr(filter2, filter3))`,
			grimoire.FilterAnd(filter1, grimoire.FilterOr(filter2, filter3)),
			grimoire.FilterClause{
				Type: grimoire.AndOp,
				Inner: []grimoire.FilterClause{
					filter1,
					{
						Type:  grimoire.OrOp,
						Inner: []grimoire.FilterClause{filter2, filter3},
					},
				},
			},
		},
		{
			`grimoire.FilterAnd(grimoire.FilterOr(filter1, filter2), filter3)`,
			grimoire.FilterAnd(grimoire.FilterOr(filter1, filter2), filter3),
			grimoire.FilterClause{
				Type: grimoire.AndOp,
				Inner: []grimoire.FilterClause{
					{
						Type:  grimoire.OrOp,
						Inner: []grimoire.FilterClause{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`grimoire.FilterAnd(grimoire.FilterOr(filter1, filter2), grimoire.FilterOr(filter3, filter4))`,
			grimoire.FilterAnd(grimoire.FilterOr(filter1, filter2), grimoire.FilterOr(filter3, filter4)),
			grimoire.FilterClause{
				Type: grimoire.AndOp,
				Inner: []grimoire.FilterClause{
					{
						Type:  grimoire.OrOp,
						Inner: []grimoire.FilterClause{filter1, filter2},
					},
					{
						Type:  grimoire.OrOp,
						Inner: []grimoire.FilterClause{filter3, filter4},
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
		Operation grimoire.FilterClause
		Result    grimoire.FilterClause
	}{
		{
			`grimoire.FilterOr()`,
			grimoire.FilterOr(),
			grimoire.FilterClause{Type: grimoire.OrOp},
		},
		{
			`grimoire.FilterOr(filter1)`,
			grimoire.FilterOr(filter1),
			filter1,
		},
		{
			`grimoire.FilterOr(filter1, filter2)`,
			grimoire.FilterOr(filter1, filter2),
			grimoire.FilterClause{
				Type:  grimoire.OrOp,
				Inner: []grimoire.FilterClause{filter1, filter2},
			},
		},
		{
			`grimoire.FilterOr(filter1, grimoire.FilterAnd(filter2, filter3))`,
			grimoire.FilterOr(filter1, grimoire.FilterAnd(filter2, filter3)),
			grimoire.FilterClause{
				Type: grimoire.OrOp,
				Inner: []grimoire.FilterClause{
					filter1,
					{
						Type:  grimoire.AndOp,
						Inner: []grimoire.FilterClause{filter2, filter3},
					},
				},
			},
		},
		{
			`grimoire.FilterOr(grimoire.FilterAnd(filter1, filter2), filter3)`,
			grimoire.FilterOr(grimoire.FilterAnd(filter1, filter2), filter3),
			grimoire.FilterClause{
				Type: grimoire.OrOp,
				Inner: []grimoire.FilterClause{
					{
						Type:  grimoire.AndOp,
						Inner: []grimoire.FilterClause{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`grimoire.FilterOr(grimoire.FilterAnd(filter1, filter2), grimoire.FilterAnd(filter3, filter4))`,
			grimoire.FilterOr(grimoire.FilterAnd(filter1, filter2), grimoire.FilterAnd(filter3, filter4)),
			grimoire.FilterClause{
				Type: grimoire.OrOp,
				Inner: []grimoire.FilterClause{
					{
						Type:  grimoire.AndOp,
						Inner: []grimoire.FilterClause{filter1, filter2},
					},
					{
						Type:  grimoire.AndOp,
						Inner: []grimoire.FilterClause{filter3, filter4},
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

func TestFilterClause_Not(t *testing.T) {
	tests := []struct {
		Case     string
		Input    grimoire.FilterOp
		Expected grimoire.FilterOp
	}{
		{
			`Not Eq`,
			grimoire.EqOp,
			grimoire.NeOp,
		},
		{
			`Not Lt`,
			grimoire.LtOp,
			grimoire.GteOp,
		},
		{
			`Not Lte`,
			grimoire.LteOp,
			grimoire.GtOp,
		},
		{
			`Not Gt`,
			grimoire.GtOp,
			grimoire.LteOp,
		},
		{
			`Not Gte`,
			grimoire.GteOp,
			grimoire.LtOp,
		},
		{
			`Not Nil`,
			grimoire.NilOp,
			grimoire.NotNilOp,
		},
		{
			`Not In`,
			grimoire.InOp,
			grimoire.NinOp,
		},
		{
			`Not Like`,
			grimoire.LikeOp,
			grimoire.NotLikeOp,
		},
		{
			`And Op`,
			grimoire.AndOp,
			grimoire.NotOp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, grimoire.FilterNot(grimoire.FilterClause{Type: tt.Input}).Type)
		})
	}
}

func TestFilterClause_AndEq(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.EqOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, grimoire.FilterClause{}.AndEq("field", "value"))
}

func TestFilterClause_AndNe(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.NeOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, grimoire.FilterClause{}.AndNe("field", "value"))
}

func TestFilterClause_AndLt(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.LtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, grimoire.FilterClause{}.AndLt("field", 10))
}

func TestFilterClause_AndLte(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.LteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, grimoire.FilterClause{}.AndLte("field", 10))
}

func TestFilterClause_AndGt(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.GtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, grimoire.FilterClause{}.AndGt("field", 10))
}

func TestFilterClause_AndGte(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.GteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, grimoire.FilterClause{}.AndGte("field", 10))
}

func TestFilterClause_AndNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:  grimoire.NilOp,
				Field: "field",
			},
		},
	}, grimoire.FilterClause{}.AndNil("field"))
}

func TestFilterClause_AndNotNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:  grimoire.NotNilOp,
				Field: "field",
			},
		},
	}, grimoire.FilterClause{}.AndNotNil("field"))
}

func TestFilterClause_AndIn(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.InOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, grimoire.FilterClause{}.AndIn("field", "value1", "value2"))
}

func TestFilterClause_AndNin(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.NinOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, grimoire.FilterClause{}.AndNin("field", "value1", "value2"))
}

func TestFilterClause_AndLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.LikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, grimoire.FilterClause{}.AndLike("field", "%expr%"))
}

func TestFilterClause_AndNotLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.NotLikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, grimoire.FilterClause{}.AndNotLike("field", "%expr%"))
}

func TestFilterClause_AndFragment(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.FragmentOp,
				Field:  "expr",
				Values: []interface{}{"value"},
			},
		},
	}, grimoire.FilterClause{}.AndFragment("expr", "value"))
}

func TestFilterClause_OrEq(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.EqOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, grimoire.FilterClause{}.OrEq("field", "value"))
}

func TestFilterClause_OrNe(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.NeOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, grimoire.FilterClause{}.OrNe("field", "value"))
}

func TestFilterClause_OrLt(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.LtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, grimoire.FilterClause{}.OrLt("field", 10))
}

func TestFilterClause_OrLte(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.LteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, grimoire.FilterClause{}.OrLte("field", 10))
}

func TestFilterClause_OrGt(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.GtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, grimoire.FilterClause{}.OrGt("field", 10))
}

func TestFilterClause_OrGte(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.GteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, grimoire.FilterClause{}.OrGte("field", 10))
}

func TestFilterClause_OrNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:  grimoire.NilOp,
				Field: "field",
			},
		},
	}, grimoire.FilterClause{}.OrNil("field"))
}

func TestFilterClause_OrNotNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:  grimoire.NotNilOp,
				Field: "field",
			},
		},
	}, grimoire.FilterClause{}.OrNotNil("field"))
}

func TestFilterClause_OrIn(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.InOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, grimoire.FilterClause{}.OrIn("field", "value1", "value2"))
}

func TestFilterClause_OrNin(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.NinOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, grimoire.FilterClause{}.OrNin("field", "value1", "value2"))
}

func TestFilterClause_OrLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.LikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, grimoire.FilterClause{}.OrLike("field", "%expr%"))
}

func TestFilterClause_OrNotLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.NotLikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, grimoire.FilterClause{}.OrNotLike("field", "%expr%"))
}

func TestFilterClause_OrFragment(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type: grimoire.OrOp,
		Inner: []grimoire.FilterClause{
			{
				Type:   grimoire.FragmentOp,
				Field:  "expr",
				Values: []interface{}{"value"},
			},
		},
	}, grimoire.FilterClause{}.OrFragment("expr", "value"))
}

func TestFilterEq(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.EqOp,
		Field:  "field",
		Values: []interface{}{"value"},
	}, grimoire.FilterEq("field", "value"))
}

func FilterNe(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.NeOp,
		Field:  "field",
		Values: []interface{}{"value"},
	}, grimoire.FilterNe("field", "value"))
}

func TestFilterLt(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.LtOp,
		Field:  "field",
		Values: []interface{}{10},
	}, grimoire.FilterLt("field", 10))
}

func TestFilterLte(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.LteOp,
		Field:  "field",
		Values: []interface{}{10},
	}, grimoire.FilterLte("field", 10))
}

func TestFilterClauseGt(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.GtOp,
		Field:  "field",
		Values: []interface{}{10},
	}, grimoire.FilterGt("field", 10))
}

func TestFilterGte(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.GteOp,
		Field:  "field",
		Values: []interface{}{10},
	}, grimoire.FilterGte("field", 10))
}

func TestFilterNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:  grimoire.NilOp,
		Field: "field",
	}, grimoire.FilterNil("field"))
}

func TestFilterNotNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:  grimoire.NotNilOp,
		Field: "field",
	}, grimoire.FilterNotNil("field"))
}

func TestFilterIn(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.InOp,
		Field:  "field",
		Values: []interface{}{"value1", "value2"},
	}, grimoire.FilterIn("field", "value1", "value2"))
}

func TestFilterNin(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.NinOp,
		Field:  "field",
		Values: []interface{}{"value1", "value2"},
	}, grimoire.FilterNin("field", "value1", "value2"))
}

func TestFilterLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.LikeOp,
		Field:  "field",
		Values: []interface{}{"%expr%"},
	}, grimoire.FilterLike("field", "%expr%"))
}

func TestFilterNotLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.NotLikeOp,
		Field:  "field",
		Values: []interface{}{"%expr%"},
	}, grimoire.FilterNotLike("field", "%expr%"))
}

func TestFilterFragment(t *testing.T) {
	assert.Equal(t, grimoire.FilterClause{
		Type:   grimoire.FragmentOp,
		Field:  "expr",
		Values: []interface{}{"value"},
	}, grimoire.FilterFragment("expr", "value"))
}

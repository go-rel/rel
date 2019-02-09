package query_test

import (
	"testing"

	"github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

var result query.FilterClause

func BenchmarkFilterClause_chain1(b *testing.B) {
	var f query.FilterClause
	for n := 0; n < b.N; n++ {
		f = query.FilterEq("id", 1)
	}
	result = f
}

func BenchmarkFilterClause_chain2(b *testing.B) {
	var f query.FilterClause
	for n := 0; n < b.N; n++ {
		f = query.FilterEq("id", 1).AndNe("name", "foo")
	}
	result = f
}

func BenchmarkFilterClause_chain3(b *testing.B) {
	var f query.FilterClause
	for n := 0; n < b.N; n++ {
		f = query.FilterEq("id", 1).AndNe("name", "foo").AndGt("score", 80)
	}
	result = f
}

func BenchmarkFilterClause_chain4(b *testing.B) {
	var f query.FilterClause
	for n := 0; n < b.N; n++ {
		f = query.FilterEq("id", 1).AndNe("name", "foo").AndGt("score", 80).AndLt("avg", 10)
	}
	result = f
}

func BenchmarkFilterClause_slice1(b *testing.B) {
	var f query.FilterClause
	for n := 0; n < b.N; n++ {
		f = query.FilterAnd(query.FilterEq("id", 1))
	}
	result = f
}

func BenchmarkFilterClause_slice2(b *testing.B) {
	var f query.FilterClause
	for n := 0; n < b.N; n++ {
		f = query.FilterAnd(query.FilterEq("id", 1), query.FilterNe("name", "foo"))
	}
	result = f
}

func BenchmarkFilterClause_slice3(b *testing.B) {
	var f query.FilterClause
	for n := 0; n < b.N; n++ {
		f = query.FilterAnd(query.FilterEq("id", 1), query.FilterNe("name", "foo"), query.FilterGt("score", 80))
	}
	result = f
}

func BenchmarkFilterClause_slice4(b *testing.B) {
	var f query.FilterClause
	for n := 0; n < b.N; n++ {
		f = query.FilterAnd(query.FilterEq("id", 1), query.FilterNe("name", "foo"), query.FilterGt("score", 80), query.FilterLt("avg", 10))
	}
	result = f
}

var filter1 = query.FilterEq("id", 1)
var filter2 = query.FilterNe("name", "foo")
var filter3 = query.FilterGt("score", 80)
var filter4 = query.FilterLt("avg", 10)

func TestFilterClause_None(t *testing.T) {
	assert.True(t, query.FilterClause{}.None())
	assert.True(t, query.FilterAnd().None())
	assert.True(t, query.FilterNot().None())

	assert.False(t, query.FilterAnd(filter1).None())
	assert.False(t, query.FilterAnd(filter1, filter2).None())
	assert.False(t, filter1.None())
}

func TestFilterClause_And(t *testing.T) {
	tests := []struct {
		Case      string
		Operation query.FilterClause
		Result    query.FilterClause
	}{
		{
			`query.FilterClause{}.And()`,
			query.FilterClause{}.And(),
			query.FilterAnd(),
		},
		{
			`query.FilterClause{}.And(filter1)`,
			query.FilterClause{}.And(filter1),
			filter1,
		},
		{
			`query.FilterClause{}.And(filter1).And()`,
			query.FilterClause{}.And(filter1).And(),
			filter1,
		},
		{
			`query.FilterClause{}.And(filter1, filter2)`,
			query.FilterClause{}.And(filter1, filter2),
			query.FilterAnd(filter1, filter2),
		},
		{
			`query.FilterClause{}.And(filter1, filter2).And()`,
			query.FilterClause{}.And(filter1, filter2).And(),
			query.FilterAnd(filter1, filter2),
		},
		{
			`query.FilterClause{}.And(filter1, filter2, filter3)`,
			query.FilterClause{}.And(filter1, filter2, filter3),
			query.FilterAnd(filter1, filter2, filter3),
		},
		{
			`query.FilterClause{}.And(filter1, filter2, filter3).And()`,
			query.FilterClause{}.And(filter1, filter2, filter3).And(),
			query.FilterAnd(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2)`,
			filter1.And(filter2),
			query.FilterAnd(filter1, filter2),
		},
		{
			`filter1.And(filter2).And()`,
			filter1.And(filter2).And(),
			query.FilterAnd(filter1, filter2),
		},
		{
			`filter1.And(filter2).And(filter3)`,
			filter1.And(filter2).And(filter3),
			query.FilterAnd(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2).And(filter3).And()`,
			filter1.And(filter2).And(filter3).And(),
			query.FilterAnd(filter1, filter2, filter3),
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
		Operation query.FilterClause
		Result    query.FilterClause
	}{
		{
			`query.FilterClause{}.Or()`,
			query.FilterClause{}.Or(),
			query.FilterOr(),
		},
		{
			`query.FilterClause{}.Or(filter1)`,
			query.FilterClause{}.Or(filter1),
			filter1,
		},
		{
			`query.FilterClause{}.Or(filter1).Or()`,
			query.FilterClause{}.Or(filter1).Or(),
			filter1,
		},
		{
			`query.FilterClause{}.Or(filter1, filter2)`,
			query.FilterClause{}.Or(filter1, filter2),
			query.FilterOr(filter1, filter2),
		},
		{
			`query.FilterClause{}.Or(filter1, filter2).Or()`,
			query.FilterClause{}.Or(filter1, filter2).Or(),
			query.FilterOr(filter1, filter2),
		},
		{
			`query.FilterClause{}.Or(filter1, filter2, filter3)`,
			query.FilterClause{}.Or(filter1, filter2, filter3),
			query.FilterOr(filter1, filter2, filter3),
		},
		{
			`query.FilterClause{}.Or(filter1, filter2, filter3).Or()`,
			query.FilterClause{}.Or(filter1, filter2, filter3).Or(),
			query.FilterOr(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2)`,
			filter1.Or(filter2),
			query.FilterOr(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or()`,
			filter1.Or(filter2).Or(),
			query.FilterOr(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or(filter3)`,
			filter1.Or(filter2).Or(filter3),
			query.FilterOr(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2).Or(filter3).Or()`,
			filter1.Or(filter2).Or(filter3).Or(),
			query.FilterOr(filter1, filter2, filter3),
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
		Operation query.FilterClause
		Result    query.FilterClause
	}{
		{
			`query.FilterAnd()`,
			query.FilterAnd(),
			query.FilterClause{Type: query.AndOp},
		},
		{
			`query.FilterAnd(filter1)`,
			query.FilterAnd(filter1),
			filter1,
		},
		{
			`query.FilterAnd(filter1, filter2)`,
			query.FilterAnd(filter1, filter2),
			query.FilterClause{
				Type:  query.AndOp,
				Inner: []query.FilterClause{filter1, filter2},
			},
		},
		{
			`query.FilterAnd(filter1, query.FilterOr(filter2, filter3))`,
			query.FilterAnd(filter1, query.FilterOr(filter2, filter3)),
			query.FilterClause{
				Type: query.AndOp,
				Inner: []query.FilterClause{
					filter1,
					{
						Type:  query.OrOp,
						Inner: []query.FilterClause{filter2, filter3},
					},
				},
			},
		},
		{
			`query.FilterAnd(query.FilterOr(filter1, filter2), filter3)`,
			query.FilterAnd(query.FilterOr(filter1, filter2), filter3),
			query.FilterClause{
				Type: query.AndOp,
				Inner: []query.FilterClause{
					{
						Type:  query.OrOp,
						Inner: []query.FilterClause{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`query.FilterAnd(query.FilterOr(filter1, filter2), query.FilterOr(filter3, filter4))`,
			query.FilterAnd(query.FilterOr(filter1, filter2), query.FilterOr(filter3, filter4)),
			query.FilterClause{
				Type: query.AndOp,
				Inner: []query.FilterClause{
					{
						Type:  query.OrOp,
						Inner: []query.FilterClause{filter1, filter2},
					},
					{
						Type:  query.OrOp,
						Inner: []query.FilterClause{filter3, filter4},
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
		Operation query.FilterClause
		Result    query.FilterClause
	}{
		{
			`query.FilterOr()`,
			query.FilterOr(),
			query.FilterClause{Type: query.OrOp},
		},
		{
			`query.FilterOr(filter1)`,
			query.FilterOr(filter1),
			filter1,
		},
		{
			`query.FilterOr(filter1, filter2)`,
			query.FilterOr(filter1, filter2),
			query.FilterClause{
				Type:  query.OrOp,
				Inner: []query.FilterClause{filter1, filter2},
			},
		},
		{
			`query.FilterOr(filter1, query.FilterAnd(filter2, filter3))`,
			query.FilterOr(filter1, query.FilterAnd(filter2, filter3)),
			query.FilterClause{
				Type: query.OrOp,
				Inner: []query.FilterClause{
					filter1,
					{
						Type:  query.AndOp,
						Inner: []query.FilterClause{filter2, filter3},
					},
				},
			},
		},
		{
			`query.FilterOr(query.FilterAnd(filter1, filter2), filter3)`,
			query.FilterOr(query.FilterAnd(filter1, filter2), filter3),
			query.FilterClause{
				Type: query.OrOp,
				Inner: []query.FilterClause{
					{
						Type:  query.AndOp,
						Inner: []query.FilterClause{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`query.FilterOr(query.FilterAnd(filter1, filter2), query.FilterAnd(filter3, filter4))`,
			query.FilterOr(query.FilterAnd(filter1, filter2), query.FilterAnd(filter3, filter4)),
			query.FilterClause{
				Type: query.OrOp,
				Inner: []query.FilterClause{
					{
						Type:  query.AndOp,
						Inner: []query.FilterClause{filter1, filter2},
					},
					{
						Type:  query.AndOp,
						Inner: []query.FilterClause{filter3, filter4},
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
		Input    query.FilterOp
		Expected query.FilterOp
	}{
		{
			`Not Eq`,
			query.EqOp,
			query.NeOp,
		},
		{
			`Not Lt`,
			query.LtOp,
			query.GteOp,
		},
		{
			`Not Lte`,
			query.LteOp,
			query.GtOp,
		},
		{
			`Not Gt`,
			query.GtOp,
			query.LteOp,
		},
		{
			`Not Gte`,
			query.GteOp,
			query.LtOp,
		},
		{
			`Not Nil`,
			query.NilOp,
			query.NotNilOp,
		},
		{
			`Not In`,
			query.InOp,
			query.NinOp,
		},
		{
			`Not Like`,
			query.LikeOp,
			query.NotLikeOp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, query.FilterNot(query.FilterClause{Type: tt.Input}).Type)
		})
	}
}

func TestFilterClause_AndEq(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.EqOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, query.FilterClause{}.AndEq("field", "value"))
}

func TestFilterClause_AndNe(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.NeOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, query.FilterClause{}.AndNe("field", "value"))
}

func TestFilterClause_AndLt(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.LtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, query.FilterClause{}.AndLt("field", 10))
}

func TestFilterClause_AndLte(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.LteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, query.FilterClause{}.AndLte("field", 10))
}

func TestFilterClause_AndGt(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.GtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, query.FilterClause{}.AndGt("field", 10))
}

func TestFilterClause_AndGte(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.GteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, query.FilterClause{}.AndGte("field", 10))
}

func TestFilterClause_AndNil(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:  query.NilOp,
				Field: "field",
			},
		},
	}, query.FilterClause{}.AndNil("field"))
}

func TestFilterClause_AndNotNil(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:  query.NotNilOp,
				Field: "field",
			},
		},
	}, query.FilterClause{}.AndNotNil("field"))
}

func TestFilterClause_AndIn(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.InOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, query.FilterClause{}.AndIn("field", "value1", "value2"))
}

func TestFilterClause_AndNin(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.NinOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, query.FilterClause{}.AndNin("field", "value1", "value2"))
}

func TestFilterClause_AndLike(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.LikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, query.FilterClause{}.AndLike("field", "%expr%"))
}

func TestFilterClause_AndNotLike(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.NotLikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, query.FilterClause{}.AndNotLike("field", "%expr%"))
}

func TestFilterClause_AndFragment(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Inner: []query.FilterClause{
			{
				Type:   query.FragmentOp,
				Field:  "expr",
				Values: []interface{}{"value"},
			},
		},
	}, query.FilterClause{}.AndFragment("expr", "value"))
}

func TestFilterClause_OrEq(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.EqOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, query.FilterClause{}.OrEq("field", "value"))
}

func TestFilterClause_OrNe(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.NeOp,
				Field:  "field",
				Values: []interface{}{"value"},
			},
		},
	}, query.FilterClause{}.OrNe("field", "value"))
}

func TestFilterClause_OrLt(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.LtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, query.FilterClause{}.OrLt("field", 10))
}

func TestFilterClause_OrLte(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.LteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, query.FilterClause{}.OrLte("field", 10))
}

func TestFilterClause_OrGt(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.GtOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, query.FilterClause{}.OrGt("field", 10))
}

func TestFilterClause_OrGte(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.GteOp,
				Field:  "field",
				Values: []interface{}{10},
			},
		},
	}, query.FilterClause{}.OrGte("field", 10))
}

func TestFilterClause_OrNil(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:  query.NilOp,
				Field: "field",
			},
		},
	}, query.FilterClause{}.OrNil("field"))
}

func TestFilterClause_OrNotNil(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:  query.NotNilOp,
				Field: "field",
			},
		},
	}, query.FilterClause{}.OrNotNil("field"))
}

func TestFilterClause_OrIn(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.InOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, query.FilterClause{}.OrIn("field", "value1", "value2"))
}

func TestFilterClause_OrNin(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.NinOp,
				Field:  "field",
				Values: []interface{}{"value1", "value2"},
			},
		},
	}, query.FilterClause{}.OrNin("field", "value1", "value2"))
}

func TestFilterClause_OrLike(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.LikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, query.FilterClause{}.OrLike("field", "%expr%"))
}

func TestFilterClause_OrNotLike(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.NotLikeOp,
				Field:  "field",
				Values: []interface{}{"%expr%"},
			},
		},
	}, query.FilterClause{}.OrNotLike("field", "%expr%"))
}

func TestFilterClause_OrFragment(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type: query.OrOp,
		Inner: []query.FilterClause{
			{
				Type:   query.FragmentOp,
				Field:  "expr",
				Values: []interface{}{"value"},
			},
		},
	}, query.FilterClause{}.OrFragment("expr", "value"))
}

func TestFilterEq(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.EqOp,
		Field:  "field",
		Values: []interface{}{"value"},
	}, query.FilterEq("field", "value"))
}

func FilterNe(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.NeOp,
		Field:  "field",
		Values: []interface{}{"value"},
	}, query.FilterNe("field", "value"))
}

func TestFilterLt(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.LtOp,
		Field:  "field",
		Values: []interface{}{10},
	}, query.FilterLt("field", 10))
}

func TestFilterLte(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.LteOp,
		Field:  "field",
		Values: []interface{}{10},
	}, query.FilterLte("field", 10))
}

func TestFilterClauseGt(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.GtOp,
		Field:  "field",
		Values: []interface{}{10},
	}, query.FilterGt("field", 10))
}

func TestFilterGte(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.GteOp,
		Field:  "field",
		Values: []interface{}{10},
	}, query.FilterGte("field", 10))
}

func TestFilterNil(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:  query.NilOp,
		Field: "field",
	}, query.FilterNil("field"))
}

func TestFilterNotNil(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:  query.NotNilOp,
		Field: "field",
	}, query.FilterNotNil("field"))
}

func TestFilterIn(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.InOp,
		Field:  "field",
		Values: []interface{}{"value1", "value2"},
	}, query.FilterIn("field", "value1", "value2"))
}

func TestFilterNin(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.NinOp,
		Field:  "field",
		Values: []interface{}{"value1", "value2"},
	}, query.FilterNin("field", "value1", "value2"))
}

func TestFilterLike(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.LikeOp,
		Field:  "field",
		Values: []interface{}{"%expr%"},
	}, query.FilterLike("field", "%expr%"))
}

func TestFilterNotLike(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.NotLikeOp,
		Field:  "field",
		Values: []interface{}{"%expr%"},
	}, query.FilterNotLike("field", "%expr%"))
}

func TestFilterFragment(t *testing.T) {
	assert.Equal(t, query.FilterClause{
		Type:   query.FragmentOp,
		Field:  "expr",
		Values: []interface{}{"value"},
	}, query.FilterFragment("expr", "value"))
}

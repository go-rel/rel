package grimoire_test

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/stretchr/testify/assert"
)

var result grimoire.Query

func BenchmarkFilterQuery_chain1(b *testing.B) {
	var query grimoire.Query
	for n := 0; n < b.N; n++ {
		query = grimoire.BuildQuery("test", grimoire.Eq("id", 1))
	}
	result = query
}

func BenchmarkFilterQuery_chain2(b *testing.B) {
	var query grimoire.Query
	for n := 0; n < b.N; n++ {
		query = grimoire.BuildQuery("test", grimoire.Eq("id", 1).AndNe("name", "foo"))
	}
	result = query
}

func BenchmarkFilterQuery_chain3(b *testing.B) {
	var query grimoire.Query
	for n := 0; n < b.N; n++ {
		query = grimoire.BuildQuery("test", grimoire.Eq("id", 1).AndNe("name", "foo").AndGt("score", 80))
	}
	result = query
}

func BenchmarkFilterQuery_chain4(b *testing.B) {
	var query grimoire.Query
	for n := 0; n < b.N; n++ {
		query = grimoire.BuildQuery("test", grimoire.Eq("id", 1).AndNe("name", "foo").AndGt("score", 80).AndLt("avg", 10))
	}
	result = query
}

func BenchmarkFilterQuery_slice1(b *testing.B) {
	var query grimoire.Query
	for n := 0; n < b.N; n++ {
		query = grimoire.BuildQuery("test", grimoire.And(grimoire.Eq("id", 1)))
	}
	result = query
}

func BenchmarkFilterQuery_slice2(b *testing.B) {
	var query grimoire.Query
	for n := 0; n < b.N; n++ {
		query = grimoire.BuildQuery("test", grimoire.And(grimoire.Eq("id", 1), grimoire.Ne("name", "foo")))
	}
	result = query
}

func BenchmarkFilterQuery_slice3(b *testing.B) {
	var query grimoire.Query
	for n := 0; n < b.N; n++ {
		query = grimoire.BuildQuery("test", grimoire.And(grimoire.Eq("id", 1), grimoire.Ne("name", "foo"), grimoire.Gt("score", 80)))
	}
	result = query
}

func BenchmarkFilterQuery_slice4(b *testing.B) {
	var query grimoire.Query
	for n := 0; n < b.N; n++ {
		query = grimoire.BuildQuery("test", grimoire.And(grimoire.Eq("id", 1), grimoire.Ne("name", "foo"), grimoire.Gt("score", 80), grimoire.Lt("avg", 10)))
	}
	result = query
}

var filter1 = grimoire.Eq("id", 1)
var filter2 = grimoire.Ne("name", "foo")
var filter3 = grimoire.Gt("score", 80)
var filter4 = grimoire.Lt("avg", 10)

func TestFilterQuery_None(t *testing.T) {
	assert.True(t, grimoire.FilterQuery{}.None())
	assert.True(t, grimoire.And().None())
	assert.True(t, grimoire.Not().None())

	assert.False(t, grimoire.And(filter1).None())
	assert.False(t, grimoire.And(filter1, filter2).None())
	assert.False(t, filter1.None())
}

func TestFilterQuery_And(t *testing.T) {
	tests := []struct {
		Case      string
		Operation grimoire.FilterQuery
		Result    grimoire.FilterQuery
	}{
		{
			`grimoire.FilterQuery{}.And()`,
			grimoire.FilterQuery{}.And(),
			grimoire.And(),
		},
		{
			`grimoire.FilterQuery{}.And(filter1)`,
			grimoire.FilterQuery{}.And(filter1),
			filter1,
		},
		{
			`grimoire.FilterQuery{}.And(filter1).And()`,
			grimoire.FilterQuery{}.And(filter1).And(),
			filter1,
		},
		{
			`grimoire.FilterQuery{}.And(filter1, filter2)`,
			grimoire.FilterQuery{}.And(filter1, filter2),
			grimoire.And(filter1, filter2),
		},
		{
			`grimoire.FilterQuery{}.And(filter1, filter2).And()`,
			grimoire.FilterQuery{}.And(filter1, filter2).And(),
			grimoire.And(filter1, filter2),
		},
		{
			`grimoire.FilterQuery{}.And(filter1, filter2, filter3)`,
			grimoire.FilterQuery{}.And(filter1, filter2, filter3),
			grimoire.And(filter1, filter2, filter3),
		},
		{
			`grimoire.FilterQuery{}.And(filter1, filter2, filter3).And()`,
			grimoire.FilterQuery{}.And(filter1, filter2, filter3).And(),
			grimoire.And(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2)`,
			filter1.And(filter2),
			grimoire.And(filter1, filter2),
		},
		{
			`filter1.And(filter2).And()`,
			filter1.And(filter2).And(),
			grimoire.And(filter1, filter2),
		},
		{
			`filter1.And(filter2).And(filter3)`,
			filter1.And(filter2).And(filter3),
			grimoire.And(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2).And(filter3).And()`,
			filter1.And(filter2).And(filter3).And(),
			grimoire.And(filter1, filter2, filter3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestFilterQuery_Or(t *testing.T) {
	tests := []struct {
		Case      string
		Operation grimoire.FilterQuery
		Result    grimoire.FilterQuery
	}{
		{
			`grimoire.FilterQuery{}.Or()`,
			grimoire.FilterQuery{}.Or(),
			grimoire.Or(),
		},
		{
			`grimoire.FilterQuery{}.Or(filter1)`,
			grimoire.FilterQuery{}.Or(filter1),
			filter1,
		},
		{
			`grimoire.FilterQuery{}.Or(filter1).Or()`,
			grimoire.FilterQuery{}.Or(filter1).Or(),
			filter1,
		},
		{
			`grimoire.FilterQuery{}.Or(filter1, filter2)`,
			grimoire.FilterQuery{}.Or(filter1, filter2),
			grimoire.Or(filter1, filter2),
		},
		{
			`grimoire.FilterQuery{}.Or(filter1, filter2).Or()`,
			grimoire.FilterQuery{}.Or(filter1, filter2).Or(),
			grimoire.Or(filter1, filter2),
		},
		{
			`grimoire.FilterQuery{}.Or(filter1, filter2, filter3)`,
			grimoire.FilterQuery{}.Or(filter1, filter2, filter3),
			grimoire.Or(filter1, filter2, filter3),
		},
		{
			`grimoire.FilterQuery{}.Or(filter1, filter2, filter3).Or()`,
			grimoire.FilterQuery{}.Or(filter1, filter2, filter3).Or(),
			grimoire.Or(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2)`,
			filter1.Or(filter2),
			grimoire.Or(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or()`,
			filter1.Or(filter2).Or(),
			grimoire.Or(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or(filter3)`,
			filter1.Or(filter2).Or(filter3),
			grimoire.Or(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2).Or(filter3).Or()`,
			filter1.Or(filter2).Or(filter3).Or(),
			grimoire.Or(filter1, filter2, filter3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		Case      string
		Operation grimoire.FilterQuery
		Result    grimoire.FilterQuery
	}{
		{
			`grimoire.And()`,
			grimoire.And(),
			grimoire.FilterQuery{Type: grimoire.FilterAndOp},
		},
		{
			`grimoire.And(filter1)`,
			grimoire.And(filter1),
			filter1,
		},
		{
			`grimoire.And(filter1, filter2)`,
			grimoire.And(filter1, filter2),
			grimoire.FilterQuery{
				Type:  grimoire.FilterAndOp,
				Inner: []grimoire.FilterQuery{filter1, filter2},
			},
		},
		{
			`grimoire.And(filter1, grimoire.Or(filter2, filter3))`,
			grimoire.And(filter1, grimoire.Or(filter2, filter3)),
			grimoire.FilterQuery{
				Type: grimoire.FilterAndOp,
				Inner: []grimoire.FilterQuery{
					filter1,
					{
						Type:  grimoire.FilterOrOp,
						Inner: []grimoire.FilterQuery{filter2, filter3},
					},
				},
			},
		},
		{
			`grimoire.And(grimoire.Or(filter1, filter2), filter3)`,
			grimoire.And(grimoire.Or(filter1, filter2), filter3),
			grimoire.FilterQuery{
				Type: grimoire.FilterAndOp,
				Inner: []grimoire.FilterQuery{
					{
						Type:  grimoire.FilterOrOp,
						Inner: []grimoire.FilterQuery{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`grimoire.And(grimoire.Or(filter1, filter2), grimoire.Or(filter3, filter4))`,
			grimoire.And(grimoire.Or(filter1, filter2), grimoire.Or(filter3, filter4)),
			grimoire.FilterQuery{
				Type: grimoire.FilterAndOp,
				Inner: []grimoire.FilterQuery{
					{
						Type:  grimoire.FilterOrOp,
						Inner: []grimoire.FilterQuery{filter1, filter2},
					},
					{
						Type:  grimoire.FilterOrOp,
						Inner: []grimoire.FilterQuery{filter3, filter4},
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

func TestOr(t *testing.T) {
	tests := []struct {
		Case      string
		Operation grimoire.FilterQuery
		Result    grimoire.FilterQuery
	}{
		{
			`grimoire.Or()`,
			grimoire.Or(),
			grimoire.FilterQuery{Type: grimoire.FilterOrOp},
		},
		{
			`grimoire.Or(filter1)`,
			grimoire.Or(filter1),
			filter1,
		},
		{
			`grimoire.Or(filter1, filter2)`,
			grimoire.Or(filter1, filter2),
			grimoire.FilterQuery{
				Type:  grimoire.FilterOrOp,
				Inner: []grimoire.FilterQuery{filter1, filter2},
			},
		},
		{
			`grimoire.Or(filter1, grimoire.And(filter2, filter3))`,
			grimoire.Or(filter1, grimoire.And(filter2, filter3)),
			grimoire.FilterQuery{
				Type: grimoire.FilterOrOp,
				Inner: []grimoire.FilterQuery{
					filter1,
					{
						Type:  grimoire.FilterAndOp,
						Inner: []grimoire.FilterQuery{filter2, filter3},
					},
				},
			},
		},
		{
			`grimoire.Or(grimoire.And(filter1, filter2), filter3)`,
			grimoire.Or(grimoire.And(filter1, filter2), filter3),
			grimoire.FilterQuery{
				Type: grimoire.FilterOrOp,
				Inner: []grimoire.FilterQuery{
					{
						Type:  grimoire.FilterAndOp,
						Inner: []grimoire.FilterQuery{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`grimoire.Or(grimoire.And(filter1, filter2), grimoire.And(filter3, filter4))`,
			grimoire.Or(grimoire.And(filter1, filter2), grimoire.And(filter3, filter4)),
			grimoire.FilterQuery{
				Type: grimoire.FilterOrOp,
				Inner: []grimoire.FilterQuery{
					{
						Type:  grimoire.FilterAndOp,
						Inner: []grimoire.FilterQuery{filter1, filter2},
					},
					{
						Type:  grimoire.FilterAndOp,
						Inner: []grimoire.FilterQuery{filter3, filter4},
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

func TestFilterQuery_Not(t *testing.T) {
	tests := []struct {
		Case     string
		Input    grimoire.FilterOp
		Expected grimoire.FilterOp
	}{
		{
			`Not Eq`,
			grimoire.FilterEqOp,
			grimoire.FilterNeOp,
		},
		{
			`Not Lt`,
			grimoire.FilterLtOp,
			grimoire.FilterGteOp,
		},
		{
			`Not Lte`,
			grimoire.FilterLteOp,
			grimoire.FilterGtOp,
		},
		{
			`Not Gt`,
			grimoire.FilterGtOp,
			grimoire.FilterLteOp,
		},
		{
			`Not Gte`,
			grimoire.FilterGteOp,
			grimoire.FilterLtOp,
		},
		{
			`Not Nil`,
			grimoire.FilterNilOp,
			grimoire.FilterNotNilOp,
		},
		{
			`Not In`,
			grimoire.FilterInOp,
			grimoire.FilterNinOp,
		},
		{
			`Not Like`,
			grimoire.FilterLikeOp,
			grimoire.FilterNotLikeOp,
		},
		{
			`And Op`,
			grimoire.FilterAndOp,
			grimoire.FilterNotOp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, grimoire.Not(grimoire.FilterQuery{Type: tt.Input}).Type)
		})
	}
}

func TestFilterQuery_AndEq(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterEqOp,
				Field: "field",
				Value: "value",
			},
		},
	}, grimoire.FilterQuery{}.AndEq("field", "value"))
}

func TestFilterQuery_AndNe(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNeOp,
				Field: "field",
				Value: "value",
			},
		},
	}, grimoire.FilterQuery{}.AndNe("field", "value"))
}

func TestFilterQuery_AndLt(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterLtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, grimoire.FilterQuery{}.AndLt("field", 10))
}

func TestFilterQuery_AndLte(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterLteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, grimoire.FilterQuery{}.AndLte("field", 10))
}

func TestFilterQuery_AndGt(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterGtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, grimoire.FilterQuery{}.AndGt("field", 10))
}

func TestFilterQuery_AndGte(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterGteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, grimoire.FilterQuery{}.AndGte("field", 10))
}

func TestFilterQuery_AndNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNilOp,
				Field: "field",
			},
		},
	}, grimoire.FilterQuery{}.AndNil("field"))
}

func TestFilterQuery_AndNotNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNotNilOp,
				Field: "field",
			},
		},
	}, grimoire.FilterQuery{}.AndNotNil("field"))
}

func TestFilterQuery_AndIn(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterInOp,
				Field: "field",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}, grimoire.FilterQuery{}.AndIn("field", "value1", "value2"))
}

func TestFilterQuery_AndNin(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNinOp,
				Field: "field",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}, grimoire.FilterQuery{}.AndNin("field", "value1", "value2"))
}

func TestFilterQuery_AndLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, grimoire.FilterQuery{}.AndLike("field", "%expr%"))
}

func TestFilterQuery_AndNotLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNotLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, grimoire.FilterQuery{}.AndNotLike("field", "%expr%"))
}

func TestFilterQuery_AndFragment(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterFragmentOp,
				Field: "expr",
				Value: []interface{}{"value"},
			},
		},
	}, grimoire.FilterQuery{}.AndFragment("expr", "value"))
}

func TestFilterQuery_OrEq(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterEqOp,
				Field: "field",
				Value: "value",
			},
		},
	}, grimoire.FilterQuery{}.OrEq("field", "value"))
}

func TestFilterQuery_OrNe(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNeOp,
				Field: "field",
				Value: "value",
			},
		},
	}, grimoire.FilterQuery{}.OrNe("field", "value"))
}

func TestFilterQuery_OrLt(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterLtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, grimoire.FilterQuery{}.OrLt("field", 10))
}

func TestFilterQuery_OrLte(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterLteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, grimoire.FilterQuery{}.OrLte("field", 10))
}

func TestFilterQuery_OrGt(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterGtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, grimoire.FilterQuery{}.OrGt("field", 10))
}

func TestFilterQuery_OrGte(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterGteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, grimoire.FilterQuery{}.OrGte("field", 10))
}

func TestFilterQuery_OrNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNilOp,
				Field: "field",
			},
		},
	}, grimoire.FilterQuery{}.OrNil("field"))
}

func TestFilterQuery_OrNotNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNotNilOp,
				Field: "field",
			},
		},
	}, grimoire.FilterQuery{}.OrNotNil("field"))
}

func TestFilterQuery_OrIn(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterInOp,
				Field: "field",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}, grimoire.FilterQuery{}.OrIn("field", "value1", "value2"))
}

func TestFilterQuery_OrNin(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNinOp,
				Field: "field",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}, grimoire.FilterQuery{}.OrNin("field", "value1", "value2"))
}

func TestFilterQuery_OrLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, grimoire.FilterQuery{}.OrLike("field", "%expr%"))
}

func TestFilterQuery_OrNotLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterNotLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, grimoire.FilterQuery{}.OrNotLike("field", "%expr%"))
}

func TestFilterQuery_OrFragment(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type: grimoire.FilterOrOp,
		Inner: []grimoire.FilterQuery{
			{
				Type:  grimoire.FilterFragmentOp,
				Field: "expr",
				Value: []interface{}{"value"},
			},
		},
	}, grimoire.FilterQuery{}.OrFragment("expr", "value"))
}

func TestEq(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterEqOp,
		Field: "field",
		Value: "value",
	}, grimoire.Eq("field", "value"))
}

func Ne(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNeOp,
		Field: "field",
		Value: "value",
	}, grimoire.Ne("field", "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterLtOp,
		Field: "field",
		Value: 10,
	}, grimoire.Lt("field", 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterLteOp,
		Field: "field",
		Value: 10,
	}, grimoire.Lte("field", 10))
}

func TestFilterQueryGt(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterGtOp,
		Field: "field",
		Value: 10,
	}, grimoire.Gt("field", 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterGteOp,
		Field: "field",
		Value: 10,
	}, grimoire.Gte("field", 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNilOp,
		Field: "field",
	}, grimoire.Nil("field"))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNotNilOp,
		Field: "field",
	}, grimoire.NotNil("field"))
}

func TestIn(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterInOp,
		Field: "field",
		Value: []interface{}{"value1", "value2"},
	}, grimoire.In("field", "value1", "value2"))
}

func TestNin(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNinOp,
		Field: "field",
		Value: []interface{}{"value1", "value2"},
	}, grimoire.Nin("field", "value1", "value2"))
}

func TestLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterLikeOp,
		Field: "field",
		Value: "%expr%",
	}, grimoire.Like("field", "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterNotLikeOp,
		Field: "field",
		Value: "%expr%",
	}, grimoire.NotLike("field", "%expr%"))
}

func TestFilterFragment(t *testing.T) {
	assert.Equal(t, grimoire.FilterQuery{
		Type:  grimoire.FilterFragmentOp,
		Field: "expr",
		Value: []interface{}{"value"},
	}, grimoire.FilterFragment("expr", "value"))
}

package rel_test

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

var result rel.Query

func BenchmarkFilterQuery_chain1(b *testing.B) {
	var query rel.Query
	for n := 0; n < b.N; n++ {
		query = rel.BuildQuery("test", rel.Eq("id", 1))
	}
	result = query
}

func BenchmarkFilterQuery_chain2(b *testing.B) {
	var query rel.Query
	for n := 0; n < b.N; n++ {
		query = rel.BuildQuery("test", rel.Eq("id", 1).AndNe("name", "foo"))
	}
	result = query
}

func BenchmarkFilterQuery_chain3(b *testing.B) {
	var query rel.Query
	for n := 0; n < b.N; n++ {
		query = rel.BuildQuery("test", rel.Eq("id", 1).AndNe("name", "foo").AndGt("score", 80))
	}
	result = query
}

func BenchmarkFilterQuery_chain4(b *testing.B) {
	var query rel.Query
	for n := 0; n < b.N; n++ {
		query = rel.BuildQuery("test", rel.Eq("id", 1).AndNe("name", "foo").AndGt("score", 80).AndLt("avg", 10))
	}
	result = query
}

func BenchmarkFilterQuery_slice1(b *testing.B) {
	var query rel.Query
	for n := 0; n < b.N; n++ {
		query = rel.BuildQuery("test", rel.And(rel.Eq("id", 1)))
	}
	result = query
}

func BenchmarkFilterQuery_slice2(b *testing.B) {
	var query rel.Query
	for n := 0; n < b.N; n++ {
		query = rel.BuildQuery("test", rel.And(rel.Eq("id", 1), rel.Ne("name", "foo")))
	}
	result = query
}

func BenchmarkFilterQuery_slice3(b *testing.B) {
	var query rel.Query
	for n := 0; n < b.N; n++ {
		query = rel.BuildQuery("test", rel.And(rel.Eq("id", 1), rel.Ne("name", "foo"), rel.Gt("score", 80)))
	}
	result = query
}

func BenchmarkFilterQuery_slice4(b *testing.B) {
	var query rel.Query
	for n := 0; n < b.N; n++ {
		query = rel.BuildQuery("test", rel.And(rel.Eq("id", 1), rel.Ne("name", "foo"), rel.Gt("score", 80), rel.Lt("avg", 10)))
	}
	result = query
}

var filter1 = rel.Eq("id", 1)
var filter2 = rel.Ne("name", "foo")
var filter3 = rel.Gt("score", 80)
var filter4 = rel.Lt("avg", 10)

func TestFilterQuery_None(t *testing.T) {
	assert.True(t, rel.FilterQuery{}.None())
	assert.True(t, rel.And().None())
	assert.True(t, rel.Not().None())

	assert.False(t, rel.And(filter1).None())
	assert.False(t, rel.And(filter1, filter2).None())
	assert.False(t, filter1.None())
}

func TestFilterQuery_And(t *testing.T) {
	tests := []struct {
		Case      string
		Operation rel.FilterQuery
		Result    rel.FilterQuery
	}{
		{
			`rel.FilterQuery{}.And()`,
			rel.FilterQuery{}.And(),
			rel.And(),
		},
		{
			`rel.FilterQuery{}.And(filter1)`,
			rel.FilterQuery{}.And(filter1),
			filter1,
		},
		{
			`rel.FilterQuery{}.And(filter1).And()`,
			rel.FilterQuery{}.And(filter1).And(),
			filter1,
		},
		{
			`rel.FilterQuery{}.And(filter1, filter2)`,
			rel.FilterQuery{}.And(filter1, filter2),
			rel.And(filter1, filter2),
		},
		{
			`rel.FilterQuery{}.And(filter1, filter2).And()`,
			rel.FilterQuery{}.And(filter1, filter2).And(),
			rel.And(filter1, filter2),
		},
		{
			`rel.FilterQuery{}.And(filter1, filter2, filter3)`,
			rel.FilterQuery{}.And(filter1, filter2, filter3),
			rel.And(filter1, filter2, filter3),
		},
		{
			`rel.FilterQuery{}.And(filter1, filter2, filter3).And()`,
			rel.FilterQuery{}.And(filter1, filter2, filter3).And(),
			rel.And(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2)`,
			filter1.And(filter2),
			rel.And(filter1, filter2),
		},
		{
			`filter1.And(filter2).And()`,
			filter1.And(filter2).And(),
			rel.And(filter1, filter2),
		},
		{
			`filter1.And(filter2).And(filter3)`,
			filter1.And(filter2).And(filter3),
			rel.And(filter1, filter2, filter3),
		},
		{
			`filter1.And(filter2).And(filter3).And()`,
			filter1.And(filter2).And(filter3).And(),
			rel.And(filter1, filter2, filter3),
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
		Operation rel.FilterQuery
		Result    rel.FilterQuery
	}{
		{
			`rel.FilterQuery{}.Or()`,
			rel.FilterQuery{}.Or(),
			rel.Or(),
		},
		{
			`rel.FilterQuery{}.Or(filter1)`,
			rel.FilterQuery{}.Or(filter1),
			filter1,
		},
		{
			`rel.FilterQuery{}.Or(filter1).Or()`,
			rel.FilterQuery{}.Or(filter1).Or(),
			filter1,
		},
		{
			`rel.FilterQuery{}.Or(filter1, filter2)`,
			rel.FilterQuery{}.Or(filter1, filter2),
			rel.Or(filter1, filter2),
		},
		{
			`rel.FilterQuery{}.Or(filter1, filter2).Or()`,
			rel.FilterQuery{}.Or(filter1, filter2).Or(),
			rel.Or(filter1, filter2),
		},
		{
			`rel.FilterQuery{}.Or(filter1, filter2, filter3)`,
			rel.FilterQuery{}.Or(filter1, filter2, filter3),
			rel.Or(filter1, filter2, filter3),
		},
		{
			`rel.FilterQuery{}.Or(filter1, filter2, filter3).Or()`,
			rel.FilterQuery{}.Or(filter1, filter2, filter3).Or(),
			rel.Or(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2)`,
			filter1.Or(filter2),
			rel.Or(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or()`,
			filter1.Or(filter2).Or(),
			rel.Or(filter1, filter2),
		},
		{
			`filter1.Or(filter2).Or(filter3)`,
			filter1.Or(filter2).Or(filter3),
			rel.Or(filter1, filter2, filter3),
		},
		{
			`filter1.Or(filter2).Or(filter3).Or()`,
			filter1.Or(filter2).Or(filter3).Or(),
			rel.Or(filter1, filter2, filter3),
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
		Operation rel.FilterQuery
		Result    rel.FilterQuery
	}{
		{
			`rel.And()`,
			rel.And(),
			rel.FilterQuery{Type: rel.FilterAndOp},
		},
		{
			`rel.And(filter1)`,
			rel.And(filter1),
			filter1,
		},
		{
			`rel.And(filter1, filter2)`,
			rel.And(filter1, filter2),
			rel.FilterQuery{
				Type:  rel.FilterAndOp,
				Inner: []rel.FilterQuery{filter1, filter2},
			},
		},
		{
			`rel.And(filter1, rel.Or(filter2, filter3))`,
			rel.And(filter1, rel.Or(filter2, filter3)),
			rel.FilterQuery{
				Type: rel.FilterAndOp,
				Inner: []rel.FilterQuery{
					filter1,
					{
						Type:  rel.FilterOrOp,
						Inner: []rel.FilterQuery{filter2, filter3},
					},
				},
			},
		},
		{
			`rel.And(rel.Or(filter1, filter2), filter3)`,
			rel.And(rel.Or(filter1, filter2), filter3),
			rel.FilterQuery{
				Type: rel.FilterAndOp,
				Inner: []rel.FilterQuery{
					{
						Type:  rel.FilterOrOp,
						Inner: []rel.FilterQuery{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`rel.And(rel.Or(filter1, filter2), rel.Or(filter3, filter4))`,
			rel.And(rel.Or(filter1, filter2), rel.Or(filter3, filter4)),
			rel.FilterQuery{
				Type: rel.FilterAndOp,
				Inner: []rel.FilterQuery{
					{
						Type:  rel.FilterOrOp,
						Inner: []rel.FilterQuery{filter1, filter2},
					},
					{
						Type:  rel.FilterOrOp,
						Inner: []rel.FilterQuery{filter3, filter4},
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
		Operation rel.FilterQuery
		Result    rel.FilterQuery
	}{
		{
			`rel.Or()`,
			rel.Or(),
			rel.FilterQuery{Type: rel.FilterOrOp},
		},
		{
			`rel.Or(filter1)`,
			rel.Or(filter1),
			filter1,
		},
		{
			`rel.Or(filter1, filter2)`,
			rel.Or(filter1, filter2),
			rel.FilterQuery{
				Type:  rel.FilterOrOp,
				Inner: []rel.FilterQuery{filter1, filter2},
			},
		},
		{
			`rel.Or(filter1, rel.And(filter2, filter3))`,
			rel.Or(filter1, rel.And(filter2, filter3)),
			rel.FilterQuery{
				Type: rel.FilterOrOp,
				Inner: []rel.FilterQuery{
					filter1,
					{
						Type:  rel.FilterAndOp,
						Inner: []rel.FilterQuery{filter2, filter3},
					},
				},
			},
		},
		{
			`rel.Or(rel.And(filter1, filter2), filter3)`,
			rel.Or(rel.And(filter1, filter2), filter3),
			rel.FilterQuery{
				Type: rel.FilterOrOp,
				Inner: []rel.FilterQuery{
					{
						Type:  rel.FilterAndOp,
						Inner: []rel.FilterQuery{filter1, filter2},
					},
					filter3,
				},
			},
		},
		{
			`rel.Or(rel.And(filter1, filter2), rel.And(filter3, filter4))`,
			rel.Or(rel.And(filter1, filter2), rel.And(filter3, filter4)),
			rel.FilterQuery{
				Type: rel.FilterOrOp,
				Inner: []rel.FilterQuery{
					{
						Type:  rel.FilterAndOp,
						Inner: []rel.FilterQuery{filter1, filter2},
					},
					{
						Type:  rel.FilterAndOp,
						Inner: []rel.FilterQuery{filter3, filter4},
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
		Input    rel.FilterOp
		Expected rel.FilterOp
	}{
		{
			`Not Eq`,
			rel.FilterEqOp,
			rel.FilterNeOp,
		},
		{
			`Not Lt`,
			rel.FilterLtOp,
			rel.FilterGteOp,
		},
		{
			`Not Lte`,
			rel.FilterLteOp,
			rel.FilterGtOp,
		},
		{
			`Not Gt`,
			rel.FilterGtOp,
			rel.FilterLteOp,
		},
		{
			`Not Gte`,
			rel.FilterGteOp,
			rel.FilterLtOp,
		},
		{
			`Not Nil`,
			rel.FilterNilOp,
			rel.FilterNotNilOp,
		},
		{
			`Not In`,
			rel.FilterInOp,
			rel.FilterNinOp,
		},
		{
			`Not Like`,
			rel.FilterLikeOp,
			rel.FilterNotLikeOp,
		},
		{
			`And Op`,
			rel.FilterAndOp,
			rel.FilterNotOp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, rel.Not(rel.FilterQuery{Type: tt.Input}).Type)
		})
	}
}

func TestFilterQuery_AndEq(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterEqOp,
				Field: "field",
				Value: "value",
			},
		},
	}, rel.FilterQuery{}.AndEq("field", "value"))
}

func TestFilterQuery_AndNe(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNeOp,
				Field: "field",
				Value: "value",
			},
		},
	}, rel.FilterQuery{}.AndNe("field", "value"))
}

func TestFilterQuery_AndLt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterLtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, rel.FilterQuery{}.AndLt("field", 10))
}

func TestFilterQuery_AndLte(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterLteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, rel.FilterQuery{}.AndLte("field", 10))
}

func TestFilterQuery_AndGt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterGtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, rel.FilterQuery{}.AndGt("field", 10))
}

func TestFilterQuery_AndGte(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterGteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, rel.FilterQuery{}.AndGte("field", 10))
}

func TestFilterQuery_AndNil(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNilOp,
				Field: "field",
			},
		},
	}, rel.FilterQuery{}.AndNil("field"))
}

func TestFilterQuery_AndNotNil(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNotNilOp,
				Field: "field",
			},
		},
	}, rel.FilterQuery{}.AndNotNil("field"))
}

func TestFilterQuery_AndIn(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterInOp,
				Field: "field",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}, rel.FilterQuery{}.AndIn("field", "value1", "value2"))
}

func TestFilterQuery_AndNin(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNinOp,
				Field: "field",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}, rel.FilterQuery{}.AndNin("field", "value1", "value2"))
}

func TestFilterQuery_AndLike(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, rel.FilterQuery{}.AndLike("field", "%expr%"))
}

func TestFilterQuery_AndNotLike(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNotLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, rel.FilterQuery{}.AndNotLike("field", "%expr%"))
}

func TestFilterQuery_AndFragment(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterFragmentOp,
				Field: "expr",
				Value: []interface{}{"value"},
			},
		},
	}, rel.FilterQuery{}.AndFragment("expr", "value"))
}

func TestFilterQuery_OrEq(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterEqOp,
				Field: "field",
				Value: "value",
			},
		},
	}, rel.FilterQuery{}.OrEq("field", "value"))
}

func TestFilterQuery_OrNe(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNeOp,
				Field: "field",
				Value: "value",
			},
		},
	}, rel.FilterQuery{}.OrNe("field", "value"))
}

func TestFilterQuery_OrLt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterLtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, rel.FilterQuery{}.OrLt("field", 10))
}

func TestFilterQuery_OrLte(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterLteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, rel.FilterQuery{}.OrLte("field", 10))
}

func TestFilterQuery_OrGt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterGtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, rel.FilterQuery{}.OrGt("field", 10))
}

func TestFilterQuery_OrGte(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterGteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, rel.FilterQuery{}.OrGte("field", 10))
}

func TestFilterQuery_OrNil(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNilOp,
				Field: "field",
			},
		},
	}, rel.FilterQuery{}.OrNil("field"))
}

func TestFilterQuery_OrNotNil(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNotNilOp,
				Field: "field",
			},
		},
	}, rel.FilterQuery{}.OrNotNil("field"))
}

func TestFilterQuery_OrIn(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterInOp,
				Field: "field",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}, rel.FilterQuery{}.OrIn("field", "value1", "value2"))
}

func TestFilterQuery_OrNin(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNinOp,
				Field: "field",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}, rel.FilterQuery{}.OrNin("field", "value1", "value2"))
}

func TestFilterQuery_OrLike(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, rel.FilterQuery{}.OrLike("field", "%expr%"))
}

func TestFilterQuery_OrNotLike(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterNotLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, rel.FilterQuery{}.OrNotLike("field", "%expr%"))
}

func TestFilterQuery_OrFragment(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type: rel.FilterOrOp,
		Inner: []rel.FilterQuery{
			{
				Type:  rel.FilterFragmentOp,
				Field: "expr",
				Value: []interface{}{"value"},
			},
		},
	}, rel.FilterQuery{}.OrFragment("expr", "value"))
}

func TestEq(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterEqOp,
		Field: "field",
		Value: "value",
	}, rel.Eq("field", "value"))
}

func Ne(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNeOp,
		Field: "field",
		Value: "value",
	}, rel.Ne("field", "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterLtOp,
		Field: "field",
		Value: 10,
	}, rel.Lt("field", 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterLteOp,
		Field: "field",
		Value: 10,
	}, rel.Lte("field", 10))
}

func TestFilterQueryGt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterGtOp,
		Field: "field",
		Value: 10,
	}, rel.Gt("field", 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterGteOp,
		Field: "field",
		Value: 10,
	}, rel.Gte("field", 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNilOp,
		Field: "field",
	}, rel.Nil("field"))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNotNilOp,
		Field: "field",
	}, rel.NotNil("field"))
}

func TestIn(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterInOp,
		Field: "field",
		Value: []interface{}{"value1", "value2"},
	}, rel.In("field", "value1", "value2"))
}

func TestInInt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterInOp,
		Field: "field",
		Value: []interface{}{1, 2},
	}, rel.InInt("field", []int{1, 2}))
}

func TestInUint(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterInOp,
		Field: "field",
		Value: []interface{}{uint(1), uint(2)},
	}, rel.InUint("field", []uint{1, 2}))
}

func TestInString(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterInOp,
		Field: "field",
		Value: []interface{}{"1", "2"},
	}, rel.InString("field", []string{"1", "2"}))
}

func TestNin(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNinOp,
		Field: "field",
		Value: []interface{}{"value1", "value2"},
	}, rel.Nin("field", "value1", "value2"))
}

func TestNinInt(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNinOp,
		Field: "field",
		Value: []interface{}{1, 2},
	}, rel.NinInt("field", []int{1, 2}))
}

func TestNinUint(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNinOp,
		Field: "field",
		Value: []interface{}{uint(1), uint(2)},
	}, rel.NinUint("field", []uint{1, 2}))
}

func TestNinString(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNinOp,
		Field: "field",
		Value: []interface{}{"1", "2"},
	}, rel.NinString("field", []string{"1", "2"}))
}
func TestLike(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterLikeOp,
		Field: "field",
		Value: "%expr%",
	}, rel.Like("field", "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterNotLikeOp,
		Field: "field",
		Value: "%expr%",
	}, rel.NotLike("field", "%expr%"))
}

func TestFilterFragment(t *testing.T) {
	assert.Equal(t, rel.FilterQuery{
		Type:  rel.FilterFragmentOp,
		Field: "expr",
		Value: []interface{}{"value"},
	}, rel.FilterFragment("expr", "value"))
}

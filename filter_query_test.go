package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var result Query

func BenchmarkFilterQuery_chain1(b *testing.B) {
	var query Query
	for n := 0; n < b.N; n++ {
		query = Build("test", Eq("id", 1))
	}
	result = query
}

func BenchmarkFilterQuery_chain2(b *testing.B) {
	var query Query
	for n := 0; n < b.N; n++ {
		query = Build("test", Eq("id", 1).AndNe("name", "foo"))
	}
	result = query
}

func BenchmarkFilterQuery_chain3(b *testing.B) {
	var query Query
	for n := 0; n < b.N; n++ {
		query = Build("test", Eq("id", 1).AndNe("name", "foo").AndGt("score", 80))
	}
	result = query
}

func BenchmarkFilterQuery_chain4(b *testing.B) {
	var query Query
	for n := 0; n < b.N; n++ {
		query = Build("test", Eq("id", 1).AndNe("name", "foo").AndGt("score", 80).AndLt("avg", 10))
	}
	result = query
}

func BenchmarkFilterQuery_slice1(b *testing.B) {
	var query Query
	for n := 0; n < b.N; n++ {
		query = Build("test", And(Eq("id", 1)))
	}
	result = query
}

func BenchmarkFilterQuery_slice2(b *testing.B) {
	var query Query
	for n := 0; n < b.N; n++ {
		query = Build("test", And(Eq("id", 1), Ne("name", "foo")))
	}
	result = query
}

func BenchmarkFilterQuery_slice3(b *testing.B) {
	var query Query
	for n := 0; n < b.N; n++ {
		query = Build("test", And(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)))
	}
	result = query
}

func BenchmarkFilterQuery_slice4(b *testing.B) {
	var query Query
	for n := 0; n < b.N; n++ {
		query = Build("test", And(Eq("id", 1), Ne("name", "foo"), Gt("score", 80), Lt("avg", 10)))
	}
	result = query
}

func TestFilterQuery_None(t *testing.T) {
	assert.True(t, FilterQuery{}.None())
	assert.True(t, And().None())
	assert.True(t, Not().None())

	assert.False(t, And(Eq("id", 1)).None())
	assert.False(t, And(Eq("id", 1), Ne("name", "foo")).None())
	assert.False(t, Eq("id", 1).None())
}

func TestFilterQuery_And(t *testing.T) {
	tests := []struct {
		name   string
		filter FilterQuery
		result FilterQuery
	}{
		{
			name:   "",
			filter: FilterQuery{}.And(),
			result: And(),
		},
		{
			name:   "where.Eq(\"id\", 1)",
			filter: FilterQuery{}.And(Eq("id", 1)),
			result: Eq("id", 1),
		},
		{
			name:   "where.Eq(\"id\", 1)",
			filter: FilterQuery{}.And(Eq("id", 1)).And(),
			result: Eq("id", 1),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: FilterQuery{}.And(Eq("id", 1), Ne("name", "foo")),
			result: And(Eq("id", 1), Ne("name", "foo")),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: FilterQuery{}.And(Eq("id", 1), Ne("name", "foo")).And(),
			result: And(Eq("id", 1), Ne("name", "foo")),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80))",
			filter: FilterQuery{}.And(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
			result: And(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80))",
			filter: FilterQuery{}.And(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)).And(),
			result: And(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: Eq("id", 1).And(Ne("name", "foo")),
			result: And(Eq("id", 1), Ne("name", "foo")),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: Eq("id", 1).And(Ne("name", "foo")).And(),
			result: And(Eq("id", 1), Ne("name", "foo")),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80))",
			filter: Eq("id", 1).And(Ne("name", "foo")).And(Gt("score", 80)),
			result: And(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80))",
			filter: Eq("id", 1).And(Ne("name", "foo")).And(Gt("score", 80)).And(),
			result: And(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Fragment(\"name = ?\", \"name\"))",
			filter: Eq("id", 1).AndFragment("name = ?", "name"),
			result: And(Eq("id", 1), FilterFragment("name = ?", "name")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, tt.filter)
			assert.Equal(t, tt.name, tt.filter.String())
		})
	}
}

func TestFilterQuery_Or(t *testing.T) {
	tests := []struct {
		name   string
		filter FilterQuery
		result FilterQuery
	}{
		{
			name:   "",
			filter: FilterQuery{}.Or(),
			result: Or(),
		},
		{
			name:   "where.Eq(\"id\", 1)",
			filter: FilterQuery{}.Or(Eq("id", 1)),
			result: Eq("id", 1),
		},
		{
			name:   "where.Eq(\"id\", 1)",
			filter: FilterQuery{}.Or(Eq("id", 1)).Or(),
			result: Eq("id", 1),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: FilterQuery{}.Or(Eq("id", 1), Ne("name", "foo")),
			result: Or(Eq("id", 1), Ne("name", "foo")),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: FilterQuery{}.Or(Eq("id", 1), Ne("name", "foo")).Or(),
			result: Or(Eq("id", 1), Ne("name", "foo")),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80))",
			filter: FilterQuery{}.Or(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
			result: Or(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80))",
			filter: FilterQuery{}.Or(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)).Or(),
			result: Or(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: Eq("id", 1).Or(Ne("name", "foo")),
			result: Or(Eq("id", 1), Ne("name", "foo")),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: Eq("id", 1).Or(Ne("name", "foo")).Or(),
			result: Or(Eq("id", 1), Ne("name", "foo")),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80))",
			filter: Eq("id", 1).Or(Ne("name", "foo")).Or(Gt("score", 80)),
			result: Or(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80))",
			filter: Eq("id", 1).Or(Ne("name", "foo")).Or(Gt("score", 80)).Or(),
			result: Or(Eq("id", 1), Ne("name", "foo"), Gt("score", 80)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, tt.filter)
			assert.Equal(t, tt.name, tt.filter.String())
		})
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		name   string
		filter FilterQuery
		result FilterQuery
	}{
		{
			name:   "",
			filter: And(),
			result: FilterQuery{Type: FilterAndOp},
		},
		{
			name:   "where.Eq(\"id\", 1)",
			filter: And(Eq("id", 1)),
			result: Eq("id", 1),
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: And(Eq("id", 1), Ne("name", "foo")),
			result: FilterQuery{
				Type:  FilterAndOp,
				Inner: []FilterQuery{Eq("id", 1), Ne("name", "foo")},
			},
		},
		{
			name:   "where.And(where.Eq(\"id\", 1), where.Or(where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80)))",
			filter: And(Eq("id", 1), Or(Ne("name", "foo"), Gt("score", 80))),
			result: FilterQuery{
				Type: FilterAndOp,
				Inner: []FilterQuery{
					Eq("id", 1),
					{
						Type:  FilterOrOp,
						Inner: []FilterQuery{Ne("name", "foo"), Gt("score", 80)},
					},
				},
			},
		},
		{
			name:   "where.And(where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\")), where.Gt(\"score\", 80))",
			filter: And(Or(Eq("id", 1), Ne("name", "foo")), Gt("score", 80)),
			result: FilterQuery{
				Type: FilterAndOp,
				Inner: []FilterQuery{
					{
						Type:  FilterOrOp,
						Inner: []FilterQuery{Eq("id", 1), Ne("name", "foo")},
					},
					Gt("score", 80),
				},
			},
		},
		{
			name:   "where.And(where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\")), where.Or(where.Gt(\"score\", 80), where.Lt(\"avg\", 10)))",
			filter: And(Or(Eq("id", 1), Ne("name", "foo")), Or(Gt("score", 80), Lt("avg", 10))),
			result: FilterQuery{
				Type: FilterAndOp,
				Inner: []FilterQuery{
					{
						Type:  FilterOrOp,
						Inner: []FilterQuery{Eq("id", 1), Ne("name", "foo")},
					},
					{
						Type:  FilterOrOp,
						Inner: []FilterQuery{Gt("score", 80), Lt("avg", 10)},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, tt.filter)
			assert.Equal(t, tt.name, tt.filter.String())
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		name   string
		filter FilterQuery
		result FilterQuery
	}{
		{
			name:   "",
			filter: Or(),
			result: FilterQuery{Type: FilterOrOp},
		},
		{
			name:   "where.Eq(\"id\", 1)",
			filter: Or(Eq("id", 1)),
			result: Eq("id", 1),
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\"))",
			filter: Or(Eq("id", 1), Ne("name", "foo")),
			result: FilterQuery{
				Type:  FilterOrOp,
				Inner: []FilterQuery{Eq("id", 1), Ne("name", "foo")},
			},
		},
		{
			name:   "where.Or(where.Eq(\"id\", 1), where.And(where.Ne(\"name\", \"foo\"), where.Gt(\"score\", 80)))",
			filter: Or(Eq("id", 1), And(Ne("name", "foo"), Gt("score", 80))),
			result: FilterQuery{
				Type: FilterOrOp,
				Inner: []FilterQuery{
					Eq("id", 1),
					{
						Type:  FilterAndOp,
						Inner: []FilterQuery{Ne("name", "foo"), Gt("score", 80)},
					},
				},
			},
		},
		{
			name:   "where.Or(where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\")), where.Gt(\"score\", 80))",
			filter: Or(And(Eq("id", 1), Ne("name", "foo")), Gt("score", 80)),
			result: FilterQuery{
				Type: FilterOrOp,
				Inner: []FilterQuery{
					{
						Type:  FilterAndOp,
						Inner: []FilterQuery{Eq("id", 1), Ne("name", "foo")},
					},
					Gt("score", 80),
				},
			},
		},
		{
			name:   "where.Or(where.And(where.Eq(\"id\", 1), where.Ne(\"name\", \"foo\")), where.And(where.Gt(\"score\", 80), where.Lt(\"avg\", 10)))",
			filter: Or(And(Eq("id", 1), Ne("name", "foo")), And(Gt("score", 80), Lt("avg", 10))),
			result: FilterQuery{
				Type: FilterOrOp,
				Inner: []FilterQuery{
					{
						Type:  FilterAndOp,
						Inner: []FilterQuery{Eq("id", 1), Ne("name", "foo")},
					},
					{
						Type:  FilterAndOp,
						Inner: []FilterQuery{Gt("score", 80), Lt("avg", 10)},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.result, tt.filter)
			assert.Equal(t, tt.name, tt.filter.String())
		})
	}
}

func TestFilterQuery_Not(t *testing.T) {
	tests := []struct {
		Case     string
		Input    FilterOp
		Expected FilterOp
	}{
		{
			`Not Eq`,
			FilterEqOp,
			FilterNeOp,
		},
		{
			`Not Lt`,
			FilterLtOp,
			FilterGteOp,
		},
		{
			`Not Lte`,
			FilterLteOp,
			FilterGtOp,
		},
		{
			`Not Gt`,
			FilterGtOp,
			FilterLteOp,
		},
		{
			`Not Gte`,
			FilterGteOp,
			FilterLtOp,
		},
		{
			`Not Nil`,
			FilterNilOp,
			FilterNotNilOp,
		},
		{
			`Not In`,
			FilterInOp,
			FilterNinOp,
		},
		{
			`Not Like`,
			FilterLikeOp,
			FilterNotLikeOp,
		},
		{
			`And Op`,
			FilterAndOp,
			FilterNotOp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, Not(FilterQuery{Type: tt.Input}).Type)
		})
	}
}

func TestFilterQuery_AndEq(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterEqOp,
				Field: "field",
				Value: "value",
			},
		},
	}, FilterQuery{}.AndEq("field", "value"))
}

func TestFilterQuery_AndNe(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterNeOp,
				Field: "field",
				Value: "value",
			},
		},
	}, FilterQuery{}.AndNe("field", "value"))
}

func TestFilterQuery_AndLt(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterLtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, FilterQuery{}.AndLt("field", 10))
}

func TestFilterQuery_AndLte(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterLteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, FilterQuery{}.AndLte("field", 10))
}

func TestFilterQuery_AndGt(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterGtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, FilterQuery{}.AndGt("field", 10))
}

func TestFilterQuery_AndGte(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterGteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, FilterQuery{}.AndGte("field", 10))
}

func TestFilterQuery_AndNil(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterNilOp,
				Field: "field",
			},
		},
	}, FilterQuery{}.AndNil("field"))
}

func TestFilterQuery_AndNotNil(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterNotNilOp,
				Field: "field",
			},
		},
	}, FilterQuery{}.AndNotNil("field"))
}

func TestFilterQuery_AndIn(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterInOp,
				Field: "field",
				Value: []any{"value1", "value2"},
			},
		},
	}, FilterQuery{}.AndIn("field", "value1", "value2"))
}

func TestFilterQuery_AndNin(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterNinOp,
				Field: "field",
				Value: []any{"value1", "value2"},
			},
		},
	}, FilterQuery{}.AndNin("field", "value1", "value2"))
}

func TestFilterQuery_AndLike(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, FilterQuery{}.AndLike("field", "%expr%"))
}

func TestFilterQuery_AndNotLike(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterNotLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, FilterQuery{}.AndNotLike("field", "%expr%"))
}

func TestFilterQuery_AndFragment(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Inner: []FilterQuery{
			{
				Type:  FilterFragmentOp,
				Field: "expr",
				Value: []any{"value"},
			},
		},
	}, FilterQuery{}.AndFragment("expr", "value"))
}

func TestFilterQuery_OrEq(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterEqOp,
				Field: "field",
				Value: "value",
			},
		},
	}, FilterQuery{}.OrEq("field", "value"))
}

func TestFilterQuery_OrNe(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterNeOp,
				Field: "field",
				Value: "value",
			},
		},
	}, FilterQuery{}.OrNe("field", "value"))
}

func TestFilterQuery_OrLt(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterLtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, FilterQuery{}.OrLt("field", 10))
}

func TestFilterQuery_OrLte(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterLteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, FilterQuery{}.OrLte("field", 10))
}

func TestFilterQuery_OrGt(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterGtOp,
				Field: "field",
				Value: 10,
			},
		},
	}, FilterQuery{}.OrGt("field", 10))
}

func TestFilterQuery_OrGte(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterGteOp,
				Field: "field",
				Value: 10,
			},
		},
	}, FilterQuery{}.OrGte("field", 10))
}

func TestFilterQuery_OrNil(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterNilOp,
				Field: "field",
			},
		},
	}, FilterQuery{}.OrNil("field"))
}

func TestFilterQuery_OrNotNil(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterNotNilOp,
				Field: "field",
			},
		},
	}, FilterQuery{}.OrNotNil("field"))
}

func TestFilterQuery_OrIn(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterInOp,
				Field: "field",
				Value: []any{"value1", "value2"},
			},
		},
	}, FilterQuery{}.OrIn("field", "value1", "value2"))
}

func TestFilterQuery_OrNin(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterNinOp,
				Field: "field",
				Value: []any{"value1", "value2"},
			},
		},
	}, FilterQuery{}.OrNin("field", "value1", "value2"))
}

func TestFilterQuery_OrLike(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, FilterQuery{}.OrLike("field", "%expr%"))
}

func TestFilterQuery_OrNotLike(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterNotLikeOp,
				Field: "field",
				Value: "%expr%",
			},
		},
	}, FilterQuery{}.OrNotLike("field", "%expr%"))
}

func TestFilterQuery_OrFragment(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type: FilterOrOp,
		Inner: []FilterQuery{
			{
				Type:  FilterFragmentOp,
				Field: "expr",
				Value: []any{"value"},
			},
		},
	}, FilterQuery{}.OrFragment("expr", "value"))
}

func TestEq(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterEqOp,
		Field: "field",
		Value: "value",
	}, Eq("field", "value"))
}

func TestNe(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterNeOp,
		Field: "field",
		Value: "value",
	}, Ne("field", "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterLtOp,
		Field: "field",
		Value: 10,
	}, Lt("field", 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterLteOp,
		Field: "field",
		Value: 10,
	}, Lte("field", 10))
}

func TestFilterQueryGt(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterGtOp,
		Field: "field",
		Value: 10,
	}, Gt("field", 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterGteOp,
		Field: "field",
		Value: 10,
	}, Gte("field", 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterNilOp,
		Field: "field",
	}, Nil("field"))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterNotNilOp,
		Field: "field",
	}, NotNil("field"))
}

func TestIn(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterInOp,
		Field: "field",
		Value: []any{"value1", "value2"},
	}, In("field", "value1", "value2"))
}

func TestInInt(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterInOp,
		Field: "field",
		Value: []any{1, 2},
	}, InInt("field", []int{1, 2}))
}

func TestInUint(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterInOp,
		Field: "field",
		Value: []any{uint(1), uint(2)},
	}, InUint("field", []uint{1, 2}))
}

func TestInString(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterInOp,
		Field: "field",
		Value: []any{"1", "2"},
	}, InString("field", []string{"1", "2"}))
}

func TestNin(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterNinOp,
		Field: "field",
		Value: []any{"value1", "value2"},
	}, Nin("field", "value1", "value2"))
}

func TestNinInt(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterNinOp,
		Field: "field",
		Value: []any{1, 2},
	}, NinInt("field", []int{1, 2}))
}

func TestNinUint(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterNinOp,
		Field: "field",
		Value: []any{uint(1), uint(2)},
	}, NinUint("field", []uint{1, 2}))
}

func TestNinString(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterNinOp,
		Field: "field",
		Value: []any{"1", "2"},
	}, NinString("field", []string{"1", "2"}))
}

func TestLike(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterLikeOp,
		Field: "field",
		Value: "%expr%",
	}, Like("field", "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterNotLikeOp,
		Field: "field",
		Value: "%expr%",
	}, NotLike("field", "%expr%"))
}

func TestFilterFragment(t *testing.T) {
	assert.Equal(t, FilterQuery{
		Type:  FilterFragmentOp,
		Field: "expr",
		Value: []any{"value"},
	}, FilterFragment("expr", "value"))
}

func TestFilterDocument(t *testing.T) {
	var (
		user = User{ID: 1}
		doc  = NewDocument(&user)
	)

	assert.Equal(t, Eq("id", 1), filterDocument(doc))
}

func TestFilterDocument_compositePrimaryKey(t *testing.T) {
	var (
		userRole = UserRole{UserID: 1, RoleID: 2}
		doc      = NewDocument(&userRole)
	)

	assert.Equal(t, Eq("user_id", 1).AndEq("role_id", 2), filterDocument(doc))
}

func TestFilterCollection(t *testing.T) {
	var (
		users = []User{
			{ID: 1},
			{ID: 2},
		}
		col = NewCollection(&users)
	)

	assert.Equal(t, In("id", 1, 2), filterCollection(col))
}

func TestFilterCollection_compositePrimaryKey(t *testing.T) {
	var (
		userRoles = []UserRole{
			{UserID: 1, RoleID: 2},
			{UserID: 3, RoleID: 4},
		}
		col = NewCollection(&userRoles)
	)

	assert.Equal(t, Or(Eq("user_id", 1).AndEq("role_id", 2), Eq("user_id", 3).AndEq("role_id", 4)), filterCollection(col))
}

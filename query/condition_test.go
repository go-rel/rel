package query_test

import (
	. "github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

var cond1 = Eq("id", 1)
var cond2 = Ne("name", "foo")
var cond3 = Gt("score", 80)
var cond4 = Lt("avg", 10)

func TestConditionNone(t *testing.T) {
	assert.True(t, Condition{}.None())
	assert.True(t, And().None())
	assert.True(t, Xor().None())
	assert.True(t, Not().None())

	assert.False(t, And(cond1).None())
	assert.False(t, And(cond1, cond2).None())
	assert.False(t, cond1.None())
}

func TestConditionAnd(t *testing.T) {
	tests := []struct {
		Case      string
		Operation Condition
		Result    Condition
	}{
		{
			`Condition{}.And()`,
			Condition{}.And(),
			And(),
		},
		{
			`Condition{}.And(cond1)`,
			Condition{}.And(cond1),
			cond1,
		},
		{
			`Condition{}.And(cond1).And()`,
			Condition{}.And(cond1).And(),
			cond1,
		},
		{
			`Condition{}.And(cond1, cond2)`,
			Condition{}.And(cond1, cond2),
			And(cond1, cond2),
		},
		{
			`Condition{}.And(cond1, cond2).And()`,
			Condition{}.And(cond1, cond2).And(),
			And(cond1, cond2),
		},
		{
			`Condition{}.And(cond1, cond2, cond3)`,
			Condition{}.And(cond1, cond2, cond3),
			And(cond1, cond2, cond3),
		},
		{
			`Condition{}.And(cond1, cond2, cond3).And()`,
			Condition{}.And(cond1, cond2, cond3).And(),
			And(cond1, cond2, cond3),
		},
		{
			`cond1.And(cond2)`,
			cond1.And(cond2),
			And(cond1, cond2),
		},
		{
			`cond1.And(cond2).And()`,
			cond1.And(cond2).And(),
			And(cond1, cond2),
		},
		{
			`cond1.And(cond2).And(cond3)`,
			cond1.And(cond2).And(cond3),
			And(cond1, cond2, cond3),
		},
		{
			`cond1.And(cond2).And(cond3).And()`,
			cond1.And(cond2).And(cond3).And(),
			And(cond1, cond2, cond3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestConditionOr(t *testing.T) {
	tests := []struct {
		Case      string
		Operation Condition
		Result    Condition
	}{
		{
			`Condition{}.Or()`,
			Condition{}.Or(),
			Or(),
		},
		{
			`Condition{}.Or(cond1)`,
			Condition{}.Or(cond1),
			cond1,
		},
		{
			`Condition{}.Or(cond1).Or()`,
			Condition{}.Or(cond1).Or(),
			cond1,
		},
		{
			`Condition{}.Or(cond1, cond2)`,
			Condition{}.Or(cond1, cond2),
			Or(cond1, cond2),
		},
		{
			`Condition{}.Or(cond1, cond2).Or()`,
			Condition{}.Or(cond1, cond2).Or(),
			Or(cond1, cond2),
		},
		{
			`Condition{}.Or(cond1, cond2, cond3)`,
			Condition{}.Or(cond1, cond2, cond3),
			Or(cond1, cond2, cond3),
		},
		{
			`Condition{}.Or(cond1, cond2, cond3).Or()`,
			Condition{}.Or(cond1, cond2, cond3).Or(),
			Or(cond1, cond2, cond3),
		},
		{
			`cond1.Or(cond2)`,
			cond1.Or(cond2),
			Or(cond1, cond2),
		},
		{
			`cond1.Or(cond2).Or()`,
			cond1.Or(cond2).Or(),
			Or(cond1, cond2),
		},
		{
			`cond1.Or(cond2).Or(cond3)`,
			cond1.Or(cond2).Or(cond3),
			Or(cond1, cond2, cond3),
		},
		{
			`cond1.Or(cond2).Or(cond3).Or()`,
			cond1.Or(cond2).Or(cond3).Or(),
			Or(cond1, cond2, cond3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Result, tt.Operation)
		})
	}
}

func TestConditionXor(t *testing.T) {
	tests := []struct {
		Case      string
		Operation Condition
		Result    Condition
	}{
		{
			`Condition{}.Xor()`,
			Condition{}.Xor(),
			Xor(),
		},
		{
			`Condition{}.Xor(cond1)`,
			Condition{}.Xor(cond1),
			cond1,
		},
		{
			`Condition{}.Xor(cond1).Xor()`,
			Condition{}.Xor(cond1).Xor(),
			cond1,
		},
		{
			`Condition{}.Xor(cond1, cond2)`,
			Condition{}.Xor(cond1, cond2),
			Xor(cond1, cond2),
		},
		{
			`Condition{}.Xor(cond1, cond2).Xor()`,
			Condition{}.Xor(cond1, cond2).Xor(),
			Xor(cond1, cond2),
		},
		{
			`Condition{}.Xor(cond1, cond2, cond3)`,
			Condition{}.Xor(cond1, cond2, cond3),
			Xor(cond1, cond2, cond3),
		},
		{
			`Condition{}.Xor(cond1, cond2, cond3).Xor()`,
			Condition{}.Xor(cond1, cond2, cond3).Xor(),
			Xor(cond1, cond2, cond3),
		},
		{
			`cond1.Xor(cond2)`,
			cond1.Xor(cond2),
			Xor(cond1, cond2),
		},
		{
			`cond1.Xor(cond2).Xor()`,
			cond1.Xor(cond2).Xor(),
			Xor(cond1, cond2),
		},
		{
			`cond1.Xor(cond2).Xor(cond3)`,
			cond1.Xor(cond2).Xor(cond3),
			Xor(cond1, cond2, cond3),
		},
		{
			`cond1.Xor(cond2).Xor(cond3).Xor()`,
			cond1.Xor(cond2).Xor(cond3).Xor(),
			Xor(cond1, cond2, cond3),
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
		Operation Condition
		Result    Condition
	}{
		{
			`And()`,
			And(),
			Condition{Type: ConditionAnd},
		},
		{
			`And(cond1)`,
			And(cond1),
			cond1,
		},
		{
			`And(cond1, cond2)`,
			And(cond1, cond2),
			Condition{
				Type:  ConditionAnd,
				Inner: []Condition{cond1, cond2},
			},
		},
		{
			`And(cond1, Or(cond2, cond3))`,
			And(cond1, Or(cond2, cond3)),
			Condition{
				Type: ConditionAnd,
				Inner: []Condition{
					cond1,
					Condition{
						Type:  ConditionOr,
						Inner: []Condition{cond2, cond3},
					},
				},
			},
		},
		{
			`And(Or(cond1, cond2), cond3)`,
			And(Or(cond1, cond2), cond3),
			Condition{
				Type: ConditionAnd,
				Inner: []Condition{
					Condition{
						Type:  ConditionOr,
						Inner: []Condition{cond1, cond2},
					},
					cond3,
				},
			},
		},
		{
			`And(Or(cond1, cond2), Or(cond3, cond4))`,
			And(Or(cond1, cond2), Or(cond3, cond4)),
			Condition{
				Type: ConditionAnd,
				Inner: []Condition{
					Condition{
						Type:  ConditionOr,
						Inner: []Condition{cond1, cond2},
					},
					Condition{
						Type:  ConditionOr,
						Inner: []Condition{cond3, cond4},
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
		Operation Condition
		Result    Condition
	}{
		{
			`Or()`,
			Or(),
			Condition{Type: ConditionOr},
		},
		{
			`Or(cond1)`,
			Or(cond1),
			cond1,
		},
		{
			`Or(cond1, cond2)`,
			Or(cond1, cond2),
			Condition{
				Type:  ConditionOr,
				Inner: []Condition{cond1, cond2},
			},
		},
		{
			`Or(cond1, And(cond2, cond3))`,
			Or(cond1, And(cond2, cond3)),
			Condition{
				Type: ConditionOr,
				Inner: []Condition{
					cond1,
					Condition{
						Type:  ConditionAnd,
						Inner: []Condition{cond2, cond3},
					},
				},
			},
		},
		{
			`Or(And(cond1, cond2), cond3)`,
			Or(And(cond1, cond2), cond3),
			Condition{
				Type: ConditionOr,
				Inner: []Condition{
					Condition{
						Type:  ConditionAnd,
						Inner: []Condition{cond1, cond2},
					},
					cond3,
				},
			},
		},
		{
			`Or(And(cond1, cond2), And(cond3, cond4))`,
			Or(And(cond1, cond2), And(cond3, cond4)),
			Condition{
				Type: ConditionOr,
				Inner: []Condition{
					Condition{
						Type:  ConditionAnd,
						Inner: []Condition{cond1, cond2},
					},
					Condition{
						Type:  ConditionAnd,
						Inner: []Condition{cond3, cond4},
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

func TestXor(t *testing.T) {
	tests := []struct {
		Case      string
		Operation Condition
		Result    Condition
	}{
		{
			`Xor()`,
			Xor(),
			Condition{Type: ConditionXor},
		},
		{
			`Xor(cond1)`,
			Xor(cond1),
			cond1,
		},
		{
			`Xor(cond1, cond2)`,
			Xor(cond1, cond2),
			Condition{
				Type:  ConditionXor,
				Inner: []Condition{cond1, cond2},
			},
		},
		{
			`Xor(cond1, And(cond2, cond3))`,
			Xor(cond1, And(cond2, cond3)),
			Condition{
				Type: ConditionXor,
				Inner: []Condition{
					cond1,
					Condition{
						Type:  ConditionAnd,
						Inner: []Condition{cond2, cond3},
					},
				},
			},
		},
		{
			`Xor(And(cond1, cond2), cond3)`,
			Xor(And(cond1, cond2), cond3),
			Condition{
				Type: ConditionXor,
				Inner: []Condition{
					Condition{
						Type:  ConditionAnd,
						Inner: []Condition{cond1, cond2},
					},
					cond3,
				},
			},
		},
		{
			`Xor(And(cond1, cond2), And(cond3, cond4))`,
			Xor(And(cond1, cond2), And(cond3, cond4)),
			Condition{
				Type: ConditionXor,
				Inner: []Condition{
					Condition{
						Type:  ConditionAnd,
						Inner: []Condition{cond1, cond2},
					},
					Condition{
						Type:  ConditionAnd,
						Inner: []Condition{cond3, cond4},
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

func TestNot(t *testing.T) {
	t.Skip("PENDING")
}

func TestEq(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionEq,
		Column: "field",
		Value:  "value",
	}, Eq("field", "value"))
}

func TestNe(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionNe,
		Column: "field",
		Value:  "value",
	}, Ne("field", "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionLt,
		Column: "field",
		Value:  10,
	}, Lt("field", 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionLte,
		Column: "field",
		Value:  10,
	}, Lte("field", 10))
}

func TestGt(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionGt,
		Column: "field",
		Value:  10,
	}, Gt("field", 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionGte,
		Column: "field",
		Value:  10,
	}, Gte("field", 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionNil,
		Column: "field",
	}, Nil("field"))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionNotNil,
		Column: "field",
	}, NotNil("field"))
}

func TestIn(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionIn,
		Column: "field",
		Value:  []string{"value"},
	}, In("field", []string{"value"}))
}

func TestNin(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionNin,
		Column: "field",
		Value:  []string{"value"},
	}, Nin("field", []string{"value"}))
}

func TestLike(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionLike,
		Column: "field",
		Expr:   "%expr%",
	}, Like("field", "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionNotLike,
		Column: "field",
		Expr:   "%expr%",
	}, NotLike("field", "%expr%"))
}

func TestFragment(t *testing.T) {
	assert.Equal(t, Condition{
		Type:   ConditionFragment,
		Column: "field",
		Value:  "value",
		Expr:   "expr",
	}, Fragment("field", "expr", "value"))
}

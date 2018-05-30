package c

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var cond1 = Eq(I("id"), 1)
var cond2 = Ne(I("name"), "foo")
var cond3 = Gt(I("score"), 80)
var cond4 = Lt(I("avg"), 10)

func TestCondition_None(t *testing.T) {
	assert.True(t, Condition{}.None())
	assert.True(t, And().None())
	assert.True(t, Not().None())

	assert.False(t, And(cond1).None())
	assert.False(t, And(cond1, cond2).None())
	assert.False(t, cond1.None())
}

func TestCondition_And(t *testing.T) {
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

func TestCondition_Or(t *testing.T) {
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
					{
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
					{
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
					{
						Type:  ConditionOr,
						Inner: []Condition{cond1, cond2},
					},
					{
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
					{
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
					{
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
					{
						Type:  ConditionAnd,
						Inner: []Condition{cond1, cond2},
					},
					{
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
	tests := []struct {
		Case     string
		Input    ConditionType
		Expected ConditionType
	}{
		{
			`Not Eq`,
			ConditionEq,
			ConditionNe,
		},
		{
			`Not Lt`,
			ConditionLt,
			ConditionGte,
		},
		{
			`Not Lte`,
			ConditionLte,
			ConditionGt,
		},
		{
			`Not Gt`,
			ConditionGt,
			ConditionLte,
		},
		{
			`Not Gte`,
			ConditionGte,
			ConditionLt,
		},
		{
			`Not Nil`,
			ConditionNil,
			ConditionNotNil,
		},
		{
			`Not In`,
			ConditionIn,
			ConditionNin,
		},
		{
			`Not Like`,
			ConditionLike,
			ConditionNotLike,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, Not(Condition{Type: tt.Input}).Type)
		})
	}
}

func TestEq(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionEq,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{"value"}},
	}, Eq(I("field"), "value"))
}

func TestNe(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionNe,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{"value"}},
	}, Ne(I("field"), "value"))
}

func TestLt(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionLt,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{10}},
	}, Lt(I("field"), 10))
}

func TestLte(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionLte,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{10}},
	}, Lte(I("field"), 10))
}

func TestGt(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionGt,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{10}},
	}, Gt(I("field"), 10))
}

func TestGte(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionGte,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{10}},
	}, Gte(I("field"), 10))
}

func TestNil(t *testing.T) {
	assert.Equal(t, Condition{
		Type: ConditionNil,
		Left: Operand{Column: I("field")},
	}, Nil(I("field")))
}

func TestNotNil(t *testing.T) {
	assert.Equal(t, Condition{
		Type: ConditionNotNil,
		Left: Operand{Column: I("field")},
	}, NotNil(I("field")))
}

func TestIn(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionIn,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{"value1", "value2"}},
	}, In(I("field"), "value1", "value2"))
}

func TestNin(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionNin,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{"value1", "value2"}},
	}, Nin(I("field"), "value1", "value2"))
}

func TestLike(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionLike,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{"%expr%"}},
	}, Like(I("field"), "%expr%"))
}

func TestNotLike(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionNotLike,
		Left:  Operand{Column: I("field")},
		Right: Operand{Values: []interface{}{"%expr%"}},
	}, NotLike(I("field"), "%expr%"))
}

func TestFragment(t *testing.T) {
	assert.Equal(t, Condition{
		Type:  ConditionFragment,
		Left:  Operand{Column: I("expr")},
		Right: Operand{Values: []interface{}{"value"}},
	}, Fragment(I("expr"), "value"))
}

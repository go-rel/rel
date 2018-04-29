// Package c defines function for building condition in query.
package c

// ConditionType defines enumeration of all supported condition types.
type ConditionType int

const (
	// ConditionAnd is condition type for and operator.
	ConditionAnd ConditionType = iota
	// ConditionOr is condition type for or operator.
	ConditionOr
	// ConditionNot is condition type for not operator.
	ConditionNot

	// ConditionEq is condition type for equal comparison.
	ConditionEq
	// ConditionNe is condition type for not equal comparison.
	ConditionNe

	// ConditionLt is condition type for less than comparison.
	ConditionLt
	// ConditionLte is condition type for less than or equal comparison.
	ConditionLte
	// ConditionGt is condition type for greater than comparison.
	ConditionGt
	// ConditionGte is condition type for greter than or equal comparison.
	ConditionGte

	// ConditionNil is condition type for nil check.
	ConditionNil
	// ConditionNotNil is condition type for not nil check.
	ConditionNotNil

	// ConditionIn is condition type for inclusion comparison.
	ConditionIn
	// ConditionNin is condition type for not inclusion comparison.
	ConditionNin

	// ConditionLike is condition type for like comparison.
	ConditionLike
	// ConditionNotLike is condition type for not like comparison.
	ConditionNotLike

	// ConditionFragment is condition type for custom condition.
	ConditionFragment
)

// I identifies database variable such as column name or table name.
type I string

// Operand defines information about condition's operand.
type Operand struct {
	Column I
	Values []interface{}
}

// NewOperand create new operand.
func NewOperand(o ...interface{}) Operand {
	if len(o) == 1 {
		if c, ok := o[0].(I); ok {
			return Operand{Column: c}
		}
	}

	return Operand{Values: o}
}

// Condition defines details of a coundition type.
type Condition struct {
	Type  ConditionType
	Left  Operand
	Right Operand
	Inner []Condition
}

// None returns true if no condition is specified.
func (c Condition) None() bool {
	return (c.Type == ConditionAnd ||
		c.Type == ConditionOr ||
		c.Type == ConditionNot) &&
		len(c.Inner) == 0
}

// And wraps conditions using and.
func (c Condition) And(condition ...Condition) Condition {
	if c.None() && len(condition) == 1 {
		return condition[0]
	} else if c.Type == ConditionAnd {
		c.Inner = append(c.Inner, condition...)
		return c
	}

	inner := append([]Condition{c}, condition...)
	return And(inner...)
}

// Or wraps conditions using or.
func (c Condition) Or(condition ...Condition) Condition {
	if c.None() && len(condition) == 1 {
		return condition[0]
	} else if c.Type == ConditionOr || c.None() {
		c.Type = ConditionOr
		c.Inner = append(c.Inner, condition...)
		return c
	}

	inner := append([]Condition{c}, condition...)
	return Or(inner...)
}

// And compares other conditions using and.
func And(inner ...Condition) Condition {
	if len(inner) == 1 {
		return inner[0]
	}

	return Condition{
		Type:  ConditionAnd,
		Inner: inner,
	}
}

// Or compares other conditions using and.
func Or(inner ...Condition) Condition {
	if len(inner) == 1 {
		return inner[0]
	}

	return Condition{
		Type:  ConditionOr,
		Inner: inner,
	}
}

// Not wraps conditions using not.
// It'll negate the condition type if possible.
func Not(inner ...Condition) Condition {
	if len(inner) == 1 {
		c := inner[0]
		switch c.Type {
		case ConditionEq:
			c.Type = ConditionNe
			return c
		case ConditionLt:
			c.Type = ConditionGte
			return c
		case ConditionLte:
			c.Type = ConditionGt
			return c
		case ConditionGt:
			c.Type = ConditionLte
			return c
		case ConditionGte:
			c.Type = ConditionLt
			return c
		case ConditionNil:
			c.Type = ConditionNotNil
			return c
		case ConditionIn:
			c.Type = ConditionNin
			return c
		case ConditionLike:
			c.Type = ConditionNotLike
			return c
		}
	}

	return Condition{
		Type:  ConditionNot,
		Inner: inner,
	}
}

// Eq compares that left value is equal to right value.
func Eq(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionEq,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

// Ne compares that left value is not equal to right value.
func Ne(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionNe,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

// Lt compares that left value is less than to right value.
func Lt(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionLt,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

// Lte compares that left value is less than or equal to right value.
func Lte(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionLte,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

// Gt compares that left value is greater than to right value.
func Gt(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionGt,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

// Gte compares that left value is greater than or equal to right value.
func Gte(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionGte,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

// Nil check whether column is nil.
func Nil(col I) Condition {
	return Condition{
		Type: ConditionNil,
		Left: NewOperand(col),
	}
}

// NotNil check whether column is not nil.
func NotNil(col I) Condition {
	return Condition{
		Type: ConditionNotNil,
		Left: NewOperand(col),
	}
}

// In check whethers value of the column is included in values.
func In(col I, values ...interface{}) Condition {
	return Condition{
		Type:  ConditionIn,
		Left:  NewOperand(col),
		Right: NewOperand(values...),
	}
}

// Nin check whethers value of the column is not included in values.
func Nin(col I, values ...interface{}) Condition {
	return Condition{
		Type:  ConditionNin,
		Left:  NewOperand(col),
		Right: NewOperand(values...),
	}
}

// Like compares value of column to match string pattern.
func Like(col I, pattern string) Condition {
	return Condition{
		Type:  ConditionLike,
		Left:  NewOperand(col),
		Right: NewOperand(pattern),
	}
}

// NotLike compares value of column to not match string pattern.
func NotLike(col I, pattern string) Condition {
	return Condition{
		Type:  ConditionNotLike,
		Left:  NewOperand(col),
		Right: NewOperand(pattern),
	}
}

// Fragment add custom condition.
func Fragment(expr I, values ...interface{}) Condition {
	return Condition{
		Type:  ConditionFragment,
		Left:  NewOperand(expr),
		Right: NewOperand(values...),
	}
}

// Package where defines function for building condition in query.
package where

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

// Condition defines details of a coundition type.
type Condition struct {
	Type   ConditionType
	Column string
	Values []interface{}
	Inner  []Condition
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

func (c Condition) and(other Condition) Condition {
	if c.Type == ConditionAnd {
		c.Inner = append(c.Inner, other)
		return c
	}

	return And(c, other)
}

func (c Condition) or(other Condition) Condition {
	if c.Type == ConditionOr {
		c.Inner = append(c.Inner, other)
		return c
	}

	return Or(c, other)
}

// AndEq append equal expression using and.
func (c Condition) AndEq(column string, value interface{}) Condition {
	return c.and(Eq(column, value))
}

// AndNe append not equal expression using and.
func (c Condition) AndNe(column string, value interface{}) Condition {
	return c.and(Ne(column, value))
}

// AndLt append lesser than expression using and.
func (c Condition) AndLt(column string, value interface{}) Condition {
	return c.and(Lt(column, value))
}

// AndLte append lesser than or equal expression using and.
func (c Condition) AndLte(column string, value interface{}) Condition {
	return c.and(Lte(column, value))
}

// AndGt append greater than expression using and.
func (c Condition) AndGt(column string, value interface{}) Condition {
	return c.and(Gt(column, value))
}

// AndGte append greater than or equal expression using and.
func (c Condition) AndGte(column string, value interface{}) Condition {
	return c.and(Gte(column, value))
}

// AndNil append is nil expression using and.
func (c Condition) AndNil(column string) Condition {
	return c.and(Nil(column))
}

// AndNotNil append is not nil expression using and.
func (c Condition) AndNotNil(column string) Condition {
	return c.and(NotNil(column))
}

// AndIn append is in expression using and.
func (c Condition) AndIn(column string, values ...interface{}) Condition {
	return c.and(In(column, values...))
}

// AndNin append is not in expression using and.
func (c Condition) AndNin(column string, values ...interface{}) Condition {
	return c.and(Nin(column, values...))
}

// AndLike append like expression using and.
func (c Condition) AndLike(column string, pattern string) Condition {
	return c.and(Like(column, pattern))
}

// AndNotLike append not like expression using and.
func (c Condition) AndNotLike(column string, pattern string) Condition {
	return c.and(Like(column, pattern))
}

// AndFragment append fragment using and.
func (c Condition) AndFragment(expr string, values ...interface{}) Condition {
	return c.and(Fragment(expr, values...))
}

// OrEq append equal expression using or.
func (c Condition) OrEq(column string, value interface{}) Condition {
	return c.or(Eq(column, value))
}

// OrNe append not equal expression using or.
func (c Condition) OrNe(column string, value interface{}) Condition {
	return c.or(Ne(column, value))
}

// OrLt append lesser than expression using or.
func (c Condition) OrLt(column string, value interface{}) Condition {
	return c.or(Lt(column, value))
}

// OrLte append lesser than or equal expression using or.
func (c Condition) OrLte(column string, value interface{}) Condition {
	return c.or(Lte(column, value))
}

// OrGt append greater than expression using or.
func (c Condition) OrGt(column string, value interface{}) Condition {
	return c.or(Gt(column, value))
}

// OrGte append greater than or equal expression using or.
func (c Condition) OrGte(column string, value interface{}) Condition {
	return c.or(Gte(column, value))
}

// OrNil append is nil expression using or.
func (c Condition) OrNil(column string) Condition {
	return c.or(Nil(column))
}

// OrNotNil append is not nil expression using or.
func (c Condition) OrNotNil(column string) Condition {
	return c.or(NotNil(column))
}

// OrIn append is in expression using or.
func (c Condition) OrIn(column string, values ...interface{}) Condition {
	return c.or(In(column, values...))
}

// OrNin append is not in expression using or.
func (c Condition) OrNin(column string, values ...interface{}) Condition {
	return c.or(Nin(column, values...))
}

// OrLike append like expression using or.
func (c Condition) OrLike(column string, pattern string) Condition {
	return c.or(Like(column, pattern))
}

// OrNotLike append not like expression using or.
func (c Condition) OrNotLike(column string, pattern string) Condition {
	return c.or(Like(column, pattern))
}

// OrFragment append fragment using or.
func (c Condition) OrFragment(expr string, values ...interface{}) Condition {
	return c.or(Fragment(expr, values...))
}

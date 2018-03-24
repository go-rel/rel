package c

type ConditionType int

const (
	ConditionAnd ConditionType = iota
	ConditionOr
	ConditionXor
	ConditionNot

	ConditionEq
	ConditionNe

	ConditionLt
	ConditionLte
	ConditionGt
	ConditionGte

	ConditionNil
	ConditionNotNil

	ConditionIn
	ConditionNin

	ConditionLike
	ConditionNotLike

	ConditionFragment
)

// column
type I string

type Operand struct {
	Column I
	Values []interface{}
}

func NewOperand(o ...interface{}) Operand {
	if len(o) == 1 {
		if c, ok := o[0].(I); ok {
			return Operand{Column: c}
		}
	}

	return Operand{Values: o}
}

type Condition struct {
	Type  ConditionType
	Left  Operand
	Right Operand
	Inner []Condition
}

func (c Condition) None() bool {
	return (c.Type == ConditionAnd ||
		c.Type == ConditionOr ||
		c.Type == ConditionXor ||
		c.Type == ConditionNot) &&
		len(c.Inner) == 0
}

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

func (c Condition) Xor(condition ...Condition) Condition {
	if c.None() && len(condition) == 1 {
		return condition[0]
	} else if c.Type == ConditionXor || c.None() {
		c.Type = ConditionXor
		c.Inner = append(c.Inner, condition...)
		return c
	}

	inner := append([]Condition{c}, condition...)
	return Xor(inner...)
}

func And(inner ...Condition) Condition {
	if len(inner) == 1 {
		return inner[0]
	}

	return Condition{
		Type:  ConditionAnd,
		Inner: inner,
	}
}

func Or(inner ...Condition) Condition {
	if len(inner) == 1 {
		return inner[0]
	}

	return Condition{
		Type:  ConditionOr,
		Inner: inner,
	}
}

func Xor(inner ...Condition) Condition {
	if len(inner) == 1 {
		return inner[0]
	}

	return Condition{
		Type:  ConditionXor,
		Inner: inner,
	}
}

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

func Eq(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionEq,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

func Ne(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionNe,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

func Lt(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionLt,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

func Lte(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionLte,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

func Gt(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionGt,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

func Gte(left, right interface{}) Condition {
	return Condition{
		Type:  ConditionGte,
		Left:  NewOperand(left),
		Right: NewOperand(right),
	}
}

func Nil(col I) Condition {
	return Condition{
		Type: ConditionNil,
		Left: NewOperand(col),
	}
}

func NotNil(col I) Condition {
	return Condition{
		Type: ConditionNotNil,
		Left: NewOperand(col),
	}
}

func In(col I, values ...interface{}) Condition {
	return Condition{
		Type:  ConditionIn,
		Left:  NewOperand(col),
		Right: NewOperand(values...),
	}
}

func Nin(col I, values ...interface{}) Condition {
	return Condition{
		Type:  ConditionNin,
		Left:  NewOperand(col),
		Right: NewOperand(values...),
	}
}

func Like(col I, pattern string) Condition {
	return Condition{
		Type:  ConditionLike,
		Left:  NewOperand(col),
		Right: NewOperand(pattern),
	}
}

func NotLike(col I, pattern string) Condition {
	return Condition{
		Type:  ConditionNotLike,
		Left:  NewOperand(col),
		Right: NewOperand(pattern),
	}
}

func Fragment(expr I, values ...interface{}) Condition {
	return Condition{
		Type:  ConditionFragment,
		Left:  NewOperand(expr),
		Right: NewOperand(values...),
	}
}

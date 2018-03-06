package query

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

type Condition struct {
	Type   ConditionType
	Column string
	Args   []interface{}
	Expr   string
	Inner  []Condition
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
		case ConditionLt:
			c.Type = ConditionGte
		case ConditionLte:
			c.Type = ConditionGt
		case ConditionGt:
			c.Type = ConditionLte
		case ConditionGte:
			c.Type = ConditionLt
		case ConditionNil:
			c.Type = ConditionNotNil
		case ConditionIn:
			c.Type = ConditionNin
		case ConditionLike:
			c.Type = ConditionNotLike
		}

		return c
	}

	return Condition{
		Type:  ConditionNot,
		Inner: inner,
	}
}

func Eq(col string, args interface{}) Condition {
	return Condition{
		Type:   ConditionEq,
		Column: col,
		Args:   []interface{}{args},
	}
}

func Ne(col string, args interface{}) Condition {
	return Condition{
		Type:   ConditionNe,
		Column: col,
		Args:   []interface{}{args},
	}
}

func Lt(col string, args interface{}) Condition {
	return Condition{
		Type:   ConditionLt,
		Column: col,
		Args:   []interface{}{args},
	}
}

func Lte(col string, args interface{}) Condition {
	return Condition{
		Type:   ConditionLte,
		Column: col,
		Args:   []interface{}{args},
	}
}

func Gt(col string, args interface{}) Condition {
	return Condition{
		Type:   ConditionGt,
		Column: col,
		Args:   []interface{}{args},
	}
}

func Gte(col string, args interface{}) Condition {
	return Condition{
		Type:   ConditionGte,
		Column: col,
		Args:   []interface{}{args},
	}
}

func Nil(col string) Condition {
	return Condition{
		Type:   ConditionNil,
		Column: col,
	}
}

func NotNil(col string) Condition {
	return Condition{
		Type:   ConditionNotNil,
		Column: col,
	}
}

func In(col string, args ...interface{}) Condition {
	return Condition{
		Type:   ConditionIn,
		Column: col,
		Args:   args,
	}
}

func Nin(col string, args ...interface{}) Condition {
	return Condition{
		Type:   ConditionNin,
		Column: col,
		Args:   args,
	}
}

func Like(col string, expr string) Condition {
	return Condition{
		Type:   ConditionLike,
		Column: col,
		Expr:   expr,
	}
}

func NotLike(col string, expr string) Condition {
	return Condition{
		Type:   ConditionNotLike,
		Column: col,
		Expr:   expr,
	}
}

func Fragment(col string, expr string, args ...interface{}) Condition {
	return Condition{
		Type:   ConditionFragment,
		Column: col,
		Args:   args,
		Expr:   expr,
	}
}

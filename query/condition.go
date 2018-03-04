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
	Value  interface{}
	Expr   string
	Inner  []Condition
}

func And(inner ...Condition) Condition {
	return Condition{
		Type:  ConditionAnd,
		Inner: inner,
	}
}

func Or(inner ...Condition) Condition {
	return Condition{
		Type:  ConditionOr,
		Inner: inner,
	}
}

func Xor(inner ...Condition) Condition {
	return Condition{
		Type:  ConditionXor,
		Inner: inner,
	}
}

func Not(inner ...Condition) Condition {
	return Condition{
		Type:  ConditionNot,
		Inner: inner,
	}
}

func Eq(col string, val interface{}) Condition {
	return Condition{
		Type:   ConditionEq,
		Column: col,
		Value:  val,
	}
}

func Ne(col string, val interface{}) Condition {
	return Condition{
		Type:   ConditionNe,
		Column: col,
		Value:  val,
	}
}

func Lt(col string, val interface{}) Condition {
	return Condition{
		Type:   ConditionLt,
		Column: col,
		Value:  val,
	}
}

func Lte(col string, val interface{}) Condition {
	return Condition{
		Type:   ConditionLte,
		Column: col,
		Value:  val,
	}
}

func Gt(col string, val interface{}) Condition {
	return Condition{
		Type:   ConditionGt,
		Column: col,
		Value:  val,
	}
}

func Gte(col string, val interface{}) Condition {
	return Condition{
		Type:   ConditionGte,
		Column: col,
		Value:  val,
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

func In(col string, val interface{}) Condition {
	return Condition{
		Type:   ConditionIn,
		Column: col,
		Value:  val,
	}
}

func Nin(col string, val interface{}) Condition {
	return Condition{
		Type:   ConditionNin,
		Column: col,
		Value:  val,
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

func Fragment(col string, expr string, val interface{}) Condition {
	return Condition{
		Type:   ConditionFrag,
		Column: col,
		Value:  pattern,
		Expr:   expr,
	}
}

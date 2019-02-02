package where

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
		default:
			return Condition{
				Type:  ConditionNot,
				Inner: inner,
			}
		}

		return c
	}

	return Condition{
		Type:  ConditionNot,
		Inner: inner,
	}
}

// Eq expression column equal to value.
func Eq(column string, value interface{}) Condition {
	return Condition{
		Type:   ConditionEq,
		Column: column,
		Values: []interface{}{value},
	}
}

// Ne compares that left value is not equal to right value.
func Ne(column string, value interface{}) Condition {
	return Condition{
		Type:   ConditionNe,
		Column: column,
		Values: []interface{}{value},
	}
}

// Lt compares that left value is less than to right value.
func Lt(column string, value interface{}) Condition {
	return Condition{
		Type:   ConditionLt,
		Column: column,
		Values: []interface{}{value},
	}
}

// Lte compares that left value is less than or equal to right value.
func Lte(column string, value interface{}) Condition {
	return Condition{
		Type:   ConditionLte,
		Column: column,
		Values: []interface{}{value},
	}
}

// Gt compares that left value is greater than to right value.
func Gt(column string, value interface{}) Condition {
	return Condition{
		Type:   ConditionGt,
		Column: column,
		Values: []interface{}{value},
	}
}

// Gte compares that left value is greater than or equal to right value.
func Gte(column string, value interface{}) Condition {
	return Condition{
		Type:   ConditionGte,
		Column: column,
		Values: []interface{}{value},
	}
}

// Nil check whether column is nil.
func Nil(column string) Condition {
	return Condition{
		Type:   ConditionNil,
		Column: column,
	}
}

// NotNil check whether column is not nil.
func NotNil(column string) Condition {
	return Condition{
		Type:   ConditionNotNil,
		Column: column,
	}
}

// In check whethers value of the column is included in values.
func In(column string, values ...interface{}) Condition {
	return Condition{
		Type:   ConditionIn,
		Column: column,
		Values: values,
	}
}

// Nin check whethers value of the column is not included in values.
func Nin(column string, values ...interface{}) Condition {
	return Condition{
		Type:   ConditionNin,
		Column: column,
		Values: values,
	}
}

// Like compares value of column to match string pattern.
func Like(column string, pattern string) Condition {
	return Condition{
		Type:   ConditionLike,
		Column: column,
		Values: []interface{}{pattern},
	}
}

// NotLike compares value of column to not match string pattern.
func NotLike(column string, pattern string) Condition {
	return Condition{
		Type:   ConditionNotLike,
		Column: column,
		Values: []interface{}{pattern},
	}
}

// Fragment add custom condition.
func Fragment(expr string, values ...interface{}) Condition {
	return Condition{
		Type:   ConditionFragment,
		Column: expr,
		Values: values,
	}
}

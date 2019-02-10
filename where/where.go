// Package where is syntatic sugar for building where query.
package where

import (
	"github.com/Fs02/grimoire/query"
)

// And compares other filters using and.
func And(inner ...query.FilterClause) query.FilterClause {
	return query.FilterAnd(inner...)
}

// Or compares other filters using and.
func Or(inner ...query.FilterClause) query.FilterClause {
	return query.FilterOr(inner...)
}

// Not wraps filters using not.
// It'll negate the filter type if possible.
func Not(inner ...query.FilterClause) query.FilterClause {
	return query.FilterNot(inner...)
}

// Eq expression field equal to value.
func Eq(field string, value interface{}) query.FilterClause {
	return query.FilterEq(field, value)
}

// Ne compares that left value is not equal to right value.
func Ne(field string, value interface{}) query.FilterClause {
	return query.FilterNe(field, value)
}

// Lt compares that left value is less than to right value.
func Lt(field string, value interface{}) query.FilterClause {
	return query.FilterLt(field, value)
}

// Lte compares that left value is less than or equal to right value.
func Lte(field string, value interface{}) query.FilterClause {
	return query.FilterLte(field, value)
}

// Gt compares that left value is greater than to right value.
func Gt(field string, value interface{}) query.FilterClause {
	return query.FilterGt(field, value)
}

// Gte compares that left value is greater than or equal to right value.
func Gte(field string, value interface{}) query.FilterClause {
	return query.FilterGte(field, value)
}

// Nil check whether field is nil.
func Nil(field string) query.FilterClause {
	return query.FilterNil(field)
}

// NotNil check whether field is not nil.
func NotNil(field string) query.FilterClause {
	return query.FilterNotNil(field)
}

// In check whethers value of the field is included in values.
func In(field string, values ...interface{}) query.FilterClause {
	return query.FilterIn(field, values...)
}

// Nin check whethers value of the field is not included in values.
func Nin(field string, values ...interface{}) query.FilterClause {
	return query.FilterNin(field, values...)
}

// Like compares value of field to match string pattern.
func Like(field string, pattern string) query.FilterClause {
	return query.FilterLike(field, pattern)
}

// NotLike compares value of field to not match string pattern.
func NotLike(field string, pattern string) query.FilterClause {
	return query.FilterNotLike(field, pattern)
}

// Fragment add custom filter.
func Fragment(expr string, values ...interface{}) query.FilterClause {
	return query.FilterFragment(expr, values...)
}

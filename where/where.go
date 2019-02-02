// Package where is syntatic sugar for building where query.
package where

import (
	"github.com/Fs02/grimoire/query"
)

// And compares other filters using and.
func And(inner ...query.Filter) query.Filter {
	return query.FilterAnd(inner...)
}

// Or compares other filters using and.
func Or(inner ...query.Filter) query.Filter {
	return query.FilterOr(inner...)
}

// Not wraps filters using not.
// It'll negate the filter type if possible.
func Not(inner ...query.Filter) query.Filter {
	return query.FilterNot(inner...)
}

// Eq expression column equal to value.
func Eq(column string, value interface{}) query.Filter {
	return query.FilterEq(column, value)
}

// Ne compares that left value is not equal to right value.
func Ne(column string, value interface{}) query.Filter {
	return query.FilterNe(column, value)
}

// Lt compares that left value is less than to right value.
func Lt(column string, value interface{}) query.Filter {
	return query.FilterLt(column, value)
}

// Lte compares that left value is less than or equal to right value.
func Lte(column string, value interface{}) query.Filter {
	return query.FilterLte(column, value)
}

// Gt compares that left value is greater than to right value.
func Gt(column string, value interface{}) query.Filter {
	return query.FilterGt(column, value)
}

// Gte compares that left value is greater than or equal to right value.
func Gte(column string, value interface{}) query.Filter {
	return query.FilterGte(column, value)
}

// Nil check whether column is nil.
func Nil(column string) query.Filter {
	return query.FilterNil(column)
}

// NotNil check whether column is not nil.
func NotNil(column string) query.Filter {
	return query.FilterNotNil(column)
}

// In check whethers value of the column is included in values.
func In(column string, values ...interface{}) query.Filter {
	return query.FilterIn(column, values...)
}

// Nin check whethers value of the column is not included in values.
func Nin(column string, values ...interface{}) query.Filter {
	return query.FilterNin(column, values...)
}

// Like compares value of column to match string pattern.
func Like(column string, pattern string) query.Filter {
	return query.FilterLike(column, pattern)
}

// NotLike compares value of column to not match string pattern.
func NotLike(column string, pattern string) query.Filter {
	return query.FilterNotLike(column, pattern)
}

// Fragment add custom filter.
func Fragment(expr string, values ...interface{}) query.Filter {
	return query.FilterFragment(expr, values...)
}

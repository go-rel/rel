// Package where is syntatic sugar for building where query.
package where

import (
	"github.com/Fs02/grimoire"
)

// And compares other filters using and.
func And(inner ...grimoire.FilterQuery) grimoire.FilterQuery {
	return grimoire.And(inner...)
}

// Or compares other filters using and.
func Or(inner ...grimoire.FilterQuery) grimoire.FilterQuery {
	return grimoire.Or(inner...)
}

// Not wraps filters using not.
// It'll negate the filter type if possible.
func Not(inner ...grimoire.FilterQuery) grimoire.FilterQuery {
	return grimoire.Not(inner...)
}

// Eq expression field equal to value.
func Eq(field string, value interface{}) grimoire.FilterQuery {
	return grimoire.Eq(field, value)
}

// Ne compares that left value is not equal to right value.
func Ne(field string, value interface{}) grimoire.FilterQuery {
	return grimoire.Ne(field, value)
}

// Lt compares that left value is less than to right value.
func Lt(field string, value interface{}) grimoire.FilterQuery {
	return grimoire.Lt(field, value)
}

// Lte compares that left value is less than or equal to right value.
func Lte(field string, value interface{}) grimoire.FilterQuery {
	return grimoire.Lte(field, value)
}

// Gt compares that left value is greater than to right value.
func Gt(field string, value interface{}) grimoire.FilterQuery {
	return grimoire.Gt(field, value)
}

// Gte compares that left value is greater than or equal to right value.
func Gte(field string, value interface{}) grimoire.FilterQuery {
	return grimoire.Gte(field, value)
}

// Nil check whether field is nil.
func Nil(field string) grimoire.FilterQuery {
	return grimoire.Nil(field)
}

// NotNil check whether field is not nil.
func NotNil(field string) grimoire.FilterQuery {
	return grimoire.NotNil(field)
}

// In check whethers value of the field is included in values.
func In(field string, values ...interface{}) grimoire.FilterQuery {
	return grimoire.In(field, values...)
}

// Nin check whethers value of the field is not included in values.
func Nin(field string, values ...interface{}) grimoire.FilterQuery {
	return grimoire.Nin(field, values...)
}

// Like compares value of field to match string pattern.
func Like(field string, pattern string) grimoire.FilterQuery {
	return grimoire.Like(field, pattern)
}

// NotLike compares value of field to not match string pattern.
func NotLike(field string, pattern string) grimoire.FilterQuery {
	return grimoire.NotLike(field, pattern)
}

// Fragment add custom filter.
func Fragment(expr string, values ...interface{}) grimoire.FilterQuery {
	return grimoire.FilterFragment(expr, values...)
}

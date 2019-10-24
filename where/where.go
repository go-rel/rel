// Package where is syntatic sugar for building where query.
package where

import (
	"github.com/Fs02/rel"
)

// And compares other filters using and.
func And(inner ...rel.FilterQuery) rel.FilterQuery {
	return rel.And(inner...)
}

// Or compares other filters using and.
func Or(inner ...rel.FilterQuery) rel.FilterQuery {
	return rel.Or(inner...)
}

// Not wraps filters using not.
// It'll negate the filter type if possible.
func Not(inner ...rel.FilterQuery) rel.FilterQuery {
	return rel.Not(inner...)
}

// Eq expression field equal to value.
func Eq(field string, value interface{}) rel.FilterQuery {
	return rel.Eq(field, value)
}

// Ne compares that left value is not equal to right value.
func Ne(field string, value interface{}) rel.FilterQuery {
	return rel.Ne(field, value)
}

// Lt compares that left value is less than to right value.
func Lt(field string, value interface{}) rel.FilterQuery {
	return rel.Lt(field, value)
}

// Lte compares that left value is less than or equal to right value.
func Lte(field string, value interface{}) rel.FilterQuery {
	return rel.Lte(field, value)
}

// Gt compares that left value is greater than to right value.
func Gt(field string, value interface{}) rel.FilterQuery {
	return rel.Gt(field, value)
}

// Gte compares that left value is greater than or equal to right value.
func Gte(field string, value interface{}) rel.FilterQuery {
	return rel.Gte(field, value)
}

// Nil check whether field is nil.
func Nil(field string) rel.FilterQuery {
	return rel.Nil(field)
}

// NotNil check whether field is not nil.
func NotNil(field string) rel.FilterQuery {
	return rel.NotNil(field)
}

// In check whethers value of the field is included in values.
func In(field string, values ...interface{}) rel.FilterQuery {
	return rel.In(field, values...)
}

// InInt check whethers integer value of the field is included in values.
func InInt(field string, values []int) rel.FilterQuery {
	return rel.InInt(field, values)
}

// InUint check whethers unsigned integer value of the field is included in values.
func InUint(field string, values []uint) rel.FilterQuery {
	return rel.InUint(field, values)
}

// InString check whethers string value of the field is included in values.
func InString(field string, values []string) rel.FilterQuery {
	return rel.InString(field, values)
}

// Nin check whethers value of the field is not included in values.
func Nin(field string, values ...interface{}) rel.FilterQuery {
	return rel.Nin(field, values...)
}

// NinInt check whethers integer value of the field is not included in values.
func NinInt(field string, values []int) rel.FilterQuery {
	return rel.NinInt(field, values)
}

// NinUint check whethers unsigned integer value of the field is not included in values.
func NinUint(field string, values []uint) rel.FilterQuery {
	return rel.NinUint(field, values)
}

// NinString check whethers string value of the field is not included in values.
func NinString(field string, values []string) rel.FilterQuery {
	return rel.NinString(field, values)
}

// Like compares value of field to match string pattern.
func Like(field string, pattern string) rel.FilterQuery {
	return rel.Like(field, pattern)
}

// NotLike compares value of field to not match string pattern.
func NotLike(field string, pattern string) rel.FilterQuery {
	return rel.NotLike(field, pattern)
}

// Fragment add custom filter.
func Fragment(expr string, values ...interface{}) rel.FilterQuery {
	return rel.FilterFragment(expr, values...)
}

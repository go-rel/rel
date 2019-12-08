// Package where is syntatic sugar for building where query.
package where

import (
	"github.com/Fs02/rel"
)

var (
	// And compares other filters using and.
	And = rel.And

	// Or compares other filters using and.
	Or = rel.Or

	// Not wraps filters using not.
	// It'll negate the filter type if possible.
	Not = rel.Not

	// Eq expression field equal to value.
	Eq = rel.Eq

	// Ne compares that left value is not equal to right value.
	Ne = rel.Ne

	// Lt compares that left value is less than to right value.
	Lt = rel.Lt

	// Lte compares that left value is less than or equal to right value.
	Lte = rel.Lte

	// Gt compares that left value is greater than to right value.
	Gt = rel.Gt

	// Gte compares that left value is greater than or equal to right value.
	Gte = rel.Gte

	// Nil check whether field is nil.
	Nil = rel.Nil

	// NotNil check whether field is not nil.
	NotNil = rel.NotNil

	// In check whethers value of the field is included in values.
	In = rel.In

	// InInt check whethers integer value of the field is included in values.
	InInt = rel.InInt

	// InUint check whethers unsigned integer value of the field is included in values.
	InUint = rel.InUint

	// InString check whethers string value of the field is included in values.
	InString = rel.InString

	// Nin check whethers value of the field is not included in values.
	Nin = rel.Nin

	// NinInt check whethers integer value of the field is not included in values.
	NinInt = rel.NinInt

	// NinUint check whethers unsigned integer value of the field is not included in values.
	NinUint = rel.NinUint

	// NinString check whethers string value of the field is not included in values.
	NinString = rel.NinString

	// Like compares value of field to match string pattern.
	Like = rel.Like

	// NotLike compares value of field to not match string pattern.
	NotLike = rel.NotLike

	// Fragment add custom filter.
	Fragment = rel.FilterFragment
)

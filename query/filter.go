// Package where defines function for building filter in query.
package query

// FilterOp defines enumeration of all supported filter types.
type FilterOp int

const (
	// AndOp is filter type for and operator.
	AndOp FilterOp = iota
	// OrOp is filter type for or operator.
	OrOp
	// NotOp is filter type for not operator.
	NotOp

	// EqOp is filter type for equal comparison.
	EqOp
	// NeOp is filter type for not equal comparison.
	NeOp

	// LtOp is filter type for less than comparison.
	LtOp
	// LteOp is filter type for less than or equal comparison.
	LteOp
	// GtOp is filter type for greater than comparison.
	GtOp
	// GteOp is filter type for greter than or equal comparison.
	GteOp

	// NilOp is filter type for nil check.
	NilOp
	// NotNilOp is filter type for not nil check.
	NotNilOp

	// InOp is filter type for inclusion comparison.
	InOp
	// NinOp is filter type for not inclusion comparison.
	NinOp

	// LikeOp is filter type for like comparison.
	LikeOp
	// NotLikeOp is filter type for not like comparison.
	NotLikeOp

	// FragmentOp is filter type for custom filter.
	FragmentOp
)

// Filter defines details of a coundition type.
type Filter struct {
	Type   FilterOp
	Field  string
	Values []interface{}
	Inner  []Filter
}

// None returns true if no filter is specified.
func (f Filter) None() bool {
	return (f.Type == AndOp ||
		f.Type == OrOp ||
		f.Type == NotOp) &&
		len(f.Inner) == 0
}

// And wraps filters using and.
func (f Filter) And(filters ...Filter) Filter {
	if f.None() && len(filters) == 1 {
		return filters[0]
	} else if f.Type == AndOp {
		f.Inner = append(f.Inner, filters...)
		return f
	}

	inner := append([]Filter{f}, filters...)
	return FilterAnd(inner...)
}

// Or wraps filters using or.
func (f Filter) Or(filter ...Filter) Filter {
	if f.None() && len(filter) == 1 {
		return filter[0]
	} else if f.Type == OrOp || f.None() {
		f.Type = OrOp
		f.Inner = append(f.Inner, filter...)
		return f
	}

	inner := append([]Filter{f}, filter...)
	return FilterOr(inner...)
}

func (f Filter) and(other Filter) Filter {
	if f.Type == AndOp {
		f.Inner = append(f.Inner, other)
		return f
	}

	return FilterAnd(f, other)
}

func (f Filter) or(other Filter) Filter {
	if f.Type == OrOp || f.None() {
		f.Type = OrOp
		f.Inner = append(f.Inner, other)
		return f
	}

	return FilterOr(f, other)
}

// AndEq append equal expression using and.
func (f Filter) AndEq(field string, value interface{}) Filter {
	return f.and(FilterEq(field, value))
}

// AndNe append not equal expression using and.
func (f Filter) AndNe(field string, value interface{}) Filter {
	return f.and(FilterNe(field, value))
}

// AndLt append lesser than expression using and.
func (f Filter) AndLt(field string, value interface{}) Filter {
	return f.and(FilterLt(field, value))
}

// AndLte append lesser than or equal expression using and.
func (f Filter) AndLte(field string, value interface{}) Filter {
	return f.and(FilterLte(field, value))
}

// AndGt append greater than expression using and.
func (f Filter) AndGt(field string, value interface{}) Filter {
	return f.and(FilterGt(field, value))
}

// AndGte append greater than or equal expression using and.
func (f Filter) AndGte(field string, value interface{}) Filter {
	return f.and(FilterGte(field, value))
}

// AndNil append is nil expression using and.
func (f Filter) AndNil(field string) Filter {
	return f.and(FilterNil(field))
}

// AndNotNil append is not nil expression using and.
func (f Filter) AndNotNil(field string) Filter {
	return f.and(FilterNotNil(field))
}

// AndIn append is in expression using and.
func (f Filter) AndIn(field string, values ...interface{}) Filter {
	return f.and(FilterIn(field, values...))
}

// AndNin append is not in expression using and.
func (f Filter) AndNin(field string, values ...interface{}) Filter {
	return f.and(FilterNin(field, values...))
}

// AndLike append like expression using and.
func (f Filter) AndLike(field string, pattern string) Filter {
	return f.and(FilterLike(field, pattern))
}

// AndNotLike append not like expression using and.
func (f Filter) AndNotLike(field string, pattern string) Filter {
	return f.and(FilterNotLike(field, pattern))
}

// AndFragment append fragment using and.
func (f Filter) AndFragment(expr string, values ...interface{}) Filter {
	return f.and(FilterFragment(expr, values...))
}

// OrEq append equal expression using or.
func (f Filter) OrEq(field string, value interface{}) Filter {
	return f.or(FilterEq(field, value))
}

// OrNe append not equal expression using or.
func (f Filter) OrNe(field string, value interface{}) Filter {
	return f.or(FilterNe(field, value))
}

// OrLt append lesser than expression using or.
func (f Filter) OrLt(field string, value interface{}) Filter {
	return f.or(FilterLt(field, value))
}

// OrLte append lesser than or equal expression using or.
func (f Filter) OrLte(field string, value interface{}) Filter {
	return f.or(FilterLte(field, value))
}

// OrGt append greater than expression using or.
func (f Filter) OrGt(field string, value interface{}) Filter {
	return f.or(FilterGt(field, value))
}

// OrGte append greater than or equal expression using or.
func (f Filter) OrGte(field string, value interface{}) Filter {
	return f.or(FilterGte(field, value))
}

// OrNil append is nil expression using or.
func (f Filter) OrNil(field string) Filter {
	return f.or(FilterNil(field))
}

// OrNotNil append is not nil expression using or.
func (f Filter) OrNotNil(field string) Filter {
	return f.or(FilterNotNil(field))
}

// OrIn append is in expression using or.
func (f Filter) OrIn(field string, values ...interface{}) Filter {
	return f.or(FilterIn(field, values...))
}

// OrNin append is not in expression using or.
func (f Filter) OrNin(field string, values ...interface{}) Filter {
	return f.or(FilterNin(field, values...))
}

// OrLike append like expression using or.
func (f Filter) OrLike(field string, pattern string) Filter {
	return f.or(FilterLike(field, pattern))
}

// OrNotLike append not like expression using or.
func (f Filter) OrNotLike(field string, pattern string) Filter {
	return f.or(FilterNotLike(field, pattern))
}

// OrFragment append fragment using or.
func (f Filter) OrFragment(expr string, values ...interface{}) Filter {
	return f.or(FilterFragment(expr, values...))
}

// FilterAnd compares other filters using and.
func FilterAnd(inner ...Filter) Filter {
	if len(inner) == 1 {
		return inner[0]
	}

	return Filter{
		Type:  AndOp,
		Inner: inner,
	}
}

// FilterOr compares other filters using and.
func FilterOr(inner ...Filter) Filter {
	if len(inner) == 1 {
		return inner[0]
	}

	return Filter{
		Type:  OrOp,
		Inner: inner,
	}
}

// FilterNot wraps filters using not.
// It'll negate the filter type if possible.
func FilterNot(inner ...Filter) Filter {
	if len(inner) == 1 {
		f := inner[0]
		switch f.Type {
		case EqOp:
			f.Type = NeOp
			return f
		case LtOp:
			f.Type = GteOp
		case LteOp:
			f.Type = GtOp
		case GtOp:
			f.Type = LteOp
		case GteOp:
			f.Type = LtOp
		case NilOp:
			f.Type = NotNilOp
		case InOp:
			f.Type = NinOp
		case LikeOp:
			f.Type = NotLikeOp
		default:
			return Filter{
				Type:  NotOp,
				Inner: inner,
			}
		}

		return f
	}

	return Filter{
		Type:  NotOp,
		Inner: inner,
	}
}

// FilterEq expression field equal to value.
func FilterEq(field string, value interface{}) Filter {
	return Filter{
		Type:   EqOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterNe compares that left value is not equal to right value.
func FilterNe(field string, value interface{}) Filter {
	return Filter{
		Type:   NeOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterLt compares that left value is less than to right value.
func FilterLt(field string, value interface{}) Filter {
	return Filter{
		Type:   LtOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterLte compares that left value is less than or equal to right value.
func FilterLte(field string, value interface{}) Filter {
	return Filter{
		Type:   LteOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterGt compares that left value is greater than to right value.
func FilterGt(field string, value interface{}) Filter {
	return Filter{
		Type:   GtOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterGte compares that left value is greater than or equal to right value.
func FilterGte(field string, value interface{}) Filter {
	return Filter{
		Type:   GteOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterNil check whether field is nil.
func FilterNil(field string) Filter {
	return Filter{
		Type:  NilOp,
		Field: field,
	}
}

// FilterNotNil check whether field is not nil.
func FilterNotNil(field string) Filter {
	return Filter{
		Type:  NotNilOp,
		Field: field,
	}
}

// FilterIn check whethers value of the field is included in values.
func FilterIn(field string, values ...interface{}) Filter {
	return Filter{
		Type:   InOp,
		Field:  field,
		Values: values,
	}
}

// FilterNin check whethers value of the field is not included in values.
func FilterNin(field string, values ...interface{}) Filter {
	return Filter{
		Type:   NinOp,
		Field:  field,
		Values: values,
	}
}

// FilterLike compares value of field to match string pattern.
func FilterLike(field string, pattern string) Filter {
	return Filter{
		Type:   LikeOp,
		Field:  field,
		Values: []interface{}{pattern},
	}
}

// FilterNotLike compares value of field to not match string pattern.
func FilterNotLike(field string, pattern string) Filter {
	return Filter{
		Type:   NotLikeOp,
		Field:  field,
		Values: []interface{}{pattern},
	}
}

// FilterFragment add custom filter.
func FilterFragment(expr string, values ...interface{}) Filter {
	return Filter{
		Type:   FragmentOp,
		Field:  expr,
		Values: values,
	}
}

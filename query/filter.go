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
	Column string
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
func (f Filter) AndEq(column string, value interface{}) Filter {
	return f.and(FilterEq(column, value))
}

// AndNe append not equal expression using and.
func (f Filter) AndNe(column string, value interface{}) Filter {
	return f.and(FilterNe(column, value))
}

// AndLt append lesser than expression using and.
func (f Filter) AndLt(column string, value interface{}) Filter {
	return f.and(FilterLt(column, value))
}

// AndLte append lesser than or equal expression using and.
func (f Filter) AndLte(column string, value interface{}) Filter {
	return f.and(FilterLte(column, value))
}

// AndGt append greater than expression using and.
func (f Filter) AndGt(column string, value interface{}) Filter {
	return f.and(FilterGt(column, value))
}

// AndGte append greater than or equal expression using and.
func (f Filter) AndGte(column string, value interface{}) Filter {
	return f.and(FilterGte(column, value))
}

// AndNil append is nil expression using and.
func (f Filter) AndNil(column string) Filter {
	return f.and(FilterNil(column))
}

// AndNotNil append is not nil expression using and.
func (f Filter) AndNotNil(column string) Filter {
	return f.and(FilterNotNil(column))
}

// AndIn append is in expression using and.
func (f Filter) AndIn(column string, values ...interface{}) Filter {
	return f.and(FilterIn(column, values...))
}

// AndNin append is not in expression using and.
func (f Filter) AndNin(column string, values ...interface{}) Filter {
	return f.and(FilterNin(column, values...))
}

// AndLike append like expression using and.
func (f Filter) AndLike(column string, pattern string) Filter {
	return f.and(FilterLike(column, pattern))
}

// AndNotLike append not like expression using and.
func (f Filter) AndNotLike(column string, pattern string) Filter {
	return f.and(FilterNotLike(column, pattern))
}

// AndFragment append fragment using and.
func (f Filter) AndFragment(expr string, values ...interface{}) Filter {
	return f.and(FilterFragment(expr, values...))
}

// OrEq append equal expression using or.
func (f Filter) OrEq(column string, value interface{}) Filter {
	return f.or(FilterEq(column, value))
}

// OrNe append not equal expression using or.
func (f Filter) OrNe(column string, value interface{}) Filter {
	return f.or(FilterNe(column, value))
}

// OrLt append lesser than expression using or.
func (f Filter) OrLt(column string, value interface{}) Filter {
	return f.or(FilterLt(column, value))
}

// OrLte append lesser than or equal expression using or.
func (f Filter) OrLte(column string, value interface{}) Filter {
	return f.or(FilterLte(column, value))
}

// OrGt append greater than expression using or.
func (f Filter) OrGt(column string, value interface{}) Filter {
	return f.or(FilterGt(column, value))
}

// OrGte append greater than or equal expression using or.
func (f Filter) OrGte(column string, value interface{}) Filter {
	return f.or(FilterGte(column, value))
}

// OrNil append is nil expression using or.
func (f Filter) OrNil(column string) Filter {
	return f.or(FilterNil(column))
}

// OrNotNil append is not nil expression using or.
func (f Filter) OrNotNil(column string) Filter {
	return f.or(FilterNotNil(column))
}

// OrIn append is in expression using or.
func (f Filter) OrIn(column string, values ...interface{}) Filter {
	return f.or(FilterIn(column, values...))
}

// OrNin append is not in expression using or.
func (f Filter) OrNin(column string, values ...interface{}) Filter {
	return f.or(FilterNin(column, values...))
}

// OrLike append like expression using or.
func (f Filter) OrLike(column string, pattern string) Filter {
	return f.or(FilterLike(column, pattern))
}

// OrNotLike append not like expression using or.
func (f Filter) OrNotLike(column string, pattern string) Filter {
	return f.or(FilterNotLike(column, pattern))
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

// FilterEq expression column equal to value.
func FilterEq(column string, value interface{}) Filter {
	return Filter{
		Type:   EqOp,
		Column: column,
		Values: []interface{}{value},
	}
}

// FilterNe compares that left value is not equal to right value.
func FilterNe(column string, value interface{}) Filter {
	return Filter{
		Type:   NeOp,
		Column: column,
		Values: []interface{}{value},
	}
}

// FilterLt compares that left value is less than to right value.
func FilterLt(column string, value interface{}) Filter {
	return Filter{
		Type:   LtOp,
		Column: column,
		Values: []interface{}{value},
	}
}

// FilterLte compares that left value is less than or equal to right value.
func FilterLte(column string, value interface{}) Filter {
	return Filter{
		Type:   LteOp,
		Column: column,
		Values: []interface{}{value},
	}
}

// FilterGt compares that left value is greater than to right value.
func FilterGt(column string, value interface{}) Filter {
	return Filter{
		Type:   GtOp,
		Column: column,
		Values: []interface{}{value},
	}
}

// FilterGte compares that left value is greater than or equal to right value.
func FilterGte(column string, value interface{}) Filter {
	return Filter{
		Type:   GteOp,
		Column: column,
		Values: []interface{}{value},
	}
}

// FilterNil check whether column is nil.
func FilterNil(column string) Filter {
	return Filter{
		Type:   NilOp,
		Column: column,
	}
}

// FilterNotNil check whether column is not nil.
func FilterNotNil(column string) Filter {
	return Filter{
		Type:   NotNilOp,
		Column: column,
	}
}

// FilterIn check whethers value of the column is included in values.
func FilterIn(column string, values ...interface{}) Filter {
	return Filter{
		Type:   InOp,
		Column: column,
		Values: values,
	}
}

// FilterNin check whethers value of the column is not included in values.
func FilterNin(column string, values ...interface{}) Filter {
	return Filter{
		Type:   NinOp,
		Column: column,
		Values: values,
	}
}

// FilterLike compares value of column to match string pattern.
func FilterLike(column string, pattern string) Filter {
	return Filter{
		Type:   LikeOp,
		Column: column,
		Values: []interface{}{pattern},
	}
}

// FilterNotLike compares value of column to not match string pattern.
func FilterNotLike(column string, pattern string) Filter {
	return Filter{
		Type:   NotLikeOp,
		Column: column,
		Values: []interface{}{pattern},
	}
}

// FilterFragment add custom filter.
func FilterFragment(expr string, values ...interface{}) Filter {
	return Filter{
		Type:   FragmentOp,
		Column: expr,
		Values: values,
	}
}

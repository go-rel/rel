package grimoire

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

// FilterClause defines details of a coundition type.
type FilterClause struct {
	Type   FilterOp
	Field  string
	Values []interface{} // TODO: don't use array (let adapter cast when needed)
	Inner  []FilterClause
}

// Build Filter query.
func (f FilterClause) Build(query *Query) {
	query.WhereClause = query.WhereClause.And(f)
}

// None returns true if no filter is specified.
func (f FilterClause) None() bool {
	return (f.Type == AndOp ||
		f.Type == OrOp ||
		f.Type == NotOp) &&
		len(f.Inner) == 0
}

// And wraps filters using and.
func (f FilterClause) And(filters ...FilterClause) FilterClause {
	if f.None() && len(filters) == 1 {
		return filters[0]
	} else if f.Type == AndOp {
		f.Inner = append(f.Inner, filters...)
		return f
	}

	inner := append([]FilterClause{f}, filters...)
	return FilterAnd(inner...)
}

// Or wraps filters using or.
func (f FilterClause) Or(filter ...FilterClause) FilterClause {
	if f.None() && len(filter) == 1 {
		return filter[0]
	} else if f.Type == OrOp || f.None() {
		f.Type = OrOp
		f.Inner = append(f.Inner, filter...)
		return f
	}

	inner := append([]FilterClause{f}, filter...)
	return FilterOr(inner...)
}

func (f FilterClause) and(other FilterClause) FilterClause {
	if f.Type == AndOp {
		f.Inner = append(f.Inner, other)
		return f
	}

	return FilterAnd(f, other)
}

func (f FilterClause) or(other FilterClause) FilterClause {
	if f.Type == OrOp || f.None() {
		f.Type = OrOp
		f.Inner = append(f.Inner, other)
		return f
	}

	return FilterOr(f, other)
}

// AndEq append equal expression using and.
func (f FilterClause) AndEq(field string, value interface{}) FilterClause {
	return f.and(FilterEq(field, value))
}

// AndNe append not equal expression using and.
func (f FilterClause) AndNe(field string, value interface{}) FilterClause {
	return f.and(FilterNe(field, value))
}

// AndLt append lesser than expression using and.
func (f FilterClause) AndLt(field string, value interface{}) FilterClause {
	return f.and(FilterLt(field, value))
}

// AndLte append lesser than or equal expression using and.
func (f FilterClause) AndLte(field string, value interface{}) FilterClause {
	return f.and(FilterLte(field, value))
}

// AndGt append greater than expression using and.
func (f FilterClause) AndGt(field string, value interface{}) FilterClause {
	return f.and(FilterGt(field, value))
}

// AndGte append greater than or equal expression using and.
func (f FilterClause) AndGte(field string, value interface{}) FilterClause {
	return f.and(FilterGte(field, value))
}

// AndNil append is nil expression using and.
func (f FilterClause) AndNil(field string) FilterClause {
	return f.and(FilterNil(field))
}

// AndNotNil append is not nil expression using and.
func (f FilterClause) AndNotNil(field string) FilterClause {
	return f.and(FilterNotNil(field))
}

// AndIn append is in expression using and.
func (f FilterClause) AndIn(field string, values ...interface{}) FilterClause {
	return f.and(FilterIn(field, values...))
}

// AndNin append is not in expression using and.
func (f FilterClause) AndNin(field string, values ...interface{}) FilterClause {
	return f.and(FilterNin(field, values...))
}

// AndLike append like expression using and.
func (f FilterClause) AndLike(field string, pattern string) FilterClause {
	return f.and(FilterLike(field, pattern))
}

// AndNotLike append not like expression using and.
func (f FilterClause) AndNotLike(field string, pattern string) FilterClause {
	return f.and(FilterNotLike(field, pattern))
}

// AndFragment append fragment using and.
func (f FilterClause) AndFragment(expr string, values ...interface{}) FilterClause {
	return f.and(FilterFragment(expr, values...))
}

// OrEq append equal expression using or.
func (f FilterClause) OrEq(field string, value interface{}) FilterClause {
	return f.or(FilterEq(field, value))
}

// OrNe append not equal expression using or.
func (f FilterClause) OrNe(field string, value interface{}) FilterClause {
	return f.or(FilterNe(field, value))
}

// OrLt append lesser than expression using or.
func (f FilterClause) OrLt(field string, value interface{}) FilterClause {
	return f.or(FilterLt(field, value))
}

// OrLte append lesser than or equal expression using or.
func (f FilterClause) OrLte(field string, value interface{}) FilterClause {
	return f.or(FilterLte(field, value))
}

// OrGt append greater than expression using or.
func (f FilterClause) OrGt(field string, value interface{}) FilterClause {
	return f.or(FilterGt(field, value))
}

// OrGte append greater than or equal expression using or.
func (f FilterClause) OrGte(field string, value interface{}) FilterClause {
	return f.or(FilterGte(field, value))
}

// OrNil append is nil expression using or.
func (f FilterClause) OrNil(field string) FilterClause {
	return f.or(FilterNil(field))
}

// OrNotNil append is not nil expression using or.
func (f FilterClause) OrNotNil(field string) FilterClause {
	return f.or(FilterNotNil(field))
}

// OrIn append is in expression using or.
func (f FilterClause) OrIn(field string, values ...interface{}) FilterClause {
	return f.or(FilterIn(field, values...))
}

// OrNin append is not in expression using or.
func (f FilterClause) OrNin(field string, values ...interface{}) FilterClause {
	return f.or(FilterNin(field, values...))
}

// OrLike append like expression using or.
func (f FilterClause) OrLike(field string, pattern string) FilterClause {
	return f.or(FilterLike(field, pattern))
}

// OrNotLike append not like expression using or.
func (f FilterClause) OrNotLike(field string, pattern string) FilterClause {
	return f.or(FilterNotLike(field, pattern))
}

// OrFragment append fragment using or.
func (f FilterClause) OrFragment(expr string, values ...interface{}) FilterClause {
	return f.or(FilterFragment(expr, values...))
}

// FilterAnd compares other filters using and.
func FilterAnd(inner ...FilterClause) FilterClause {
	if len(inner) == 1 {
		return inner[0]
	}

	return FilterClause{
		Type:  AndOp,
		Inner: inner,
	}
}

// FilterOr compares other filters using and.
func FilterOr(inner ...FilterClause) FilterClause {
	if len(inner) == 1 {
		return inner[0]
	}

	return FilterClause{
		Type:  OrOp,
		Inner: inner,
	}
}

// FilterNot wraps filters using not.
// It'll negate the filter type if possible.
func FilterNot(inner ...FilterClause) FilterClause {
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
			return FilterClause{
				Type:  NotOp,
				Inner: inner,
			}
		}

		return f
	}

	return FilterClause{
		Type:  NotOp,
		Inner: inner,
	}
}

// FilterEq expression field equal to value.
func FilterEq(field string, value interface{}) FilterClause {
	return FilterClause{
		Type:   EqOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterNe compares that left value is not equal to right value.
func FilterNe(field string, value interface{}) FilterClause {
	return FilterClause{
		Type:   NeOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterLt compares that left value is less than to right value.
func FilterLt(field string, value interface{}) FilterClause {
	return FilterClause{
		Type:   LtOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterLte compares that left value is less than or equal to right value.
func FilterLte(field string, value interface{}) FilterClause {
	return FilterClause{
		Type:   LteOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterGt compares that left value is greater than to right value.
func FilterGt(field string, value interface{}) FilterClause {
	return FilterClause{
		Type:   GtOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterGte compares that left value is greater than or equal to right value.
func FilterGte(field string, value interface{}) FilterClause {
	return FilterClause{
		Type:   GteOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterNil check whether field is nil.
func FilterNil(field string) FilterClause {
	return FilterClause{
		Type:  NilOp,
		Field: field,
	}
}

// FilterNotNil check whether field is not nil.
func FilterNotNil(field string) FilterClause {
	return FilterClause{
		Type:  NotNilOp,
		Field: field,
	}
}

// FilterIn check whethers value of the field is included in values.
func FilterIn(field string, values ...interface{}) FilterClause {
	return FilterClause{
		Type:   InOp,
		Field:  field,
		Values: values,
	}
}

// FilterNin check whethers value of the field is not included in values.
func FilterNin(field string, values ...interface{}) FilterClause {
	return FilterClause{
		Type:   NinOp,
		Field:  field,
		Values: values,
	}
}

// FilterLike compares value of field to match string pattern.
func FilterLike(field string, pattern string) FilterClause {
	return FilterClause{
		Type:   LikeOp,
		Field:  field,
		Values: []interface{}{pattern},
	}
}

// FilterNotLike compares value of field to not match string pattern.
func FilterNotLike(field string, pattern string) FilterClause {
	return FilterClause{
		Type:   NotLikeOp,
		Field:  field,
		Values: []interface{}{pattern},
	}
}

// FilterFragment add custom filter.
func FilterFragment(expr string, values ...interface{}) FilterClause {
	return FilterClause{
		Type:   FragmentOp,
		Field:  expr,
		Values: values,
	}
}

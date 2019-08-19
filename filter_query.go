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

// FilterQuery defines details of a coundition type.
type FilterQuery struct {
	Type   FilterOp
	Field  string
	Values []interface{} // TODO: don't use array (let adapter cast when needed)
	Inner  []FilterQuery
}

// Build Filter query.
func (fq FilterQuery) Build(query *Query) {
	query.WhereClause = query.WhereClause.And(fq)
}

// None returns true if no filter is specified.
func (fq FilterQuery) None() bool {
	return (fq.Type == AndOp ||
		fq.Type == OrOp ||
		fq.Type == NotOp) &&
		len(fq.Inner) == 0
}

// And wraps filters using and.
func (fq FilterQuery) And(filters ...FilterQuery) FilterQuery {
	if fq.None() && len(filters) == 1 {
		return filters[0]
	} else if fq.Type == AndOp {
		fq.Inner = append(fq.Inner, filters...)
		return fq
	}

	inner := append([]FilterQuery{fq}, filters...)
	return FilterAnd(inner...)
}

// Or wraps filters using or.
func (fq FilterQuery) Or(filter ...FilterQuery) FilterQuery {
	if fq.None() && len(filter) == 1 {
		return filter[0]
	} else if fq.Type == OrOp || fq.None() {
		fq.Type = OrOp
		fq.Inner = append(fq.Inner, filter...)
		return fq
	}

	inner := append([]FilterQuery{fq}, filter...)
	return FilterOr(inner...)
}

func (fq FilterQuery) and(other FilterQuery) FilterQuery {
	if fq.Type == AndOp {
		fq.Inner = append(fq.Inner, other)
		return fq
	}

	return FilterAnd(fq, other)
}

func (fq FilterQuery) or(other FilterQuery) FilterQuery {
	if fq.Type == OrOp || fq.None() {
		fq.Type = OrOp
		fq.Inner = append(fq.Inner, other)
		return fq
	}

	return FilterOr(fq, other)
}

// AndEq append equal expression using and.
func (fq FilterQuery) AndEq(field string, value interface{}) FilterQuery {
	return fq.and(FilterEq(field, value))
}

// AndNe append not equal expression using and.
func (fq FilterQuery) AndNe(field string, value interface{}) FilterQuery {
	return fq.and(FilterNe(field, value))
}

// AndLt append lesser than expression using and.
func (fq FilterQuery) AndLt(field string, value interface{}) FilterQuery {
	return fq.and(FilterLt(field, value))
}

// AndLte append lesser than or equal expression using and.
func (fq FilterQuery) AndLte(field string, value interface{}) FilterQuery {
	return fq.and(FilterLte(field, value))
}

// AndGt append greater than expression using and.
func (fq FilterQuery) AndGt(field string, value interface{}) FilterQuery {
	return fq.and(FilterGt(field, value))
}

// AndGte append greater than or equal expression using and.
func (fq FilterQuery) AndGte(field string, value interface{}) FilterQuery {
	return fq.and(FilterGte(field, value))
}

// AndNil append is nil expression using and.
func (fq FilterQuery) AndNil(field string) FilterQuery {
	return fq.and(FilterNil(field))
}

// AndNotNil append is not nil expression using and.
func (fq FilterQuery) AndNotNil(field string) FilterQuery {
	return fq.and(FilterNotNil(field))
}

// AndIn append is in expression using and.
func (fq FilterQuery) AndIn(field string, values ...interface{}) FilterQuery {
	return fq.and(FilterIn(field, values...))
}

// AndNin append is not in expression using and.
func (fq FilterQuery) AndNin(field string, values ...interface{}) FilterQuery {
	return fq.and(FilterNin(field, values...))
}

// AndLike append like expression using and.
func (fq FilterQuery) AndLike(field string, pattern string) FilterQuery {
	return fq.and(FilterLike(field, pattern))
}

// AndNotLike append not like expression using and.
func (fq FilterQuery) AndNotLike(field string, pattern string) FilterQuery {
	return fq.and(FilterNotLike(field, pattern))
}

// AndFragment append fragment using and.
func (fq FilterQuery) AndFragment(expr string, values ...interface{}) FilterQuery {
	return fq.and(FilterFragment(expr, values...))
}

// OrEq append equal expression using or.
func (fq FilterQuery) OrEq(field string, value interface{}) FilterQuery {
	return fq.or(FilterEq(field, value))
}

// OrNe append not equal expression using or.
func (fq FilterQuery) OrNe(field string, value interface{}) FilterQuery {
	return fq.or(FilterNe(field, value))
}

// OrLt append lesser than expression using or.
func (fq FilterQuery) OrLt(field string, value interface{}) FilterQuery {
	return fq.or(FilterLt(field, value))
}

// OrLte append lesser than or equal expression using or.
func (fq FilterQuery) OrLte(field string, value interface{}) FilterQuery {
	return fq.or(FilterLte(field, value))
}

// OrGt append greater than expression using or.
func (fq FilterQuery) OrGt(field string, value interface{}) FilterQuery {
	return fq.or(FilterGt(field, value))
}

// OrGte append greater than or equal expression using or.
func (fq FilterQuery) OrGte(field string, value interface{}) FilterQuery {
	return fq.or(FilterGte(field, value))
}

// OrNil append is nil expression using or.
func (fq FilterQuery) OrNil(field string) FilterQuery {
	return fq.or(FilterNil(field))
}

// OrNotNil append is not nil expression using or.
func (fq FilterQuery) OrNotNil(field string) FilterQuery {
	return fq.or(FilterNotNil(field))
}

// OrIn append is in expression using or.
func (fq FilterQuery) OrIn(field string, values ...interface{}) FilterQuery {
	return fq.or(FilterIn(field, values...))
}

// OrNin append is not in expression using or.
func (fq FilterQuery) OrNin(field string, values ...interface{}) FilterQuery {
	return fq.or(FilterNin(field, values...))
}

// OrLike append like expression using or.
func (fq FilterQuery) OrLike(field string, pattern string) FilterQuery {
	return fq.or(FilterLike(field, pattern))
}

// OrNotLike append not like expression using or.
func (fq FilterQuery) OrNotLike(field string, pattern string) FilterQuery {
	return fq.or(FilterNotLike(field, pattern))
}

// OrFragment append fragment using or.
func (fq FilterQuery) OrFragment(expr string, values ...interface{}) FilterQuery {
	return fq.or(FilterFragment(expr, values...))
}

// FilterAnd compares other filters using and.
func FilterAnd(inner ...FilterQuery) FilterQuery {
	if len(inner) == 1 {
		return inner[0]
	}

	return FilterQuery{
		Type:  AndOp,
		Inner: inner,
	}
}

// FilterOr compares other filters using and.
func FilterOr(inner ...FilterQuery) FilterQuery {
	if len(inner) == 1 {
		return inner[0]
	}

	return FilterQuery{
		Type:  OrOp,
		Inner: inner,
	}
}

// FilterNot wraps filters using not.
// It'll negate the filter type if possible.
func FilterNot(inner ...FilterQuery) FilterQuery {
	if len(inner) == 1 {
		fq := inner[0]
		switch fq.Type {
		case EqOp:
			fq.Type = NeOp
			return fq
		case LtOp:
			fq.Type = GteOp
		case LteOp:
			fq.Type = GtOp
		case GtOp:
			fq.Type = LteOp
		case GteOp:
			fq.Type = LtOp
		case NilOp:
			fq.Type = NotNilOp
		case InOp:
			fq.Type = NinOp
		case LikeOp:
			fq.Type = NotLikeOp
		default:
			return FilterQuery{
				Type:  NotOp,
				Inner: inner,
			}
		}

		return fq
	}

	return FilterQuery{
		Type:  NotOp,
		Inner: inner,
	}
}

// FilterEq expression field equal to value.
func FilterEq(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:   EqOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterNe compares that left value is not equal to right value.
func FilterNe(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:   NeOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterLt compares that left value is less than to right value.
func FilterLt(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:   LtOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterLte compares that left value is less than or equal to right value.
func FilterLte(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:   LteOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterGt compares that left value is greater than to right value.
func FilterGt(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:   GtOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterGte compares that left value is greater than or equal to right value.
func FilterGte(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:   GteOp,
		Field:  field,
		Values: []interface{}{value},
	}
}

// FilterNil check whether field is nil.
func FilterNil(field string) FilterQuery {
	return FilterQuery{
		Type:  NilOp,
		Field: field,
	}
}

// FilterNotNil check whether field is not nil.
func FilterNotNil(field string) FilterQuery {
	return FilterQuery{
		Type:  NotNilOp,
		Field: field,
	}
}

// FilterIn check whethers value of the field is included in values.
func FilterIn(field string, values ...interface{}) FilterQuery {
	return FilterQuery{
		Type:   InOp,
		Field:  field,
		Values: values,
	}
}

// FilterNin check whethers value of the field is not included in values.
func FilterNin(field string, values ...interface{}) FilterQuery {
	return FilterQuery{
		Type:   NinOp,
		Field:  field,
		Values: values,
	}
}

// FilterLike compares value of field to match string pattern.
func FilterLike(field string, pattern string) FilterQuery {
	return FilterQuery{
		Type:   LikeOp,
		Field:  field,
		Values: []interface{}{pattern},
	}
}

// FilterNotLike compares value of field to not match string pattern.
func FilterNotLike(field string, pattern string) FilterQuery {
	return FilterQuery{
		Type:   NotLikeOp,
		Field:  field,
		Values: []interface{}{pattern},
	}
}

// FilterFragment add custom filter.
func FilterFragment(expr string, values ...interface{}) FilterQuery {
	return FilterQuery{
		Type:   FragmentOp,
		Field:  expr,
		Values: values,
	}
}

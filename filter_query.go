package rel

// FilterOp defines enumeration of all supported filter types.
type FilterOp int

const (
	// FilterAndOp is filter type for and operator.
	FilterAndOp FilterOp = iota
	// FilterOrOp is filter type for or operator.
	FilterOrOp
	// FilterNotOp is filter type for not operator.
	FilterNotOp

	// FilterEqOp is filter type for equal comparison.
	FilterEqOp
	// FilterNeOp is filter type for not equal comparison.
	FilterNeOp

	// FilterLtOp is filter type for less than comparison.
	FilterLtOp
	// FilterLteOp is filter type for less than or equal comparison.
	FilterLteOp
	// FilterGtOp is filter type for greater than comparison.
	FilterGtOp
	// FilterGteOp is filter type for greter than or equal comparison.
	FilterGteOp

	// FilterNilOp is filter type for nil check.
	FilterNilOp
	// FilterNotNilOp is filter type for not nil check.
	FilterNotNilOp

	// FilterInOp is filter type for inclusion comparison.
	FilterInOp
	// FilterNinOp is filter type for not inclusion comparison.
	FilterNinOp

	// FilterLikeOp is filter type for like comparison.
	FilterLikeOp
	// FilterNotLikeOp is filter type for not like comparison.
	FilterNotLikeOp

	// FilterFragmentOp is filter type for custom filter.
	FilterFragmentOp
)

// FilterQuery defines details of a coundition type.
type FilterQuery struct {
	Type  FilterOp
	Field string
	Value interface{}
	Inner []FilterQuery
}

// Build Filter query.
func (fq FilterQuery) Build(query *Query) {
	query.WhereQuery = query.WhereQuery.And(fq)
}

// None returns true if no filter is specified.
func (fq FilterQuery) None() bool {
	return (fq.Type == FilterAndOp ||
		fq.Type == FilterOrOp ||
		fq.Type == FilterNotOp) &&
		len(fq.Inner) == 0
}

// And wraps filters using and.
func (fq FilterQuery) And(filters ...FilterQuery) FilterQuery {
	if fq.None() && len(filters) == 1 {
		return filters[0]
	} else if fq.Type == FilterAndOp {
		fq.Inner = append(fq.Inner, filters...)
		return fq
	}

	inner := append([]FilterQuery{fq}, filters...)
	return And(inner...)
}

// Or wraps filters using or.
func (fq FilterQuery) Or(filter ...FilterQuery) FilterQuery {
	if fq.None() && len(filter) == 1 {
		return filter[0]
	} else if fq.Type == FilterOrOp || fq.None() {
		fq.Type = FilterOrOp
		fq.Inner = append(fq.Inner, filter...)
		return fq
	}

	inner := append([]FilterQuery{fq}, filter...)
	return Or(inner...)
}

func (fq FilterQuery) and(other FilterQuery) FilterQuery {
	if fq.Type == FilterAndOp {
		fq.Inner = append(fq.Inner, other)
		return fq
	}

	return And(fq, other)
}

func (fq FilterQuery) or(other FilterQuery) FilterQuery {
	if fq.Type == FilterOrOp || fq.None() {
		fq.Type = FilterOrOp
		fq.Inner = append(fq.Inner, other)
		return fq
	}

	return Or(fq, other)
}

// AndEq append equal expression using and.
func (fq FilterQuery) AndEq(field string, value interface{}) FilterQuery {
	return fq.and(Eq(field, value))
}

// AndNe append not equal expression using and.
func (fq FilterQuery) AndNe(field string, value interface{}) FilterQuery {
	return fq.and(Ne(field, value))
}

// AndLt append lesser than expression using and.
func (fq FilterQuery) AndLt(field string, value interface{}) FilterQuery {
	return fq.and(Lt(field, value))
}

// AndLte append lesser than or equal expression using and.
func (fq FilterQuery) AndLte(field string, value interface{}) FilterQuery {
	return fq.and(Lte(field, value))
}

// AndGt append greater than expression using and.
func (fq FilterQuery) AndGt(field string, value interface{}) FilterQuery {
	return fq.and(Gt(field, value))
}

// AndGte append greater than or equal expression using and.
func (fq FilterQuery) AndGte(field string, value interface{}) FilterQuery {
	return fq.and(Gte(field, value))
}

// AndNil append is nil expression using and.
func (fq FilterQuery) AndNil(field string) FilterQuery {
	return fq.and(Nil(field))
}

// AndNotNil append is not nil expression using and.
func (fq FilterQuery) AndNotNil(field string) FilterQuery {
	return fq.and(NotNil(field))
}

// AndIn append is in expression using and.
func (fq FilterQuery) AndIn(field string, values ...interface{}) FilterQuery {
	return fq.and(In(field, values...))
}

// AndNin append is not in expression using and.
func (fq FilterQuery) AndNin(field string, values ...interface{}) FilterQuery {
	return fq.and(Nin(field, values...))
}

// AndLike append like expression using and.
func (fq FilterQuery) AndLike(field string, pattern string) FilterQuery {
	return fq.and(Like(field, pattern))
}

// AndNotLike append not like expression using and.
func (fq FilterQuery) AndNotLike(field string, pattern string) FilterQuery {
	return fq.and(NotLike(field, pattern))
}

// AndFragment append fragment using and.
func (fq FilterQuery) AndFragment(expr string, values ...interface{}) FilterQuery {
	return fq.and(FilterFragment(expr, values...))
}

// OrEq append equal expression using or.
func (fq FilterQuery) OrEq(field string, value interface{}) FilterQuery {
	return fq.or(Eq(field, value))
}

// OrNe append not equal expression using or.
func (fq FilterQuery) OrNe(field string, value interface{}) FilterQuery {
	return fq.or(Ne(field, value))
}

// OrLt append lesser than expression using or.
func (fq FilterQuery) OrLt(field string, value interface{}) FilterQuery {
	return fq.or(Lt(field, value))
}

// OrLte append lesser than or equal expression using or.
func (fq FilterQuery) OrLte(field string, value interface{}) FilterQuery {
	return fq.or(Lte(field, value))
}

// OrGt append greater than expression using or.
func (fq FilterQuery) OrGt(field string, value interface{}) FilterQuery {
	return fq.or(Gt(field, value))
}

// OrGte append greater than or equal expression using or.
func (fq FilterQuery) OrGte(field string, value interface{}) FilterQuery {
	return fq.or(Gte(field, value))
}

// OrNil append is nil expression using or.
func (fq FilterQuery) OrNil(field string) FilterQuery {
	return fq.or(Nil(field))
}

// OrNotNil append is not nil expression using or.
func (fq FilterQuery) OrNotNil(field string) FilterQuery {
	return fq.or(NotNil(field))
}

// OrIn append is in expression using or.
func (fq FilterQuery) OrIn(field string, values ...interface{}) FilterQuery {
	return fq.or(In(field, values...))
}

// OrNin append is not in expression using or.
func (fq FilterQuery) OrNin(field string, values ...interface{}) FilterQuery {
	return fq.or(Nin(field, values...))
}

// OrLike append like expression using or.
func (fq FilterQuery) OrLike(field string, pattern string) FilterQuery {
	return fq.or(Like(field, pattern))
}

// OrNotLike append not like expression using or.
func (fq FilterQuery) OrNotLike(field string, pattern string) FilterQuery {
	return fq.or(NotLike(field, pattern))
}

// OrFragment append fragment using or.
func (fq FilterQuery) OrFragment(expr string, values ...interface{}) FilterQuery {
	return fq.or(FilterFragment(expr, values...))
}

// And compares other filters using and.
func And(inner ...FilterQuery) FilterQuery {
	if len(inner) == 1 {
		return inner[0]
	}

	return FilterQuery{
		Type:  FilterAndOp,
		Inner: inner,
	}
}

// Or compares other filters using and.
func Or(inner ...FilterQuery) FilterQuery {
	if len(inner) == 1 {
		return inner[0]
	}

	return FilterQuery{
		Type:  FilterOrOp,
		Inner: inner,
	}
}

// Not wraps filters using not.
// It'll negate the filter type if possible.
func Not(inner ...FilterQuery) FilterQuery {
	if len(inner) == 1 {
		fq := inner[0]
		switch fq.Type {
		case FilterEqOp:
			fq.Type = FilterNeOp
			return fq
		case FilterLtOp:
			fq.Type = FilterGteOp
		case FilterLteOp:
			fq.Type = FilterGtOp
		case FilterGtOp:
			fq.Type = FilterLteOp
		case FilterGteOp:
			fq.Type = FilterLtOp
		case FilterNilOp:
			fq.Type = FilterNotNilOp
		case FilterInOp:
			fq.Type = FilterNinOp
		case FilterLikeOp:
			fq.Type = FilterNotLikeOp
		default:
			return FilterQuery{
				Type:  FilterNotOp,
				Inner: inner,
			}
		}

		return fq
	}

	return FilterQuery{
		Type:  FilterNotOp,
		Inner: inner,
	}
}

// Eq expression field equal to value.
func Eq(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterEqOp,
		Field: field,
		Value: value,
	}
}

// Ne compares that left value is not equal to right value.
func Ne(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterNeOp,
		Field: field,
		Value: value,
	}
}

// Lt compares that left value is less than to right value.
func Lt(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterLtOp,
		Field: field,
		Value: value,
	}
}

// Lte compares that left value is less than or equal to right value.
func Lte(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterLteOp,
		Field: field,
		Value: value,
	}
}

// Gt compares that left value is greater than to right value.
func Gt(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterGtOp,
		Field: field,
		Value: value,
	}
}

// Gte compares that left value is greater than or equal to right value.
func Gte(field string, value interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterGteOp,
		Field: field,
		Value: value,
	}
}

// Nil check whether field is nil.
func Nil(field string) FilterQuery {
	return FilterQuery{
		Type:  FilterNilOp,
		Field: field,
	}
}

// NotNil check whether field is not nil.
func NotNil(field string) FilterQuery {
	return FilterQuery{
		Type:  FilterNotNilOp,
		Field: field,
	}
}

// In check whethers value of the field is included in values.
func In(field string, values ...interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterInOp,
		Field: field,
		Value: values,
	}
}

// InInt check whethers integer values of the field is included.
func InInt(field string, values []int) FilterQuery {
	var (
		ivalues = make([]interface{}, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return In(field, ivalues...)
}

// InUint check whethers unsigned integer values of the field is included.
func InUint(field string, values []uint) FilterQuery {
	var (
		ivalues = make([]interface{}, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return In(field, ivalues...)
}

// InString check whethers string values of the field is included.
func InString(field string, values []string) FilterQuery {
	var (
		ivalues = make([]interface{}, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return In(field, ivalues...)
}

// Nin check whethers value of the field is not included in values.
func Nin(field string, values ...interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterNinOp,
		Field: field,
		Value: values,
	}
}

// NinInt check whethers integer values of the is not included.
func NinInt(field string, values []int) FilterQuery {
	var (
		ivalues = make([]interface{}, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return Nin(field, ivalues...)
}

// NinUint check whethers unsigned integer values of the is not included.
func NinUint(field string, values []uint) FilterQuery {
	var (
		ivalues = make([]interface{}, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return Nin(field, ivalues...)
}

// NinString check whethers string values of the is not included.
func NinString(field string, values []string) FilterQuery {
	var (
		ivalues = make([]interface{}, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return Nin(field, ivalues...)
}

// Like compares value of field to match string pattern.
func Like(field string, pattern string) FilterQuery {
	return FilterQuery{
		Type:  FilterLikeOp,
		Field: field,
		Value: pattern,
	}
}

// NotLike compares value of field to not match string pattern.
func NotLike(field string, pattern string) FilterQuery {
	return FilterQuery{
		Type:  FilterNotLikeOp,
		Field: field,
		Value: pattern,
	}
}

// FilterFragment add custom filter.
func FilterFragment(expr string, values ...interface{}) FilterQuery {
	return FilterQuery{
		Type:  FilterFragmentOp,
		Field: expr,
		Value: values,
	}
}

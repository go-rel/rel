package rel

import (
	"errors"
	"strings"
)

// FilterOp defines enumeration of all supported filter types.
type FilterOp int

func (fo FilterOp) String() string {
	return [...]string{
		"And",
		"Or",
		"Not",
		"Eq",
		"Ne",
		"Lt",
		"Lte",
		"Gt",
		"Gte",
		"Nil",
		"NotNil",
		"In",
		"Nin",
		"Like",
		"NotLike",
		"Fragment",
	}[fo]
}

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

// FilterQuery defines details of a condition type.
type FilterQuery struct {
	Type  FilterOp
	Field string
	Value any
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

func (fq FilterQuery) String() string {
	if fq.None() {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("where.")
	builder.WriteString(fq.Type.String())
	builder.WriteByte('(')

	switch fq.Type {
	case FilterAndOp, FilterOrOp, FilterNotOp:
		for i := range fq.Inner {
			if i > 0 {
				builder.WriteString(", ")
			}

			builder.WriteString(fq.Inner[i].String())
		}
	case FilterEqOp, FilterNeOp, FilterLtOp, FilterLteOp, FilterGtOp, FilterGteOp:
		builder.WriteByte('"')
		builder.WriteString(fq.Field)
		builder.WriteString("\", ")
		builder.WriteString(fmtAny(fq.Value))
	case FilterNilOp, FilterNotNilOp, FilterLikeOp, FilterNotLikeOp:
		builder.WriteByte('"')
		builder.WriteString(fq.Field)
		builder.WriteByte('"')
	case FilterInOp, FilterNinOp:
		builder.WriteByte('"')
		builder.WriteString(fq.Field)
		builder.WriteString("\", ")
		builder.WriteString(fmtAnys(fq.Value.([]any)))
	case FilterFragmentOp:
		v := fq.Value.([]any)
		builder.WriteByte('"')
		builder.WriteString(fq.Field)
		builder.WriteByte('"')

		if len(v) > 0 {
			builder.WriteString(", ")
			builder.WriteString(fmtAnys(v))
		}
	}

	builder.WriteByte(')')

	return builder.String()
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

func (fq FilterQuery) applyIndex(index *Index) {
	index.Filter = fq
}

// AndEq append equal expression using and.
func (fq FilterQuery) AndEq(field string, value any) FilterQuery {
	return fq.and(Eq(field, value))
}

// AndNe append not equal expression using and.
func (fq FilterQuery) AndNe(field string, value any) FilterQuery {
	return fq.and(Ne(field, value))
}

// AndLt append lesser than expression using and.
func (fq FilterQuery) AndLt(field string, value any) FilterQuery {
	return fq.and(Lt(field, value))
}

// AndLte append lesser than or equal expression using and.
func (fq FilterQuery) AndLte(field string, value any) FilterQuery {
	return fq.and(Lte(field, value))
}

// AndGt append greater than expression using and.
func (fq FilterQuery) AndGt(field string, value any) FilterQuery {
	return fq.and(Gt(field, value))
}

// AndGte append greater than or equal expression using and.
func (fq FilterQuery) AndGte(field string, value any) FilterQuery {
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
func (fq FilterQuery) AndIn(field string, values ...any) FilterQuery {
	return fq.and(In(field, values...))
}

// AndNin append is not in expression using and.
func (fq FilterQuery) AndNin(field string, values ...any) FilterQuery {
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
func (fq FilterQuery) AndFragment(expr string, values ...any) FilterQuery {
	return fq.and(FilterFragment(expr, values...))
}

// OrEq append equal expression using or.
func (fq FilterQuery) OrEq(field string, value any) FilterQuery {
	return fq.or(Eq(field, value))
}

// OrNe append not equal expression using or.
func (fq FilterQuery) OrNe(field string, value any) FilterQuery {
	return fq.or(Ne(field, value))
}

// OrLt append lesser than expression using or.
func (fq FilterQuery) OrLt(field string, value any) FilterQuery {
	return fq.or(Lt(field, value))
}

// OrLte append lesser than or equal expression using or.
func (fq FilterQuery) OrLte(field string, value any) FilterQuery {
	return fq.or(Lte(field, value))
}

// OrGt append greater than expression using or.
func (fq FilterQuery) OrGt(field string, value any) FilterQuery {
	return fq.or(Gt(field, value))
}

// OrGte append greater than or equal expression using or.
func (fq FilterQuery) OrGte(field string, value any) FilterQuery {
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
func (fq FilterQuery) OrIn(field string, values ...any) FilterQuery {
	return fq.or(In(field, values...))
}

// OrNin append is not in expression using or.
func (fq FilterQuery) OrNin(field string, values ...any) FilterQuery {
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
func (fq FilterQuery) OrFragment(expr string, values ...any) FilterQuery {
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

// Or compares other filters using or.
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
func Eq(field string, value any) FilterQuery {
	return FilterQuery{
		Type:  FilterEqOp,
		Field: field,
		Value: value,
	}
}

func lockVersion(version int) FilterQuery {
	return Eq("lock_version", version)
}

// Ne compares that left value is not equal to right value.
func Ne(field string, value any) FilterQuery {
	return FilterQuery{
		Type:  FilterNeOp,
		Field: field,
		Value: value,
	}
}

// Lt compares that left value is less than to right value.
func Lt(field string, value any) FilterQuery {
	return FilterQuery{
		Type:  FilterLtOp,
		Field: field,
		Value: value,
	}
}

// Lte compares that left value is less than or equal to right value.
func Lte(field string, value any) FilterQuery {
	return FilterQuery{
		Type:  FilterLteOp,
		Field: field,
		Value: value,
	}
}

// Gt compares that left value is greater than to right value.
func Gt(field string, value any) FilterQuery {
	return FilterQuery{
		Type:  FilterGtOp,
		Field: field,
		Value: value,
	}
}

// Gte compares that left value is greater than or equal to right value.
func Gte(field string, value any) FilterQuery {
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
func In(field string, values ...any) FilterQuery {
	return FilterQuery{
		Type:  FilterInOp,
		Field: field,
		Value: values,
	}
}

// InInt check whethers integer values of the field is included.
func InInt(field string, values []int) FilterQuery {
	var (
		ivalues = make([]any, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return In(field, ivalues...)
}

// InUint check whethers unsigned integer values of the field is included.
func InUint(field string, values []uint) FilterQuery {
	var (
		ivalues = make([]any, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return In(field, ivalues...)
}

// InString check whethers string values of the field is included.
func InString(field string, values []string) FilterQuery {
	var (
		ivalues = make([]any, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return In(field, ivalues...)
}

// Nin check whethers value of the field is not included in values.
func Nin(field string, values ...any) FilterQuery {
	return FilterQuery{
		Type:  FilterNinOp,
		Field: field,
		Value: values,
	}
}

// NinInt check whethers integer values of the is not included.
func NinInt(field string, values []int) FilterQuery {
	var (
		ivalues = make([]any, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return Nin(field, ivalues...)
}

// NinUint check whethers unsigned integer values of the is not included.
func NinUint(field string, values []uint) FilterQuery {
	var (
		ivalues = make([]any, len(values))
	)

	for i := range values {
		ivalues[i] = values[i]
	}

	return Nin(field, ivalues...)
}

// NinString check whethers string values of the is not included.
func NinString(field string, values []string) FilterQuery {
	var (
		ivalues = make([]any, len(values))
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
func FilterFragment(expr string, values ...any) FilterQuery {
	return FilterQuery{
		Type:  FilterFragmentOp,
		Field: expr,
		Value: values,
	}
}

func filterDocument(doc *Document) FilterQuery {
	var (
		pFields = doc.PrimaryFields()
		pValues = doc.PrimaryValues()
	)

	return filterDocumentPrimary(pFields, pValues, FilterEqOp)
}

func filterDocumentPrimary(pFields []string, pValues []any, op FilterOp) FilterQuery {
	var filter FilterQuery

	for i := range pFields {
		filter = filter.And(FilterQuery{
			Type:  op,
			Field: pFields[i],
			Value: pValues[i],
		})
	}

	return filter

}

func filterCollection(col *Collection) FilterQuery {
	var (
		pFields = col.PrimaryFields()
		pValues = col.PrimaryValues()
		length  = col.Len()
	)

	return filterCollectionPrimary(pFields, pValues, length)
}

func filterCollectionPrimary(pFields []string, pValues []any, length int) FilterQuery {
	var filter FilterQuery

	if len(pFields) == 1 {
		filter = In(pFields[0], pValues[0].([]any)...)
	} else {
		var (
			andFilters = make([]FilterQuery, length)
		)

		for i := range pValues {
			var (
				values = pValues[i].([]any)
			)

			for j := range values {
				andFilters[j] = andFilters[j].AndEq(pFields[i], values[j])
			}
		}

		filter = Or(andFilters...)
	}

	return filter
}

func filterBelongsTo(assoc Association) (FilterQuery, error) {
	var (
		rValue = assoc.ReferenceValue()
		fValue = assoc.ForeignValue()
		filter = Eq(assoc.ForeignField(), fValue)
	)

	if rValue != fValue {
		return filter, ConstraintError{
			Key:  assoc.ReferenceField(),
			Type: ForeignKeyConstraint,
			Err:  errors.New("rel: inconsistent belongs to ref and fk"),
		}
	}

	return filter, nil
}

func filterHasOne(assoc Association, asssocDoc *Document) (FilterQuery, error) {
	var (
		fField = assoc.ForeignField()
		fValue = assoc.ForeignValue()
		rValue = assoc.ReferenceValue()
		filter = filterDocument(asssocDoc).AndEq(fField, rValue)
	)

	if rValue != fValue {
		return filter, ConstraintError{
			Key:  fField,
			Type: ForeignKeyConstraint,
			Err:  errors.New("rel: inconsistent has one ref and fk"),
		}
	}

	return filter, nil
}

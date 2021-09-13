package rel

// SelectQuery defines select clause of the query.
type SelectQuery struct {
	OnlyDistinct bool
	Fields       []string
}

// Distinct select query.
func (sq SelectQuery) Distinct() SelectQuery {
	sq.OnlyDistinct = true
	return sq
}

// NewSelect query.
//
// Deprecated: use Select instead
func NewSelect(fields ...string) SelectQuery {
	return SelectQuery{
		Fields: fields,
	}
}

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
func NewSelect(fields ...string) SelectQuery {
	return SelectQuery{
		Fields: fields,
	}
}

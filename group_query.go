package rel

// GroupQuery defines group clause of the query.
type GroupQuery struct {
	Fields []string
	Filter FilterQuery
}

// Build query.
func (gq GroupQuery) Build(query *Query) {
	query.GroupQuery = gq
}

// Having appends filter for group query with and operand.
func (gq GroupQuery) Having(filters ...FilterQuery) GroupQuery {
	gq.Filter = gq.Filter.And(filters...)
	return gq
}

// OrHaving appends filter for group query with or operand.
func (gq GroupQuery) OrHaving(filters ...FilterQuery) GroupQuery {
	gq.Filter = gq.Filter.Or(And(filters...))
	return gq
}

// Where is alias for having.
func (gq GroupQuery) Where(filters ...FilterQuery) GroupQuery {
	return gq.Having(filters...)
}

// OrWhere is alias for OrHaving.
func (gq GroupQuery) OrWhere(filters ...FilterQuery) GroupQuery {
	return gq.OrHaving(filters...)
}

// NewGroup query.
func NewGroup(fields ...string) GroupQuery {
	return GroupQuery{
		Fields: fields,
	}
}

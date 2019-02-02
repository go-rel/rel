package query

type Query struct {
	built        bool
	OffsetResult Offset
	LimitResult  Limit
}

// Offset the result returned by database.
func (query Query) Offset(offset Offset) Query {
	query.OffsetResult = offset
	return query
}

// Limit result returned by database.
func (query Query) Limit(limit Limit) Query {
	query.LimitResult = limit
	return query
}

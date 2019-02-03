package query

type Query struct {
	built        bool
	WhereClause  Filter
	SortClause   []Sort
	OffsetClause Offset
	LimitClause  Limit
}

func (query Query) Where(filters ...Filter) Query {
	query.WhereClause = FilterAnd(filters...)
	return query
}

func (query Query) Sort(sorts ...Sort) Query {
	query.SortClause = sorts
	return query
}

// Offset the result returned by database.
func (query Query) Offset(offset Offset) Query {
	query.OffsetClause = offset
	return query
}

// Limit result returned by database.
func (query Query) Limit(limit Limit) Query {
	query.LimitClause = limit
	return query
}

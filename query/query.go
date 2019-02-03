package query

type Query struct {
	built        bool
	WhereClause  Filter
	GroupClause  Group
	SortClause   []Sort
	OffsetClause Offset
	LimitClause  Limit
}

func (query Query) Where(filters ...Filter) Query {
	query.WhereClause = FilterAnd(filters...)
	return query
}

func (query Query) Group(fields ...string) Query {
	query.GroupClause.Fields = fields
	return query
}

func (query Query) Having(filters ...Filter) Query {
	query.GroupClause.Filters = filters
	return query
}

func (query Query) Sort(fields ...string) Query {
	return query.SortAsc(fields...)
}

func (query Query) SortAsc(fields ...string) Query {
	sorts := make([]Sort, len(fields))
	for i := range fields {
		sorts[i] = SortAsc(fields[i])
	}

	query.SortClause = append(query.SortClause, sorts...)
	return query
}

func (query Query) SortDesc(fields ...string) Query {
	sorts := make([]Sort, len(fields))
	for i := range fields {
		sorts[i] = SortDesc(fields[i])
	}

	query.SortClause = append(query.SortClause, sorts...)
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

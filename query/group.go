package query

type GroupClause struct {
	Fields  []string
	Filters []FilterClause
}

func (g GroupClause) Build(query *Query) {
	query.GroupClause = g
}

func (g GroupClause) Having(filters ...FilterClause) GroupClause {
	g.Filters = filters
	return g
}

func (g GroupClause) Where(filters ...FilterClause) GroupClause {
	return g.Having(filters...)
}

func Group(fields ...string) GroupClause {
	return GroupClause{
		Fields: fields,
	}
}

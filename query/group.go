package query

type GroupClause struct {
	Fields []string
	Filter FilterClause
}

func (g GroupClause) Build(query *Query) {
	query.GroupClause = g
}

func (g GroupClause) Having(filters ...FilterClause) GroupClause {
	g.Filter = g.Filter.And(filters...)
	return g
}

func (g GroupClause) OrHaving(filters ...FilterClause) GroupClause {
	g.Filter = g.Filter.Or(FilterAnd(filters...))
	return g
}

func (g GroupClause) Where(filters ...FilterClause) GroupClause {
	return g.Having(filters...)
}

func (g GroupClause) OrWhere(filters ...FilterClause) GroupClause {
	return g.OrHaving(filters...)
}

func NewGroup(fields ...string) GroupClause {
	return GroupClause{
		Fields: fields,
	}
}

package grimoire

type GroupQuery struct {
	Fields []string
	Filter FilterQuery
}

func (gq GroupQuery) Build(query *Query) {
	query.GroupQuery = gq
}

func (gq GroupQuery) Having(filters ...FilterQuery) GroupQuery {
	gq.Filter = gq.Filter.And(filters...)
	return gq
}

func (gq GroupQuery) OrHaving(filters ...FilterQuery) GroupQuery {
	gq.Filter = gq.Filter.Or(FilterAnd(filters...))
	return gq
}

func (gq GroupQuery) Where(filters ...FilterQuery) GroupQuery {
	return gq.Having(filters...)
}

func (gq GroupQuery) OrWhere(filters ...FilterQuery) GroupQuery {
	return gq.OrHaving(filters...)
}

func NewGroup(fields ...string) GroupQuery {
	return GroupQuery{
		Fields: fields,
	}
}

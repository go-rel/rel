package query

type Group struct {
	Fields  []string
	Filters []Filter
}

func (group Group) Build(query *Query) {
	query.GroupClause = group
}

func (group Group) Having(filters ...Filter) Group {
	group.Filters = filters
	return group
}

func (group Group) Where(filters ...Filter) Group {
	return group.Having(filters...)
}

func GroupBy(fields ...string) Group {
	return Group{
		Fields: fields,
	}
}

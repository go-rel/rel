package query

type Limit int

func (limit Limit) Build(query *Query) {
	query.LimitClause = limit
}

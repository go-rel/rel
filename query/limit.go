package query

type LimitClause int

func (l LimitClause) Build(query *Query) {
	query.LimitClause = l
}

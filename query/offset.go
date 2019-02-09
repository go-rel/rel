package query

type OffsetClause int

func (o OffsetClause) Build(query *Query) {
	query.OffsetClause = o
}

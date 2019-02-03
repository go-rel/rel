package query

type Offset int

func (offset Offset) Build(query *Query) {
	query.OffsetClause = offset
}

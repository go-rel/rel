package query

type AggregateClause struct {
	Field string
	Mode  string
}

func (a AggregateClause) Build(query *Query) {
	query.AggregateClause = a
}

func Count() AggregateClause {
	return AggregateClause{
		Mode: "COUNT",
	}
}

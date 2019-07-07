package query

type Builder interface {
	Build(*Query)
}

func Build(collection string, builders ...Builder) Query {
	q := Query{
		empty: true,
	}

	for _, builder := range builders {
		builder.Build(&q)
		q.empty = false
	}

	if q.Collection == "" {
		q.Collection = collection
		q.empty = false
	}

	for i := range q.JoinClause {
		q.JoinClause[i].buildJoin(q)
	}

	return q
}

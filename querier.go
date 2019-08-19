package grimoire

type Querier interface {
	Build(*Query)
}

func BuildQuery(collection string, queriers ...Querier) Query {
	q := Query{
		empty: true,
	}

	for _, querier := range queriers {
		querier.Build(&q)
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

package query

type Builder interface {
	Build(*Query)
}

func Build(builders ...Builder) Query {
	var query Query
	for _, builder := range builders {
		builder.Build(&query)
		query.built = true
	}

	return query
}

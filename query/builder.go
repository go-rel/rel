package query

type Builder interface {
	Build(*Query)
}

func Build(builders ...Builder) Query {
	q := Query{
		empty: true,
	}

	for _, builder := range builders {
		builder.Build(&q)
	}

	return q
}

package query

type SelectClause struct {
	Distinct bool
	Fields   []string
}

func Select(fields ...string) SelectClause {
	return SelectClause{
		Fields: fields,
	}
}

func SelectDistinct(fields ...string) SelectClause {
	return SelectClause{
		Distinct: true,
		Fields:   fields,
	}
}

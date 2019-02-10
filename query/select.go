package query

type SelectClause struct {
	OnlyDistinct bool
	Fields       []string
}

func (s SelectClause) Distinct() SelectClause {
	s.OnlyDistinct = true
	return s
}

func Select(fields ...string) SelectClause {
	return SelectClause{
		Fields: fields,
	}
}

func Distinct() SelectClause {
	return SelectClause{
		OnlyDistinct: true,
	}
}

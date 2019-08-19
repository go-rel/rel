package grimoire

type SelectQuery struct {
	OnlyDistinct bool
	Fields       []string
}

func (sq SelectQuery) Distinct() SelectQuery {
	sq.OnlyDistinct = true
	return sq
}

func NewSelect(fields ...string) SelectQuery {
	return SelectQuery{
		Fields: fields,
	}
}

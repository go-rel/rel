package query

// SortClause defines sort information of query.
type SortClause struct {
	Field string
	Sort  int
}

func (s SortClause) Build(query *Query) {
	query.SortClause = append(query.SortClause, s)
}

// Asc returns true if sort is ascending.
func (s SortClause) Asc() bool {
	return s.Sort >= 0
}

// Desc returns true if s is descending.
func (s SortClause) Desc() bool {
	return s.Sort < 0
}

// SortAsc sorts field with ascending sort.
func SortAsc(field string) SortClause {
	return SortClause{
		Field: field,
		Sort:  1,
	}
}

// SortDesc sorts field with descending sort.
func SortDesc(field string) SortClause {
	return SortClause{
		Field: field,
		Sort:  -1,
	}
}

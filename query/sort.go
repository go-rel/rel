package query

// Sort defines sort information of query.
type Sort struct {
	Field string
	Sort  int
}

func (sort Sort) Build(query *Query) {
	query.SortClause = append(query.SortClause, sort)
}

// SortAsc sorts field with ascending sort.
func SortAsc(field string) Sort {
	return Sort{
		Field: field,
		Sort:  1,
	}
}

// SortDesc sorts field with descending sort.
func SortDesc(field string) Sort {
	return Sort{
		Field: field,
		Sort:  -1,
	}
}

// Asc returns true if sort is ascending.
func (sort Sort) Asc() bool {
	return sort.Sort >= 0
}

// Desc returns true if sort is descending.
func (sort Sort) Desc() bool {
	return sort.Sort < 0
}

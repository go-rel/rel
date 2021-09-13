package rel

// SortQuery defines sort information of query.
type SortQuery struct {
	Field string
	Sort  int
}

// Build sort query.
func (sq SortQuery) Build(query *Query) {
	query.SortQuery = append(query.SortQuery, sq)
}

// Asc returns true if sort is ascending.
func (sq SortQuery) Asc() bool {
	return sq.Sort >= 0
}

// Desc returns true if s is descending.
func (sq SortQuery) Desc() bool {
	return sq.Sort < 0
}

// SortAsc sorts field with ascending sort.
func SortAsc(field string) SortQuery {
	return SortQuery{
		Field: field,
		Sort:  1,
	}
}

// NewSortDesc sorts field with descending sort.
func SortDesc(field string) SortQuery {
	return SortQuery{
		Field: field,
		Sort:  -1,
	}
}

var (
	// NewSortAsc sorts field with ascending sort.
	//
	// Deprecated: use SortAsc instead
	NewSortAsc = SortAsc

	// NewSortDesc sorts field with descending sort.
	//
	// Deprecated: use SortDesc instead
	NewSortDesc = SortDesc
)

package rel

// Offset  Query.
type Offset int

// Build query.
func (o Offset) Build(query *Query) {
	query.OffsetQuery = o
}

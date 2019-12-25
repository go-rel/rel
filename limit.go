package rel

// Limit query.
type Limit int

// Build query.
func (l Limit) Build(query *Query) {
	query.LimitQuery = l
}

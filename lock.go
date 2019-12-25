package rel

// Lock query.
// This query will be ignored if used outside of transaction.
type Lock string

// Build query.
func (l Lock) Build(query *Query) {
	query.LockQuery = l
}

// ForUpdate lock query.
func ForUpdate() Lock {
	return "FOR UPDATE"
}

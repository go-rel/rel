package rel

// SQLQuery allows querying using native query supported by database.
type SQLQuery struct {
	Statement string
	Values    []interface{}
}

// Build Raw Query.
func (sq SQLQuery) Build(query *Query) {
	query.SQLQuery = sq
}

// SQL Query.
func SQL(statement string, values ...interface{}) SQLQuery {
	return SQLQuery{
		Statement: statement,
		Values:    values,
	}
}

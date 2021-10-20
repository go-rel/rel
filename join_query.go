package rel

// JoinQuery defines join clause in query.
type JoinQuery struct {
	Mode      string
	Table     string
	From      string
	To        string
	Filter    FilterQuery
	Arguments []interface{}
}

// Build query.
func (jq JoinQuery) Build(query *Query) {
	query.JoinQuery = append(query.JoinQuery, jq)
}

// NewJoinWith query with custom join mode, table, field and additional filters with AND condition.
func NewJoinWith(mode string, table string, from string, to string, filter ...FilterQuery) JoinQuery {
	return JoinQuery{
		Mode:   mode,
		Table:  table,
		From:   from,
		To:     to,
		Filter: And(filter...),
	}
}

// NewJoinFragment defines a join clause using raw query.
func NewJoinFragment(expr string, args ...interface{}) JoinQuery {
	if args == nil {
		// prevent buildJoin to populate From and To variable.
		args = []interface{}{}
	}

	return JoinQuery{
		Mode:      expr,
		Arguments: args,
	}
}

// NewJoin with given table.
func NewJoin(table string, filter ...FilterQuery) JoinQuery {
	return NewJoinWith("JOIN", table, "", "", filter...)
}

// NewJoinOn table with given field and optional additional filter.
func NewJoinOn(table string, from string, to string, filter ...FilterQuery) JoinQuery {
	return NewJoinWith("JOIN", table, from, to, filter...)
}

// NewInnerJoin with given table and optional filter.
func NewInnerJoin(table string, filter ...FilterQuery) JoinQuery {
	return NewInnerJoinOn(table, "", "", filter...)
}

// NewInnerJoinOn table with given field and optional additional filter.
func NewInnerJoinOn(table string, from string, to string, filter ...FilterQuery) JoinQuery {
	return NewJoinWith("INNER JOIN", table, from, to, filter...)
}

// NewLeftJoin with given table and optional filter.
func NewLeftJoin(table string, filter ...FilterQuery) JoinQuery {
	return NewLeftJoinOn(table, "", "", filter...)
}

// NewLeftJoinOn table with given field and optional additional filter.
func NewLeftJoinOn(table string, from string, to string, filter ...FilterQuery) JoinQuery {
	return NewJoinWith("LEFT JOIN", table, from, to, filter...)
}

// NewRightJoin with given table and optional filter.
func NewRightJoin(table string, filter ...FilterQuery) JoinQuery {
	return NewRightJoinOn(table, "", "", filter...)
}

// NewRightJoinOn table with given field and optional additional filter.
func NewRightJoinOn(table string, from string, to string, filter ...FilterQuery) JoinQuery {
	return NewJoinWith("RIGHT JOIN", table, from, to, filter...)
}

// NewFullJoin with given table and optional filter.
func NewFullJoin(table string, filter ...FilterQuery) JoinQuery {
	return NewFullJoinOn(table, "", "", filter...)
}

// NewFullJoinOn table with given field and optional additional filter.
func NewFullJoinOn(table string, from string, to string, filter ...FilterQuery) JoinQuery {
	return NewJoinWith("FULL JOIN", table, from, to, filter...)
}

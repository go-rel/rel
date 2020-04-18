package rel

// JoinQuery defines join clause in query.
type JoinQuery struct {
	Mode      string
	Table     string
	From      string
	To        string
	Arguments []interface{}
}

// Build query.
func (jq JoinQuery) Build(query *Query) {
	query.JoinQuery = append(query.JoinQuery, jq)
}

// NewJoinWith query with custom join mode, table and field.
func NewJoinWith(mode string, table string, from string, to string) JoinQuery {
	return JoinQuery{
		Mode:  mode,
		Table: table,
		From:  from,
		To:    to,
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
func NewJoin(table string) JoinQuery {
	return NewJoinWith("JOIN", table, "", "")
}

// NewJoinOn table with given field.
func NewJoinOn(table string, from string, to string) JoinQuery {
	return NewJoinWith("JOIN", table, from, to)
}

// NewInnerJoin with given table.
func NewInnerJoin(table string) JoinQuery {
	return NewInnerJoinOn(table, "", "")
}

// NewInnerJoinOn table with given field.
func NewInnerJoinOn(table string, from string, to string) JoinQuery {
	return NewJoinWith("INNER JOIN", table, from, to)
}

// NewLeftJoin with given table.
func NewLeftJoin(table string) JoinQuery {
	return NewLeftJoinOn(table, "", "")
}

// NewLeftJoinOn table with given field.
func NewLeftJoinOn(table string, from string, to string) JoinQuery {
	return NewJoinWith("LEFT JOIN", table, from, to)
}

// NewRightJoin with given table.
func NewRightJoin(table string) JoinQuery {
	return NewRightJoinOn(table, "", "")
}

// NewRightJoinOn table with given field.
func NewRightJoinOn(table string, from string, to string) JoinQuery {
	return NewJoinWith("RIGHT JOIN", table, from, to)
}

// NewFullJoin with given table.
func NewFullJoin(table string) JoinQuery {
	return NewFullJoinOn(table, "", "")
}

// NewFullJoinOn table with given field.
func NewFullJoinOn(table string, from string, to string) JoinQuery {
	return NewJoinWith("FULL JOIN", table, from, to)
}

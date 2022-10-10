package rel

// JoinQuery defines join clause in query.
type JoinQuery struct {
	Mode      string
	Table     string
	From      string
	To        string
	Assoc     string
	Filter    FilterQuery
	Arguments []any
}

// Build query.
func (jq JoinQuery) Build(query *Query) {
	query.JoinQuery = append(query.JoinQuery, jq)

	if jq.Assoc != "" {
		query.AddPopulator(&query.JoinQuery[len(query.JoinQuery)-1])
	}
}

func (jq *JoinQuery) Populate(query *Query, docMeta DocumentMeta) {
	var (
		assocMeta    = docMeta.Association(jq.Assoc)
		assocDocMeta = assocMeta.DocumentMeta()
	)

	jq.Table = assocDocMeta.Table() + " as " + jq.Assoc
	jq.To = jq.Assoc + "." + assocMeta.ForeignField()
	jq.From = docMeta.Table() + "." + assocMeta.ReferenceField()

	// load association if defined and supported
	if assocMeta.Type() == HasOne || assocMeta.Type() == BelongsTo {
		var (
			load        = false
			selectField = jq.Assoc + ".*"
		)

		for i := range query.SelectQuery.Fields {
			if load && i > 0 {
				query.SelectQuery.Fields[i-1] = query.SelectQuery.Fields[i]
			}
			if query.SelectQuery.Fields[i] == selectField {
				load = true
			}
		}

		if load {
			fields := make([]string, len(assocDocMeta.Fields()))
			for i, f := range assocDocMeta.Fields() {
				fields[i] = jq.Assoc + "." + f + " as " + jq.Assoc + "." + f
			}
			query.SelectQuery.Fields = append(query.SelectQuery.Fields[:(len(query.SelectQuery.Fields)-1)], fields...)
		}
	}
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
func NewJoinFragment(expr string, args ...any) JoinQuery {
	if args == nil {
		// prevent buildJoin to populate From and To variable.
		args = []any{}
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

// NewJoinAssocWith with given association field and optional additional filters.
func NewJoinAssocWith(mode string, assoc string, filter ...FilterQuery) JoinQuery {
	return JoinQuery{
		Mode:   mode,
		Assoc:  assoc,
		Filter: And(filter...),
	}
}

// NewJoinAssoc with given association field and optional additional filters.
func NewJoinAssoc(assoc string, filter ...FilterQuery) JoinQuery {
	return NewJoinAssocWith("JOIN", assoc, filter...)
}

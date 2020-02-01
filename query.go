package rel

// Querier interface defines contract to be used for query builder.
type Querier interface {
	Build(*Query)
}

// Build for given table using given queriers.
func Build(table string, queriers ...Querier) Query {
	var (
		query Query
	)

	if len(queriers) > 0 {
		_, query.empty = queriers[0].(Query)
	}

	for _, querier := range queriers {
		// avoid using indirect call to avoid heap allocation
		switch q := querier.(type) {
		case Query:
			q.Build(&query)
		case JoinQuery:
			q.Build(&query)
		case FilterQuery:
			q.Build(&query)
		case GroupQuery:
			q.Build(&query)
		case SortQuery:
			q.Build(&query)
		case Offset:
			q.Build(&query)
		case Limit:
			q.Build(&query)
		case Lock:
			q.Build(&query)
		}
	}

	if query.Table == "" {
		query.Table = table
	}

	for i := range query.JoinQuery {
		query.JoinQuery[i].buildJoin(query)
	}

	return query
}

// Query defines information about query generated by query builder.
type Query struct {
	empty         bool // todo: use bit to mark what is updated and use it when building
	Table         string
	SelectQuery   SelectQuery
	JoinQuery     []JoinQuery
	WhereQuery    FilterQuery
	GroupQuery    GroupQuery
	SortQuery     []SortQuery
	OffsetQuery   Offset
	LimitQuery    Limit
	LockQuery     Lock
	UnscopedQuery Unscoped
}

// Build query.
func (q Query) Build(query *Query) {
	if query.empty {
		*query = q
	} else {
		// manual merge
		if q.Table != "" {
			query.Table = q.Table
		}

		if q.SelectQuery.Fields != nil {
			query.SelectQuery = q.SelectQuery
		}

		query.JoinQuery = append(query.JoinQuery, q.JoinQuery...)

		query.WhereQuery = query.WhereQuery.And(q.WhereQuery)

		if q.GroupQuery.Fields != nil {
			query.GroupQuery = q.GroupQuery
		}

		q.SortQuery = append(q.SortQuery, query.SortQuery...)

		if q.OffsetQuery != 0 {
			query.OffsetQuery = q.OffsetQuery
		}

		if q.LimitQuery != 0 {
			query.LimitQuery = q.LimitQuery
		}

		if q.LockQuery != "" {
			query.LockQuery = q.LockQuery
		}
	}
}

// Select filter fields to be selected from database.
func (q Query) Select(fields ...string) Query {
	q.SelectQuery = NewSelect(fields...)
	return q
}

// From set the table to be used for query.
func (q Query) From(table string) Query {
	q.Table = table
	return q
}

// Distinct sets select query to be distinct.
func (q Query) Distinct() Query {
	q.SelectQuery.OnlyDistinct = true
	return q
}

// Join current table with other table.
func (q Query) Join(table string) Query {
	return q.JoinOn(table, "", "")
}

// JoinOn current table with other table.
func (q Query) JoinOn(table string, from string, to string) Query {
	return q.JoinWith("JOIN", table, from, to)
}

// JoinWith current table with other table with custom join mode.
func (q Query) JoinWith(mode string, table string, from string, to string) Query {
	NewJoinWith(mode, table, from, to).Build(&q) // TODO: ensure this always called last

	return q
}

// Joinf create join query using a raw query.
func (q Query) Joinf(expr string, args ...interface{}) Query {
	NewJoinFragment(expr, args...).Build(&q) // TODO: ensure this always called last

	return q
}

// Where query.
func (q Query) Where(filters ...FilterQuery) Query {
	q.WhereQuery = q.WhereQuery.And(filters...)
	return q
}

// Wheref create where query using a raw query.
func (q Query) Wheref(expr string, args ...interface{}) Query {
	q.WhereQuery = q.WhereQuery.And(FilterFragment(expr, args...))
	return q
}

// OrWhere query.
func (q Query) OrWhere(filters ...FilterQuery) Query {
	q.WhereQuery = q.WhereQuery.Or(And(filters...))
	return q
}

// OrWheref create where query using a raw query.
func (q Query) OrWheref(expr string, args ...interface{}) Query {
	q.WhereQuery = q.WhereQuery.Or(FilterFragment(expr, args...))
	return q
}

// Group query.
func (q Query) Group(fields ...string) Query {
	q.GroupQuery.Fields = fields
	return q
}

// Having query.
func (q Query) Having(filters ...FilterQuery) Query {
	q.GroupQuery.Filter = q.GroupQuery.Filter.And(filters...)
	return q
}

// Havingf create having query using a raw query.
func (q Query) Havingf(expr string, args ...interface{}) Query {
	q.GroupQuery.Filter = q.GroupQuery.Filter.And(FilterFragment(expr, args...))
	return q
}

// OrHaving query.
func (q Query) OrHaving(filters ...FilterQuery) Query {
	q.GroupQuery.Filter = q.GroupQuery.Filter.Or(And(filters...))
	return q
}

// OrHavingf create having query using a raw query.
func (q Query) OrHavingf(expr string, args ...interface{}) Query {
	q.GroupQuery.Filter = q.GroupQuery.Filter.Or(FilterFragment(expr, args...))
	return q
}

// Sort query.
func (q Query) Sort(fields ...string) Query {
	return q.SortAsc(fields...)
}

// SortAsc query.
func (q Query) SortAsc(fields ...string) Query {
	var (
		offset = len(q.SortQuery)
	)

	q.SortQuery = append(q.SortQuery, make([]SortQuery, len(fields))...)
	for i := range fields {
		q.SortQuery[offset+i] = NewSortAsc(fields[i])
	}

	return q
}

// SortDesc query.
func (q Query) SortDesc(fields ...string) Query {
	var (
		offset = len(q.SortQuery)
	)

	q.SortQuery = append(q.SortQuery, make([]SortQuery, len(fields))...)
	for i := range fields {
		q.SortQuery[offset+i] = NewSortDesc(fields[i])
	}

	return q
}

// Offset the result returned by database.
func (q Query) Offset(offset Offset) Query {
	q.OffsetQuery = offset
	return q
}

// Limit result returned by database.
func (q Query) Limit(limit Limit) Query {
	q.LimitQuery = limit
	return q
}

// Lock query expression.
func (q Query) Lock(lock Lock) Query {
	q.LockQuery = lock
	return q
}

// Unscoped allows soft-delete to be ignored.
func (q Query) Unscoped() Query {
	q.UnscopedQuery = true
	return q
}

// Select query create a query with chainable syntax, using select as the starting point.
func Select(fields ...string) Query {
	return Query{
		SelectQuery: SelectQuery{
			Fields: fields,
		},
	}
}

// From create a query with chainable syntax, using from as the starting point.
func From(table string) Query {
	return Query{
		Table: table,
	}
}

// Join create a query with chainable syntax, using join as the starting point.
func Join(table string) Query {
	return JoinOn(table, "", "")
}

// JoinOn create a query with chainable syntax, using join as the starting point.
func JoinOn(table string, from string, to string) Query {
	return JoinWith("JOIN", table, from, to)
}

// JoinWith create a query with chainable syntax, using join as the starting point.
func JoinWith(mode string, table string, from string, to string) Query {
	return Query{
		JoinQuery: []JoinQuery{
			NewJoinWith(mode, table, from, to),
		},
	}
	// var q Query
	// NewJoinWith(mode, table, from, to).Build(&q) // TODO: ensure this always called last

	// return q
}

// Joinf create a query with chainable syntax, using join as the starting point.
func Joinf(expr string, args ...interface{}) Query {
	return Query{
		JoinQuery: []JoinQuery{
			NewJoinFragment(expr, args...),
		},
	}

	// var q Query
	// NewJoinFragment(expr, args...).Build(&q) // TODO: ensure this always called last

	// return q
}

// Where create a query with chainable syntax, using where as the starting point.
func Where(filters ...FilterQuery) Query {
	return Query{
		WhereQuery: And(filters...),
	}
}

// Offset  Query.
type Offset int

// Build query.
func (o Offset) Build(query *Query) {
	query.OffsetQuery = o
}

// Limit query.
type Limit int

// Build query.
func (l Limit) Build(query *Query) {
	query.LimitQuery = l
}

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

// Unscoped query.
type Unscoped bool

// Build query.
func (u Unscoped) Build(query *Query) {
	query.UnscopedQuery = u
}

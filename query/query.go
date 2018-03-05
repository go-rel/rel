package query

type Query struct {
	Collection      string
	Fields          []string
	JoinQuery       []JoinQuery
	Condition       Condition
	GroupFields     []string
	HavingCondition Condition
	OrderQuery      []OrderQuery
	OffsetResult    int
	LimitResut      int
}

type JoinQuery struct {
	Mode       string
	Collection string
	Condition  Condition
}

type OrderQuery struct {
	Field string
	Order int
}

func From(collection string) Query {
	return Query{
		Collection: collection,
		Fields:     []string{"*"},
	}
}

func (q Query) Select(fields ...string) Query {
	q.Fields = fields
	return q
}

func (q Query) Join(collection string, condition ...Condition) Query {
	return q.JoinWith("JOIN", collection, condition...)
}

func (q Query) JoinWith(mode string, collection string, condition ...Condition) Query {
	q.JoinQuery = append(q.JoinQuery, JoinQuery{
		Mode:       mode,
		Collection: collection,
		Condition:  And(condition...),
	})

	return q
}

// Where expressions are used to filter the result set. If there is more than one where expression, they are combined with an and operator
func (q Query) Where(condition ...Condition) Query {
	q.Condition = q.Condition.And(condition...)
	return q
}

// OrWhere behaves exactly the same as where except it combines with any previous expression by using an OR
func (q Query) OrWhere(condition ...Condition) Query {
	q.Condition = q.Condition.Or(condition...)
	return q
}

func (q Query) GroupBy(fields ...string) Query {
	q.GroupFields = fields
	return q
}

func (q Query) Having(condition ...Condition) Query {
	q.HavingCondition = q.HavingCondition.And(condition...)
	return q
}

func (q Query) OrHaving(condition ...Condition) Query {
	q.HavingCondition = q.HavingCondition.Or(condition...)
	return q
}

func (q Query) OrderBy(order ...OrderQuery) Query {
	q.OrderQuery = append(q.OrderQuery, order...)
	return q
}

func (q Query) Offset(offset int) Query {
	q.OffsetResult = offset
	return q
}

func (q Query) Limit(limit int) Query {
	q.LimitResut = limit
	return q
}

func Asc(field string) OrderQuery {
	return OrderQuery{
		Field: field,
		Order: 1,
	}
}

func Desc(field string) OrderQuery {
	return OrderQuery{
		Field: field,
		Order: -1,
	}
}

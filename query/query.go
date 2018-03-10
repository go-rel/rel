package query

type Query struct {
	Collection      string
	Fields          []string
	AsDistinct      bool
	JoinClause      []JoinClause
	Condition       Condition
	GroupFields     []string
	HavingCondition Condition
	OrderClause     []OrderClause
	OffsetResult    int
	LimitResult     int
}

type JoinClause struct {
	Mode       string
	Collection string
	Condition  Condition
}

type OrderClause struct {
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

func (q Query) Distinct() Query {
	q.AsDistinct = true
	return q
}

func (q Query) Join(collection string, condition ...Condition) Query {
	return q.JoinWith("JOIN", collection, condition...)
}

func (q Query) JoinWith(mode string, collection string, condition ...Condition) Query {
	q.JoinClause = append(q.JoinClause, JoinClause{
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
	q.Condition = q.Condition.Or(And(condition...))
	return q
}

func (q Query) Group(fields ...string) Query {
	q.GroupFields = fields
	return q
}

func (q Query) Having(condition ...Condition) Query {
	q.HavingCondition = q.HavingCondition.And(condition...)
	return q
}

func (q Query) OrHaving(condition ...Condition) Query {
	q.HavingCondition = q.HavingCondition.Or(And(condition...))
	return q
}

func (q Query) Order(order ...OrderClause) Query {
	q.OrderClause = append(q.OrderClause, order...)
	return q
}

func (q Query) Offset(offset int) Query {
	q.OffsetResult = offset
	return q
}

func (q Query) Limit(limit int) Query {
	q.LimitResult = limit
	return q
}

func Asc(field string) OrderClause {
	return OrderClause{
		Field: field,
		Order: 1,
	}
}

func Desc(field string) OrderClause {
	return OrderClause{
		Field: field,
		Order: -1,
	}
}

func (o OrderClause) Asc() bool {
	return o.Order >= 0
}

func (o OrderClause) Desc() bool {
	return o.Order < 0
}

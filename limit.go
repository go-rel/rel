package rel

type Limit int

func (l Limit) Build(query *Query) {
	query.LimitQuery = l
}

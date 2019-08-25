package grimoire

type Offset int

func (o Offset) Build(query *Query) {
	query.OffsetQuery = o
}

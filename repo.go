package grimoire

type Repo struct {
	adapter Adapter
}

func (repo Repo) From(collection string) Query {
	return Query{
		Collection: collection,
		Fields:     []string{"*"},
	}
}

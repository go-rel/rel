package grimoire

type Repo struct {
	adapter Adapter
}

func New(adapter Adapter) Repo {
	return Repo{
		adapter: adapter,
	}
}

func (repo Repo) Close() error {
	return repo.adapter.Close()
}

func (repo Repo) From(collection string) Query {
	return Query{
		Collection: collection,
		Fields:     []string{"*"},
	}
}

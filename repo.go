package grimoire

type Repo struct {
	adapter       Adapter
	inTransaction bool
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

func (repo Repo) Transaction(fn func(Repo) error) error {
	err := repo.adapter.Begin()
	if err != nil {
		return err
	}

	repo.inTransaction = true

	defer func() {
		if p := recover(); p != nil {
			repo.adapter.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			repo.adapter.Rollback()
		} else {
			err = repo.adapter.Commit()
		}
		repo.inTransaction = false
	}()

	err = fn(repo)
	return err
}

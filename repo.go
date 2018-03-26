package grimoire

import (
	"github.com/Fs02/grimoire/errors"
)

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
		repo:       &repo,
		Collection: collection,
		Fields:     []string{"*"},
	}
}

func (repo Repo) Transaction(fn func(Repo) error) error {
	err := repo.adapter.Begin()
	if err != nil {
		return err
	}

	func() {
		defer func() {
			if p := recover(); p != nil {
				repo.adapter.Rollback()

				if e, ok := p.(errors.Error); ok && !e.UnexpectedError() {
					err = e
				} else {
					panic(p) // re-throw panic after Rollback
				}
			} else if err != nil {
				repo.adapter.Rollback()
			} else {
				err = repo.adapter.Commit()
			}
		}()

		err = fn(repo)
	}()

	return err
}

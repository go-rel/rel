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

func (repo Repo) From(collection string) Query {
	return Query{
		repo:       &repo,
		Collection: collection,
		Fields:     []string{"*"},
	}
}

func (repo Repo) Transaction(fn func(Repo) error) error {
	adp, err := repo.adapter.Begin()
	if err != nil {
		return err
	}

	txRepo := New(adp)

	func() {
		defer func() {
			if p := recover(); p != nil {
				txRepo.adapter.Rollback()

				if e, ok := p.(errors.Error); ok && !e.UnexpectedError() {
					err = e
				} else {
					panic(p) // re-throw panic after Rollback
				}
			} else if err != nil {
				txRepo.adapter.Rollback()
			} else {
				err = txRepo.adapter.Commit()
			}
		}()

		err = fn(txRepo)
	}()

	return err
}

package grimoire

import (
	"github.com/Fs02/grimoire/errors"
)

// Repo defines grimoire repository.
type Repo struct {
	adapter Adapter
	logger  []Logger
}

// New create new repo using adapter.
func New(adapter Adapter) Repo {
	return Repo{
		adapter: adapter,
		logger:  []Logger{DefaultLogger},
	}
}

// Adapter returns adapter of repo.
func (repo *Repo) Adapter() Adapter {
	return repo.adapter
}

// SetLogger replace default logger with custom logger.
func (repo *Repo) SetLogger(logger ...Logger) {
	repo.logger = logger
}

// From initiates a query for a collection.
func (repo Repo) From(collection string) Query {
	return Query{
		repo:       &repo,
		Collection: collection,
		Fields:     []string{collection + ".*"},
	}
}

// Transaction performs transaction with given function argument.
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

				if e, ok := p.(errors.Error); ok && e.Kind() != errors.Unexpected {
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

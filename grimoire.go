package grimoire

import (
	"github.com/Fs02/grimoire/query"
)

// Begin() (*Repo, error)
// Commit() (*Repo, error)
// Rollback() (*Repo, error)
// CommitOrRollback() (*Repo, error)

// All(*structs, query) error
// One(*struct, query) error
// Insert(*struct, changeset) error
// Update(*struct, changeset) error
// UpdateAll(*struct, changeset, condition) error
// Delete(*struct) error
// DeleteAll(*struct, condition) error

// res := struct{id int; name string}{}
// Repo.All(&res, query)

// Adapter abstraction
// accepts struct and query or changeset
// returns query string and arguments
type Adapter interface {
	Open(string) error
	Close() error

	All(query.Query) (string, []interface{})
	Insert(*Changeset) (string, []interface{})
	Update(*Changeset, query.Condition) (string, []interface{})
	Delete(query.Condition) (string, []interface{})

	// Begin() (*Repo, error)
	// Commit() (*Repo, error)
	// Rollback() (*Repo, error)

	// Query exec query string with it's arguments
	// reurns results and an error if any
	Query(interface{}, string, []interface{}) error

	// Query exec query string with it's arguments
	// returns last inserted id, rows affected and error
	Exec(string, []interface{}) (int64, int64, error)
}

type Repo struct {
	adapter Adapter
}

func (r Repo) All(entities interface{}, q query.Query) error {
	qs, args := r.adapter.All(q)
	return r.adapter.Query(entities, qs, args)
}

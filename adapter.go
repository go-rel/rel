package grimoire

import (
	"github.com/Fs02/grimoire/changeset"
)

// Adapter abstraction
// accepts struct and query or changeset
// returns query string and arguments
type Adapter interface {
	Open(string) error
	Close() error

	Find(Query) (string, []interface{})
	Insert(Query, changeset.Changeset) (string, []interface{})
	Update(Query, changeset.Changeset) (string, []interface{})
	Delete(Query) (string, []interface{})

	Begin() error
	Commit() error
	Rollback() error

	// Query exec query string with it's arguments
	// reurns results and an error if any
	Query(interface{}, string, []interface{}) (int64, error)

	// Query exec query string with it's arguments
	// returns last inserted id, rows affected and error
	Exec(string, []interface{}) (int64, int64, error)
}

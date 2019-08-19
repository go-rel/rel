package grimoire

import (
	"github.com/Fs02/grimoire/change"
)

// Adapter interface
type Adapter interface {
	Aggregate(Query, string, string, ...Logger) (int, error)
	Query(Query, ...Logger) (Cursor, error)
	Insert(Query, change.Changes, ...Logger) (interface{}, error)
	InsertAll(Query, []string, []change.Changes, ...Logger) ([]interface{}, error)
	Update(Query, change.Changes, ...Logger) error
	Delete(Query, ...Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

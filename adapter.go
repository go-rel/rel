package grimoire

import (
	"github.com/Fs02/grimoire/change"
	"github.com/Fs02/grimoire/query"
)

// Adapter interface
type Adapter interface {
	Aggregate(query.Query, string, string, ...Logger) (int, error)
	Query(query.Query, ...Logger) (Cursor, error)
	Insert(query.Query, change.Changes, ...Logger) (interface{}, error)
	InsertAll(query.Query, []string, []change.Changes, ...Logger) ([]interface{}, error)
	Update(query.Query, change.Changes, ...Logger) error
	Delete(query.Query, ...Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

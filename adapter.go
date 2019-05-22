package grimoire

import (
	"github.com/Fs02/grimoire/query"
)

// Adapter interface
type Adapter interface {
	Aggregate(query.Query, interface{}, string, string, ...Logger) error
	All(query.Query, interface{}, ...Logger) (int, error)
	Insert(query.Query, map[string]interface{}, ...Logger) (interface{}, error)
	InsertAll(query.Query, []string, []map[string]interface{}, ...Logger) ([]interface{}, error)
	Update(query.Query, map[string]interface{}, ...Logger) error
	Delete(query.Query, ...Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

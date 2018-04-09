package grimoire

// Adapter abstraction
// accepts struct and query or changeset
// returns query string and arguments
type Adapter interface {
	All(Query, interface{}) (int, error)
	Delete(Query) error
	Insert(Query, map[string]interface{}) (interface{}, error)
	InsertAll(Query, []string, []map[string]interface{}) ([]interface{}, error)
	Update(Query, map[string]interface{}) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

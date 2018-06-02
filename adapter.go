package grimoire

// Adapter interface
type Adapter interface {
	All(Query, interface{}, ...Logger) (int, error)
	Aggregate(Query, interface{}, ...Logger) error
	Insert(Query, map[string]interface{}, ...Logger) (interface{}, error)
	InsertAll(Query, []string, []map[string]interface{}, ...Logger) ([]interface{}, error)
	Update(Query, map[string]interface{}, ...Logger) error
	Delete(Query, ...Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

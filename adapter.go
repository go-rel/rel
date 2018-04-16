package grimoire

// Adapter interface
type Adapter interface {
	Count(Query, Logger) (int, error)
	All(Query, interface{}, Logger) (int, error)
	Delete(Query, Logger) error
	Insert(Query, map[string]interface{}, Logger) (interface{}, error)
	InsertAll(Query, []string, []map[string]interface{}, Logger) ([]interface{}, error)
	Update(Query, map[string]interface{}, Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

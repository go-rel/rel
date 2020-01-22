package rel

// Adapter interface
type Adapter interface {
	Aggregate(query Query, mode string, field string, loggers ...Logger) (int, error)
	Query(query Query, loggers ...Logger) (Cursor, error)
	Insert(query Query, modifies map[string]Modify, loggers ...Logger) (interface{}, error)
	InsertAll(query Query, fields []string, bulkModifies []map[string]Modify, loggers ...Logger) ([]interface{}, error)
	Update(query Query, modifies map[string]Modify, loggers ...Logger) error
	Delete(query Query, loggers ...Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

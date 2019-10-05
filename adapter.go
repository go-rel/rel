package rel

// Adapter interface
type Adapter interface {
	Aggregate(query Query, mode string, field string, loggers ...Logger) (int, error)
	Query(query Query, loggers ...Logger) (Cursor, error)
	Insert(query Query, changes Changes, loggers ...Logger) (interface{}, error)
	InsertAll(query Query, fields []string, changes []Changes, loggers ...Logger) ([]interface{}, error)
	Update(query Query, changes Changes, loggers ...Logger) error
	Delete(query Query, loggers ...Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

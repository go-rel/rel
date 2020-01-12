package rel

// Adapter interface
type Adapter interface {
	Aggregate(query Query, mode string, field string, loggers ...Logger) (int, error)
	Query(query Query, loggers ...Logger) (Cursor, error)
	Insert(query Query, modification Modification, loggers ...Logger) (interface{}, error)
	InsertAll(query Query, fields []string, modification []Modification, loggers ...Logger) ([]interface{}, error)
	Update(query Query, modification Modification, loggers ...Logger) error
	Delete(query Query, loggers ...Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

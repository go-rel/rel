package grimoire

// Adapter interface
type Adapter interface {
	Aggregate(Query, string, string, ...Logger) (int, error)
	Query(Query, ...Logger) (Cursor, error)
	Insert(Query, Changes, ...Logger) (interface{}, error)
	InsertAll(Query, []string, []Changes, ...Logger) ([]interface{}, error)
	Update(Query, Changes, ...Logger) error
	Delete(Query, ...Logger) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error
}

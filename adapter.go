package grimoire

// Adapter abstraction
// accepts struct and query or changeset
// returns query string and arguments
type Adapter interface {
	All(Query, interface{}) (int, error)
	Delete(Query) error
	Insert(Query, map[string]interface{}) (int, error)
	Update(Query, map[string]interface{}) error

	Begin() (Adapter, error)
	Commit() error
	Rollback() error

	// Query exec query string with it's arguments
	// reurns results and an error if any
	Query(interface{}, string, []interface{}) (int64, error)

	// Query exec query string with it's arguments
	// returns last inserted id, rows affected and error
	Exec(string, []interface{}) (int64, int64, error)
}

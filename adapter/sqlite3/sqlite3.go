package sqlite3

import (
	db "database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sql"
	_ "github.com/mattn/go-sqlite3"
)

// Adapter definition for mysql database.
type Adapter struct {
	*sql.Adapter
}

var _ grimoire.Adapter = (*Adapter)(nil)

// Open mysql connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var err error
	adapter := &Adapter{
		&sql.Adapter{
			Placeholder:   "?",
			IsOrdinal:     false,
			IncrementFunc: incrementFunc,
			ErrorFunc:     errorFunc,
		},
	}

	adapter.DB, err = db.Open("sqlite3", dsn)
	return adapter, err
}

func incrementFunc(adapter sql.Adapter) int {
	// decrement
	return -1
}

func errorFunc(err error) error {
	if err == nil {
		return nil
		// } else if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		// 	return errors.DuplicateError(e.Message, "")
	}

	return err
}

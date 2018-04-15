package sqlite3

import (
	db "database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sql"
	"github.com/Fs02/grimoire/errors"
	"github.com/mattn/go-sqlite3"
)

// Adapter definition for mysql database.
type Adapter struct {
	*sql.Adapter
}

var _ grimoire.Adapter = (*Adapter)(nil)

// Open mysql connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var err error

	adapter := &Adapter{sql.New("?", false, errorFunc, incrementFunc)}
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
	} else if e, ok := err.(sqlite3.Error); ok && e.ExtendedCode == sqlite3.ErrConstraintUnique {
		return errors.DuplicateError(e.Error(), "")
	}

	return err
}

package mysql

import (
	db "database/sql"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sql"
	"github.com/Fs02/grimoire/errors"
	"github.com/go-sql-driver/mysql"
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
	adapter.DB, err = db.Open("mysql", dsn)

	return adapter, err
}

func incrementFunc(adapter sql.Adapter) int {
	var variable string
	var increment int
	var err error
	if adapter.Tx != nil {
		err = adapter.Tx.QueryRow("SHOW VARIABLES LIKE 'auto_increment_increment';").Scan(&variable, &increment)
	} else {
		err = adapter.DB.QueryRow("SHOW VARIABLES LIKE 'auto_increment_increment';").Scan(&variable, &increment)
	}
	paranoid.Panic(err)

	return increment
}

func errorFunc(err error) error {
	if err == nil {
		return nil
	} else if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return errors.DuplicateError(e.Message, "")
	}

	return err
}

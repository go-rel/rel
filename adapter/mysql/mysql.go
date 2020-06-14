// Package mysql wraps mysql driver as an adapter for REL.
//
// Usage:
//	// open mysql connection.
//	adapter, err := mysql.Open("root@(127.0.0.1:3306)/rel_test?charset=utf8&parseTime=True&loc=Local")
//	if err != nil {
//		panic(err)
//	}
//	defer adapter.Close()
//
//	// initialize REL's repo.
//	repo := rel.New(adapter)
package mysql

import (
	db "database/sql"
	"strings"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/sql"
)

// Adapter definition for mysql database.
type Adapter struct {
	*sql.Adapter
}

var _ rel.Adapter = (*Adapter)(nil)

// New is mysql adapter constructor.
func New(database *db.DB) *Adapter {
	return &Adapter{
		Adapter: &sql.Adapter{
			Config: &sql.Config{
				Placeholder:   "?",
				EscapeChar:    "`",
				IncrementFunc: incrementFunc,
				ErrorFunc:     errorFunc,
			},
			DB: database,
		},
	}
}

// Open mysql connection using dsn.
func Open(dsn string) (*Adapter, error) {
	// force clientFoundRows=true
	// this allows not found record check when updating a record.
	if strings.ContainsRune(dsn, '?') {
		dsn += "&clientFoundRows=true"
	} else {
		dsn += "?clientFoundRows=true"
	}

	var database, err = db.Open("mysql", dsn)
	return New(database), err
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

	check(err)

	return increment
}

func errorFunc(err error) error {
	if err == nil {
		return nil
	}

	var (
		msg          = err.Error()
		errCodeSep   = ':'
		errCodeIndex = strings.IndexRune(msg, errCodeSep)
	)

	if errCodeIndex < 0 {
		errCodeIndex = 0
	}

	switch msg[:errCodeIndex] {
	case "Error 1062":
		return rel.ConstraintError{
			Key:  sql.ExtractString(msg, "key '", "'"),
			Type: rel.UniqueConstraint,
			Err:  err,
		}
	case "Error 1452":
		return rel.ConstraintError{
			Key:  sql.ExtractString(msg, "CONSTRAINT `", "`"),
			Type: rel.ForeignKeyConstraint,
			Err:  err,
		}
	default:
		return err
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

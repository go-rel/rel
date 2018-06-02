// Package sqlite3 wraps go-sqlite3 driver as an adapter for grimoire.
//
// Usage:
//	// open sqlite3 connection.
//	adapter, err := sqlite3.Open("dev.db")
//	if err != nil {
//		panic(err)
//	}
//	defer adapter.Close()
//
//	// initialize grimoire's repo.
//	repo := grimoire.New(adapter)
package sqlite3

import (
	db "database/sql"
	"strings"

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

	adapter := &Adapter{
		Adapter: &sql.Adapter{
			Config: &sql.Config{
				Placeholder:         "?",
				EscapeChar:          "`",
				InsertDefaultValues: true,
				IncrementFunc:       incrementFunc,
				ErrorFunc:           errorFunc,
			},
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
	}

	e, _ := err.(sqlite3.Error)
	switch e.ExtendedCode {
	case sqlite3.ErrConstraintUnique:
		return errors.New(e.Error(), strings.Split(e.Error(), "failed: ")[1], errors.UniqueConstraint)
	case sqlite3.ErrConstraintCheck:
		return errors.New(e.Error(), strings.Split(e.Error(), "failed: ")[1], errors.CheckConstraint)
	default:
		return err
	}
}

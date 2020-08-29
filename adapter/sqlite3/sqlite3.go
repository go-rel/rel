// Package sqlite3 wraps go-sqlite3 driver as an adapter for rel.
//
// Usage:
//	// open sqlite3 connection.
//	adapter, err := sqlite3.Open("dev.db")
//	if err != nil {
//		panic(err)
//	}
//	defer adapter.Close()
//
//	// initialize rel's repo.
//	repo := rel.New(adapter)
package sqlite3

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
				Placeholder:         "?",
				EscapeChar:          "`",
				InsertDefaultValues: true,
				IncrementFunc:       incrementFunc,
				ErrorFunc:           errorFunc,
				MapColumnTypeFunc:   mapColumnTypeFunc,
			},
			DB: database,
		},
	}
}

// Open mysql connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var database, err = db.Open("sqlite3", dsn)
	return New(database), err
}

func incrementFunc(adapter sql.Adapter) int {
	// decrement
	return -1
}

func errorFunc(err error) error {
	if err == nil {
		return nil
	}

	var (
		msg         = err.Error()
		failedSep   = " failed: "
		failedIndex = strings.Index(msg, failedSep)
		failedLen   = 9 // len(failedSep)
	)

	if failedIndex < 0 {
		failedIndex = 0
	}

	switch msg[:failedIndex] {
	case "UNIQUE constraint":
		return rel.ConstraintError{
			Key:  msg[failedIndex+failedLen:],
			Type: rel.UniqueConstraint,
			Err:  err,
		}
	case "CHECK constraint":
		return rel.ConstraintError{
			Key:  msg[failedIndex+failedLen:],
			Type: rel.CheckConstraint,
			Err:  err,
		}
	default:
		return err
	}
}

func mapColumnTypeFunc(column rel.Column) (string, int, int) {
	var (
		typ  string
		m, n int
	)

	switch column.Type {
	case rel.ID:
		typ = "INTEGER PRIMARY KEY"
	case rel.Int:
		typ = "INTEGER"
		m = column.Limit
	default:
		typ, m, n = sql.MapColumnType(column)
	}

	return typ, m, n
}

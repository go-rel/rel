// Package postgres wraps postgres (pq) driver as an adapter for rel.
//
// Usage:
//	// open postgres connection.
//	adapter, err := postgres.Open("postgres://postgres@localhost/rel_test?sslmode=disable")
//	if err != nil {
//		panic(err)
//	}
//	defer adapter.Close()
//
//	// initialize rel's repo.
//	repo := rel.New(adapter)
package postgres

import (
	db "database/sql"
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/sql"
)

// Adapter definition for postgrees database.
type Adapter struct {
	*sql.Adapter
}

var _ rel.Adapter = (*Adapter)(nil)

// Open postgrees connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var err error

	adapter := &Adapter{
		Adapter: &sql.Adapter{
			Config: &sql.Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
				ErrorFunc:           errorFunc,
			},
		},
	}
	adapter.DB, err = db.Open("postgres", dsn)

	return adapter, err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(query rel.Query, modifies map[string]rel.Modify, loggers ...rel.Logger) (interface{}, error) {
	var (
		id              int64
		statement, args = sql.NewBuilder(adapter.Config).Returning("id").Insert(query.Table, modifies)
		rows, err       = adapter.query(statement, args, loggers)
	)

	if err == nil && rows.Next() {
		defer rows.Close()
		rows.Scan(&id)
	}

	return id, err
}

// InsertAll inserts multiple records to database and returns its ids.
func (adapter *Adapter) InsertAll(query rel.Query, fields []string, bulkModifies []map[string]rel.Modify, loggers ...rel.Logger) ([]interface{}, error) {
	var (
		ids             []interface{}
		statement, args = sql.NewBuilder(adapter.Config).Returning("id").InsertAll(query.Table, fields, bulkModifies)
		rows, err       = adapter.query(statement, args, loggers)
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int64
			rows.Scan(&id)
			ids = append(ids, id)
		}
	}

	return ids, err
}

func (adapter *Adapter) query(statement string, args []interface{}, loggers []rel.Logger) (*db.Rows, error) {
	var (
		err   error
		rows  *db.Rows
		start = time.Now()
	)

	if adapter.Tx != nil {
		rows, err = adapter.Tx.Query(statement, args...)
	} else {
		rows, err = adapter.DB.Query(statement, args...)
	}

	go rel.Log(loggers, statement, time.Since(start), err)

	return rows, errorFunc(err)
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin() (rel.Adapter, error) {
	newAdapter, err := adapter.Adapter.Begin()

	return &Adapter{
		Adapter: newAdapter.(*sql.Adapter),
	}, err
}

func errorFunc(err error) error {
	if err == nil {
		return nil
	}

	var (
		msg            = err.Error()
		constraintType = sql.ExtractString(msg, "violates ", " constraint")
	)

	switch constraintType {
	case "unique":
		return rel.ConstraintError{
			Key:  sql.ExtractString(err.Error(), "constraint \"", "\""),
			Type: rel.UniqueConstraint,
			Err:  err,
		}
	case "foreign key":
		return rel.ConstraintError{
			Key:  sql.ExtractString(err.Error(), "constraint \"", "\""),
			Type: rel.ForeignKeyConstraint,
			Err:  err,
		}
	case "check":
		return rel.ConstraintError{
			Key:  sql.ExtractString(err.Error(), "constraint \"", "\""),
			Type: rel.CheckConstraint,
			Err:  err,
		}
	default:
		return err
	}
}

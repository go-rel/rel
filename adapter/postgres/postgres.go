// Package postgres wraps postgres (pq) driver as an adapter for REL.
//
// Usage:
//	// open postgres connection.
//	adapter, err := postgres.Open("postgres://postgres@localhost/rel_test?sslmode=disable")
//	if err != nil {
//		panic(err)
//	}
//	defer adapter.Close()
//
//	// initialize REL's repo.
//	repo := rel.New(adapter)
package postgres

import (
	"context"
	db "database/sql"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/sql"
)

// Adapter definition for postgres database.
type Adapter struct {
	*sql.Adapter
}

var _ rel.Adapter = (*Adapter)(nil)

// New is postgres adapter constructor.
func New(database *db.DB) *Adapter {
	return &Adapter{
		Adapter: &sql.Adapter{
			Config: &sql.Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
				ErrorFunc:           errorFunc,
			},
			DB: database,
		},
	}
}

// Open postgres connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var database, err = db.Open("postgres", dsn)
	return New(database), err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(ctx context.Context, query rel.Query, mutates map[string]rel.Mutate) (interface{}, error) {
	var (
		id              int64
		statement, args = sql.NewBuilder(adapter.Config).Returning("id").Insert(query.Table, mutates)
		rows, err       = adapter.query(ctx, statement, args)
	)

	if err == nil && rows.Next() {
		defer rows.Close()
		rows.Scan(&id)
	}

	return id, err
}

// InsertAll inserts multiple records to database and returns its ids.
func (adapter *Adapter) InsertAll(ctx context.Context, query rel.Query, fields []string, bulkMutates []map[string]rel.Mutate) ([]interface{}, error) {
	var (
		ids             []interface{}
		statement, args = sql.NewBuilder(adapter.Config).Returning("id").InsertAll(query.Table, fields, bulkMutates)
		rows, err       = adapter.query(ctx, statement, args)
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

func (adapter *Adapter) query(ctx context.Context, statement string, args []interface{}) (*db.Rows, error) {
	var (
		err  error
		rows *db.Rows
	)

	finish := adapter.Instrument(ctx, "adapter-query", statement)
	if adapter.Tx != nil {
		rows, err = adapter.Tx.QueryContext(ctx, statement, args...)
	} else {
		rows, err = adapter.DB.QueryContext(ctx, statement, args...)
	}
	finish(err)

	return rows, errorFunc(err)
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin(ctx context.Context) (rel.Adapter, error) {
	newAdapter, err := adapter.Adapter.Begin(ctx)

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

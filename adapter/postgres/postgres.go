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
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/sql"
)

// Adapter definition for postgres database.
type Adapter struct {
	*sql.Adapter
}

var (
	_ rel.Adapter = (*Adapter)(nil)

	// Config for postgres adapter.
	Config = sql.Config{
		Placeholder:         "$",
		EscapeChar:          "\"",
		Ordinal:             true,
		InsertDefaultValues: true,
		ErrorFunc:           errorFunc,
		MapColumnFunc:       mapColumnFunc,
	}
)

// New postgres adapter using existing connection.
func New(database *db.DB) *Adapter {
	return &Adapter{
		Adapter: &sql.Adapter{
			Config: Config,
			DB:     database,
		},
	}
}

// Open postgres connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var database, err = db.Open("postgres", dsn)
	return New(database), err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(ctx context.Context, query rel.Query, primaryField string, mutates map[string]rel.Mutate) (interface{}, error) {
	var (
		id              int64
		statement, args = sql.NewBuilder(adapter.Config).Returning(primaryField).Insert(query.Table, mutates)
		rows, err       = adapter.query(ctx, statement, args)
	)

	if err == nil && rows.Next() {
		defer rows.Close()
		rows.Scan(&id)
	}

	return id, err
}

// InsertAll inserts multiple records to database and returns its ids.
func (adapter *Adapter) InsertAll(ctx context.Context, query rel.Query, primaryField string, fields []string, bulkMutates []map[string]rel.Mutate) ([]interface{}, error) {
	var (
		ids             []interface{}
		statement, args = sql.NewBuilder(adapter.Config).Returning(primaryField).InsertAll(query.Table, fields, bulkMutates)
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

	finish := adapter.Instrumenter.Observe(ctx, "adapter-query", statement)
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

func mapColumnFunc(column *rel.Column) (string, int, int) {
	var (
		typ  string
		m, n int
	)

	// postgres specific
	column.Unsigned = false
	if column.Default == "" {
		column.Default = nil
	}

	switch column.Type {
	case rel.ID:
		typ = "SERIAL NOT NULL PRIMARY KEY"
	case rel.DateTime:
		typ = "TIMESTAMPTZ"
		if t, ok := column.Default.(time.Time); ok {
			column.Default = t.Format("2006-01-02 15:04:05")
		}
	case rel.Int, rel.BigInt, rel.Text:
		column.Limit = 0
		typ, m, n = sql.MapColumn(column)
	default:
		typ, m, n = sql.MapColumn(column)
	}

	return typ, m, n
}

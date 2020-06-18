// Package sql is general sql adapter that wraps database/sql.
package sql

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/Fs02/rel"
)

// Config holds configuration for adapter.
type Config struct {
	Placeholder         string
	Ordinal             bool
	InsertDefaultValues bool
	EscapeChar          string
	ErrorFunc           func(error) error
	IncrementFunc       func(Adapter) int
}

// Adapter definition for database database.
type Adapter struct {
	Instrumenter rel.Instrumenter
	Config       *Config
	DB           *sql.DB
	Tx           *sql.Tx
	savepoint    int
}

var _ rel.Adapter = (*Adapter)(nil)

// Close database connection.
func (adapter *Adapter) Close() error {
	return adapter.DB.Close()
}

// Instrumentation set instrumenter for this adapter.
func (adapter *Adapter) Instrumentation(instrumenter rel.Instrumenter) {
	adapter.Instrumenter = instrumenter
}

// Instrument call instrumenter, if no instrumenter is set, this will be a no op.
func (adapter *Adapter) Instrument(ctx context.Context, op string, message string) func(err error) {
	if adapter.Instrumenter != nil {
		return adapter.Instrumenter(ctx, op, message)
	}

	return func(err error) {}
}

// Ping database.
func (adapter *Adapter) Ping(ctx context.Context) error {
	return adapter.DB.PingContext(ctx)
}

// Aggregate record using given query.
func (adapter *Adapter) Aggregate(ctx context.Context, query rel.Query, mode string, field string) (int, error) {
	var (
		err             error
		out             sql.NullInt64
		statement, args = NewBuilder(adapter.Config).Aggregate(query, mode, field)
	)

	finish := adapter.Instrument(ctx, "adapter-aggregate", statement)
	if adapter.Tx != nil {
		err = adapter.Tx.QueryRowContext(ctx, statement, args...).Scan(&out)
	} else {
		err = adapter.DB.QueryRowContext(ctx, statement, args...).Scan(&out)
	}
	finish(err)

	return int(out.Int64), err
}

// Query performs query operation.
func (adapter *Adapter) Query(ctx context.Context, query rel.Query) (rel.Cursor, error) {
	var (
		statement, args = NewBuilder(adapter.Config).Find(query)
	)

	finish := adapter.Instrument(ctx, "adapter-query", statement)
	rows, err := adapter.query(ctx, statement, args)
	finish(err)

	return &Cursor{rows}, adapter.Config.ErrorFunc(err)
}

func (adapter *Adapter) query(ctx context.Context, statement string, args []interface{}) (*sql.Rows, error) {
	if adapter.Tx != nil {
		return adapter.Tx.QueryContext(ctx, statement, args...)
	}

	return adapter.DB.QueryContext(ctx, statement, args...)
}

// Exec performs exec operation.
func (adapter *Adapter) Exec(ctx context.Context, statement string, args []interface{}) (int64, int64, error) {
	finish := adapter.Instrument(ctx, "adapter-exec", statement)
	res, err := adapter.exec(ctx, statement, args)
	finish(err)

	if err != nil {
		return 0, 0, adapter.Config.ErrorFunc(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

func (adapter *Adapter) exec(ctx context.Context, statement string, args []interface{}) (sql.Result, error) {
	if adapter.Tx != nil {
		return adapter.Tx.ExecContext(ctx, statement, args...)
	}

	return adapter.DB.ExecContext(ctx, statement, args...)
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(ctx context.Context, query rel.Query, mutates map[string]rel.Mutate) (interface{}, error) {
	var (
		statement, args = NewBuilder(adapter.Config).Insert(query.Table, mutates)
		id, _, err      = adapter.Exec(ctx, statement, args)
	)

	return id, err
}

// InsertAll inserts all record to database and returns its ids.
func (adapter *Adapter) InsertAll(ctx context.Context, query rel.Query, fields []string, bulkMutates []map[string]rel.Mutate) ([]interface{}, error) {
	statement, args := NewBuilder(adapter.Config).InsertAll(query.Table, fields, bulkMutates)
	id, _, err := adapter.Exec(ctx, statement, args)
	if err != nil {
		return nil, err
	}

	var (
		ids = make([]interface{}, len(bulkMutates))
		inc = 1
	)

	if adapter.Config.IncrementFunc != nil {
		inc = adapter.Config.IncrementFunc(*adapter)
	}

	if inc < 0 {
		id = id + int64((len(bulkMutates)-1)*inc)
		inc *= -1
	}

	for i := range ids {
		ids[i] = id + int64(i*inc)
	}

	return ids, nil
}

// Update updates a record in database.
func (adapter *Adapter) Update(ctx context.Context, query rel.Query, mutates map[string]rel.Mutate) (int, error) {
	var (
		statement, args      = NewBuilder(adapter.Config).Update(query.Table, mutates, query.WhereQuery)
		_, updatedCount, err = adapter.Exec(ctx, statement, args)
	)

	return int(updatedCount), err
}

// Delete deletes all results that match the query.
func (adapter *Adapter) Delete(ctx context.Context, query rel.Query) (int, error) {
	var (
		statement, args      = NewBuilder(adapter.Config).Delete(query.Table, query.WhereQuery)
		_, deletedCount, err = adapter.Exec(ctx, statement, args)
	)

	return int(deletedCount), err
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin(ctx context.Context) (rel.Adapter, error) {
	var (
		tx        *sql.Tx
		savepoint int
		err       error
	)

	finish := adapter.Instrument(ctx, "adapter-begin", "begin transaction")

	if adapter.Tx != nil {
		tx = adapter.Tx
		savepoint = adapter.savepoint + 1
		_, _, err = adapter.Exec(ctx, "SAVEPOINT s"+strconv.Itoa(savepoint)+";", []interface{}{})
	} else {
		tx, err = adapter.DB.BeginTx(ctx, nil)
	}

	finish(err)

	return &Adapter{
		Instrumenter: adapter.Instrumenter,
		Config:       adapter.Config,
		Tx:           tx,
		savepoint:    savepoint,
	}, err
}

// Commit commits current transaction.
func (adapter *Adapter) Commit(ctx context.Context) error {
	var err error

	finish := adapter.Instrument(ctx, "adapter-commit", "commit transaction")

	if adapter.Tx == nil {
		err = errors.New("unable to commit outside transaction")
	} else if adapter.savepoint > 0 {
		_, _, err = adapter.Exec(ctx, "RELEASE SAVEPOINT s"+strconv.Itoa(adapter.savepoint)+";", []interface{}{})
	} else {
		err = adapter.Tx.Commit()
	}

	finish(err)

	return adapter.Config.ErrorFunc(err)
}

// Rollback revert current transaction.
func (adapter *Adapter) Rollback(ctx context.Context) error {
	var err error

	finish := adapter.Instrument(ctx, "adapter-rollback", "rollback transaction")

	if adapter.Tx == nil {
		err = errors.New("unable to rollback outside transaction")
	} else if adapter.savepoint > 0 {
		_, _, err = adapter.Exec(ctx, "ROLLBACK TO SAVEPOINT s"+strconv.Itoa(adapter.savepoint)+";", []interface{}{})
	} else {
		err = adapter.Tx.Rollback()
	}

	finish(err)

	return adapter.Config.ErrorFunc(err)
}

// New initialize adapter without db.
func New(config *Config) *Adapter {
	adapter := &Adapter{
		Config: config,
	}

	return adapter
}

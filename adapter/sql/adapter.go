// Package sql is general sql adapter that wraps database/sql.
package sql

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

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

// Adapter definition for mysql database.
type Adapter struct {
	Config    *Config
	DB        *sql.DB
	Tx        *sql.Tx
	savepoint int
}

var _ rel.Adapter = (*Adapter)(nil)

// Close mysql connection.
func (adapter *Adapter) Close() error {
	return adapter.DB.Close()
}

// Aggregate record using given query.
func (adapter *Adapter) Aggregate(ctx context.Context, query rel.Query, mode string, field string, loggers ...rel.Logger) (int, error) {
	var (
		err             error
		out             sql.NullInt64
		statement, args = NewBuilder(adapter.Config).Aggregate(query, mode, field)
	)

	start := time.Now()
	if adapter.Tx != nil {
		err = adapter.Tx.QueryRowContext(ctx, statement, args...).Scan(&out)
	} else {
		err = adapter.DB.QueryRowContext(ctx, statement, args...).Scan(&out)
	}

	go rel.Log(loggers, statement, time.Since(start), err)

	return int(out.Int64), err
}

// Query performs query operation.
func (adapter *Adapter) Query(ctx context.Context, query rel.Query, loggers ...rel.Logger) (rel.Cursor, error) {
	var (
		rows            *sql.Rows
		err             error
		statement, args = NewBuilder(adapter.Config).Find(query)
	)

	start := time.Now()
	if adapter.Tx != nil {
		rows, err = adapter.Tx.QueryContext(ctx, statement, args...)
	} else {
		rows, err = adapter.DB.QueryContext(ctx, statement, args...)
	}

	go rel.Log(loggers, statement, time.Since(start), err)

	return &Cursor{rows}, adapter.Config.ErrorFunc(err)
}

// Exec performs exec operation.
func (adapter *Adapter) Exec(ctx context.Context, statement string, args []interface{}, loggers ...rel.Logger) (int64, int64, error) {
	var (
		res sql.Result
		err error
	)

	start := time.Now()
	if adapter.Tx != nil {
		res, err = adapter.Tx.ExecContext(ctx, statement, args...)
	} else {
		res, err = adapter.DB.ExecContext(ctx, statement, args...)
	}

	go rel.Log(loggers, statement, time.Since(start), err)

	if err != nil {
		return 0, 0, adapter.Config.ErrorFunc(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(ctx context.Context, query rel.Query, modifies map[string]rel.Modify, loggers ...rel.Logger) (interface{}, error) {
	var (
		statement, args = NewBuilder(adapter.Config).Insert(query.Table, modifies)
		id, _, err      = adapter.Exec(ctx, statement, args, loggers...)
	)

	return id, err
}

// InsertAll inserts all record to database and returns its ids.
func (adapter *Adapter) InsertAll(ctx context.Context, query rel.Query, fields []string, bulkModifies []map[string]rel.Modify, loggers ...rel.Logger) ([]interface{}, error) {
	statement, args := NewBuilder(adapter.Config).InsertAll(query.Table, fields, bulkModifies)
	id, _, err := adapter.Exec(ctx, statement, args, loggers...)
	if err != nil {
		return nil, err
	}

	var (
		ids = make([]interface{}, len(bulkModifies))
		inc = 1
	)

	if adapter.Config.IncrementFunc != nil {
		inc = adapter.Config.IncrementFunc(*adapter)
	}

	if inc < 0 {
		id = id + int64((len(bulkModifies)-1)*inc)
		inc *= -1
	}

	for i := range ids {
		ids[i] = id + int64(i*inc)
	}

	return ids, nil
}

// Update updates a record in database.
func (adapter *Adapter) Update(ctx context.Context, query rel.Query, modifies map[string]rel.Modify, loggers ...rel.Logger) (int, error) {
	var (
		statement, args      = NewBuilder(adapter.Config).Update(query.Table, modifies, query.WhereQuery)
		_, updatedCount, err = adapter.Exec(ctx, statement, args, loggers...)
	)

	return int(updatedCount), err
}

// Delete deletes all results that match the query.
func (adapter *Adapter) Delete(ctx context.Context, query rel.Query, loggers ...rel.Logger) (int, error) {
	var (
		statement, args      = NewBuilder(adapter.Config).Delete(query.Table, query.WhereQuery)
		_, deletedCount, err = adapter.Exec(ctx, statement, args, loggers...)
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

	if adapter.Tx != nil {
		tx = adapter.Tx
		savepoint = adapter.savepoint + 1
		_, _, err = adapter.Exec(ctx, "SAVEPOINT s"+strconv.Itoa(savepoint)+";", []interface{}{})
	} else {
		tx, err = adapter.DB.BeginTx(ctx, nil)
	}

	return &Adapter{
		Config:    adapter.Config,
		Tx:        tx,
		savepoint: savepoint,
	}, err
}

// Commit commits current transaction.
func (adapter *Adapter) Commit(ctx context.Context) error {
	var err error

	if adapter.Tx == nil {
		err = errors.New("unable to commit outside transaction")
	} else if adapter.savepoint > 0 {
		_, _, err = adapter.Exec(ctx, "RELEASE SAVEPOINT s"+strconv.Itoa(adapter.savepoint)+";", []interface{}{})
	} else {
		err = adapter.Tx.Commit()
	}

	return adapter.Config.ErrorFunc(err)
}

// Rollback revert current transaction.
func (adapter *Adapter) Rollback(ctx context.Context) error {
	var err error

	if adapter.Tx == nil {
		err = errors.New("unable to rollback outside transaction")
	} else if adapter.savepoint > 0 {
		_, _, err = adapter.Exec(ctx, "ROLLBACK TO SAVEPOINT s"+strconv.Itoa(adapter.savepoint)+";", []interface{}{})
	} else {
		err = adapter.Tx.Rollback()
	}

	return adapter.Config.ErrorFunc(err)
}

// New initialize adapter without db.
func New(config *Config) *Adapter {
	adapter := &Adapter{
		Config: config,
	}

	return adapter
}

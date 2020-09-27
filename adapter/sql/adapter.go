// Package sql is general sql adapter that wraps database/sql.
package sql

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/go-rel/rel"
)

// Adapter definition for database database.
type Adapter struct {
	Instrumenter rel.Instrumenter
	Config       Config
	DB           *sql.DB
	Tx           *sql.Tx
	savepoint    int
}

var _ rel.Adapter = (*Adapter)(nil)

// Close database connection.
func (a *Adapter) Close() error {
	return a.DB.Close()
}

// Instrumentation set instrumenter for this adapter.
func (a *Adapter) Instrumentation(instrumenter rel.Instrumenter) {
	a.Instrumenter = instrumenter
}

// Ping database.
func (a *Adapter) Ping(ctx context.Context) error {
	return a.DB.PingContext(ctx)
}

// Aggregate record using given query.
func (a *Adapter) Aggregate(ctx context.Context, query rel.Query, mode string, field string) (int, error) {
	var (
		err             error
		out             sql.NullInt64
		statement, args = NewBuilder(a.Config).Aggregate(query, mode, field)
	)

	finish := a.Instrumenter.Observe(ctx, "adapter-aggregate", statement)
	if a.Tx != nil {
		err = a.Tx.QueryRowContext(ctx, statement, args...).Scan(&out)
	} else {
		err = a.DB.QueryRowContext(ctx, statement, args...).Scan(&out)
	}
	finish(err)

	return int(out.Int64), err
}

// Query performs query operation.
func (a *Adapter) Query(ctx context.Context, query rel.Query) (rel.Cursor, error) {
	var (
		statement, args = NewBuilder(a.Config).Find(query)
	)

	finish := a.Instrumenter.Observe(ctx, "adapter-query", statement)
	rows, err := a.query(ctx, statement, args)
	finish(err)

	return &Cursor{rows}, a.Config.ErrorFunc(err)
}

func (a *Adapter) query(ctx context.Context, statement string, args []interface{}) (*sql.Rows, error) {
	if a.Tx != nil {
		return a.Tx.QueryContext(ctx, statement, args...)
	}

	return a.DB.QueryContext(ctx, statement, args...)
}

// Exec performs exec operation.
func (a *Adapter) Exec(ctx context.Context, statement string, args []interface{}) (int64, int64, error) {
	finish := a.Instrumenter.Observe(ctx, "adapter-exec", statement)
	res, err := a.exec(ctx, statement, args)
	finish(err)

	if err != nil {
		return 0, 0, a.Config.ErrorFunc(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

func (a *Adapter) exec(ctx context.Context, statement string, args []interface{}) (sql.Result, error) {
	if a.Tx != nil {
		return a.Tx.ExecContext(ctx, statement, args...)
	}

	return a.DB.ExecContext(ctx, statement, args...)
}

// Insert inserts a record to database and returns its id.
func (a *Adapter) Insert(ctx context.Context, query rel.Query, primaryField string, mutates map[string]rel.Mutate) (interface{}, error) {
	var (
		statement, args = NewBuilder(a.Config).Insert(query.Table, mutates)
		id, _, err      = a.Exec(ctx, statement, args)
	)

	return id, err
}

// InsertAll inserts all record to database and returns its ids.
func (a *Adapter) InsertAll(ctx context.Context, query rel.Query, primaryField string, fields []string, bulkMutates []map[string]rel.Mutate) ([]interface{}, error) {
	statement, args := NewBuilder(a.Config).InsertAll(query.Table, fields, bulkMutates)
	id, _, err := a.Exec(ctx, statement, args)
	if err != nil {
		return nil, err
	}

	var (
		ids = make([]interface{}, len(bulkMutates))
		inc = 1
	)

	if a.Config.IncrementFunc != nil {
		inc = a.Config.IncrementFunc(*a)
	}

	if inc < 0 {
		id = id + int64((len(bulkMutates)-1)*inc)
		inc *= -1
	}

	if primaryField != "" {
		counter := 0
		for i := range ids {
			if mut, ok := bulkMutates[i][primaryField]; ok {
				ids[i] = mut.Value
				id = toInt64(ids[i])
				counter = 1
			} else {
				ids[i] = id + int64(counter*inc)
				counter++
			}
		}
	}

	return ids, nil
}

// Update updates a record in database.
func (a *Adapter) Update(ctx context.Context, query rel.Query, mutates map[string]rel.Mutate) (int, error) {
	var (
		statement, args      = NewBuilder(a.Config).Update(query.Table, mutates, query.WhereQuery)
		_, updatedCount, err = a.Exec(ctx, statement, args)
	)

	return int(updatedCount), err
}

// Delete deletes all results that match the query.
func (a *Adapter) Delete(ctx context.Context, query rel.Query) (int, error) {
	var (
		statement, args      = NewBuilder(a.Config).Delete(query.Table, query.WhereQuery)
		_, deletedCount, err = a.Exec(ctx, statement, args)
	)

	return int(deletedCount), err
}

// Begin begins a new transaction.
func (a *Adapter) Begin(ctx context.Context) (rel.Adapter, error) {
	var (
		tx        *sql.Tx
		savepoint int
		err       error
	)

	finish := a.Instrumenter.Observe(ctx, "adapter-begin", "begin transaction")

	if a.Tx != nil {
		tx = a.Tx
		savepoint = a.savepoint + 1
		_, _, err = a.Exec(ctx, "SAVEPOINT s"+strconv.Itoa(savepoint)+";", []interface{}{})
	} else {
		tx, err = a.DB.BeginTx(ctx, nil)
	}

	finish(err)

	return &Adapter{
		Instrumenter: a.Instrumenter,
		Config:       a.Config,
		Tx:           tx,
		savepoint:    savepoint,
	}, err
}

// Commit commits current transaction.
func (a *Adapter) Commit(ctx context.Context) error {
	var err error

	finish := a.Instrumenter.Observe(ctx, "adapter-commit", "commit transaction")

	if a.Tx == nil {
		err = errors.New("unable to commit outside transaction")
	} else if a.savepoint > 0 {
		_, _, err = a.Exec(ctx, "RELEASE SAVEPOINT s"+strconv.Itoa(a.savepoint)+";", []interface{}{})
	} else {
		err = a.Tx.Commit()
	}

	finish(err)

	return a.Config.ErrorFunc(err)
}

// Rollback revert current transaction.
func (a *Adapter) Rollback(ctx context.Context) error {
	var err error

	finish := a.Instrumenter.Observe(ctx, "adapter-rollback", "rollback transaction")

	if a.Tx == nil {
		err = errors.New("unable to rollback outside transaction")
	} else if a.savepoint > 0 {
		_, _, err = a.Exec(ctx, "ROLLBACK TO SAVEPOINT s"+strconv.Itoa(a.savepoint)+";", []interface{}{})
	} else {
		err = a.Tx.Rollback()
	}

	finish(err)

	return a.Config.ErrorFunc(err)
}

// Apply table.
func (a *Adapter) Apply(ctx context.Context, migration rel.Migration) error {
	var (
		statement string
		builder   = NewBuilder(a.Config)
	)

	switch v := migration.(type) {
	case rel.Table:
		statement = builder.Table(v)
	case rel.Index:
		statement = builder.Index(v)
	case rel.Raw:
		statement = string(v)
	}

	_, _, err := a.Exec(ctx, statement, nil)
	return err
}

// New initialize adapter without db.
func New(config Config) *Adapter {
	adapter := &Adapter{
		Config: config,
	}

	return adapter
}

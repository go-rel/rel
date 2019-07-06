// Package sql is general sql adapter that wraps database/sql.
package sql

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/change"
	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/query"
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

var _ grimoire.Adapter = (*Adapter)(nil)

// Close mysql connection.
func (adapter *Adapter) Close() error {
	return adapter.DB.Close()
}

// All retrieves all record that match the query.
func (adapter *Adapter) All(q query.Query, doc interface{}, loggers ...grimoire.Logger) (int, error) {
	statement, args := NewBuilder(adapter.Config).Find(q)
	count, err := adapter.Query(doc, statement, args, loggers...)
	return int(count), err
}

// Aggregate record using given query.
func (adapter *Adapter) Aggregate(q query.Query, doc interface{}, todo1 string, todo2 string, loggers ...grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Config).Aggregate(q)
	_, err := adapter.Query(doc, statement, args, loggers...)
	return err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(q query.Query, changes change.Changes, loggers ...grimoire.Logger) (interface{}, error) {
	statement, args := NewBuilder(adapter.Config).Insert(q.Collection, changes)
	id, _, err := adapter.Exec(statement, args, loggers...)
	return id, err
}

// InsertAll inserts all record to database and returns its ids.
func (adapter *Adapter) InsertAll(q query.Query, fields []string, allchanges []change.Changes, loggers ...grimoire.Logger) ([]interface{}, error) {
	statement, args := NewBuilder(adapter.Config).InsertAll(q.Collection, fields, allchanges)
	id, _, err := adapter.Exec(statement, args, loggers...)
	if err != nil {
		return nil, err
	}

	var (
		ids = []interface{}{id}
		inc = 1
	)

	if adapter.Config.IncrementFunc != nil {
		inc = adapter.Config.IncrementFunc(*adapter)
	}

	for i := 1; i < len(allchanges); i++ {
		ids = append(ids, id+int64(inc*i))
	}

	return ids, nil
}

// Update updates a record in database.
func (adapter *Adapter) Update(q query.Query, changes change.Changes, loggers ...grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Config).Update(q.Collection, changes, q.WhereClause)
	_, _, err := adapter.Exec(statement, args, loggers...)
	return err
}

// Delete deletes all results that match the query.
func (adapter *Adapter) Delete(q query.Query, loggers ...grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Config).Delete(q.Collection, q.WhereClause)
	_, _, err := adapter.Exec(statement, args, loggers...)
	return err
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin() (grimoire.Adapter, error) {
	var (
		tx        *sql.Tx
		savepoint int
		err       error
	)

	if adapter.Tx != nil {
		tx = adapter.Tx
		savepoint = adapter.savepoint + 1
		_, _, err = adapter.Exec("SAVEPOINT s"+strconv.Itoa(savepoint)+";", []interface{}{})
	} else {
		tx, err = adapter.DB.Begin()
	}

	return &Adapter{
		Config:    adapter.Config,
		Tx:        tx,
		savepoint: savepoint,
	}, err
}

// Commit commits current transaction.
func (adapter *Adapter) Commit() error {
	var err error

	if adapter.Tx == nil {
		err = errors.NewUnexpected("unable to commit outside transaction")
	} else if adapter.savepoint > 0 {
		_, _, err = adapter.Exec("RELEASE SAVEPOINT s"+strconv.Itoa(adapter.savepoint)+";", []interface{}{})
	} else {
		err = adapter.Tx.Commit()
	}

	return adapter.Config.ErrorFunc(err)
}

// Rollback revert current transaction.
func (adapter *Adapter) Rollback() error {
	var err error

	if adapter.Tx == nil {
		err = errors.NewUnexpected("unable to rollback outside transaction")
	} else if adapter.savepoint > 0 {
		_, _, err = adapter.Exec("ROLLBACK TO SAVEPOINT s"+strconv.Itoa(adapter.savepoint)+";", []interface{}{})
	} else {
		err = adapter.Tx.Rollback()
	}

	return adapter.Config.ErrorFunc(err)
}

// Query performs query operation.
func (adapter *Adapter) Query(out interface{}, statement string, args []interface{}, loggers ...grimoire.Logger) (int64, error) {
	var (
		rows *sql.Rows
		err  error
	)

	start := time.Now()
	if adapter.Tx != nil {
		rows, err = adapter.Tx.Query(statement, args...)
	} else {
		rows, err = adapter.DB.Query(statement, args...)
	}

	go grimoire.Log(loggers, statement, time.Since(start), err)

	if err != nil {
		return 0, adapter.Config.ErrorFunc(err)
	}

	defer rows.Close()
	count, err := Scan(out, rows)
	return count, adapter.Config.ErrorFunc(err)
}

// Exec performs exec operation.
func (adapter *Adapter) Exec(statement string, args []interface{}, loggers ...grimoire.Logger) (int64, int64, error) {
	var (
		res sql.Result
		err error
	)

	start := time.Now()
	if adapter.Tx != nil {
		res, err = adapter.Tx.Exec(statement, args...)
	} else {
		res, err = adapter.DB.Exec(statement, args...)
	}

	go grimoire.Log(loggers, statement, time.Since(start), err)

	if err != nil {
		return 0, 0, adapter.Config.ErrorFunc(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

// New initialize adapter without db.
func New(config *Config) *Adapter {
	adapter := &Adapter{
		Config: config,
	}

	return adapter
}

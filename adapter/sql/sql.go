// Package sql is general sql adapter that wraps database/sql.
package sql

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/errors"
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
func (adapter *Adapter) All(query grimoire.Query, doc interface{}, loggers ...grimoire.Logger) (int, error) {
	statement, args := NewBuilder(adapter.Config).Find(query)
	count, err := adapter.Query(doc, statement, args, loggers...)
	return int(count), err
}

// Aggregate record using given query.
func (adapter *Adapter) Aggregate(query grimoire.Query, doc interface{}, loggers ...grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Config).Aggregate(query)
	_, err := adapter.Query(doc, statement, args, loggers...)
	return err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(query grimoire.Query, changes map[string]interface{}, loggers ...grimoire.Logger) (interface{}, error) {
	statement, args := NewBuilder(adapter.Config).Insert(query.Collection, changes)
	id, _, err := adapter.Exec(statement, args, loggers...)
	return id, err
}

// InsertAll inserts all record to database and returns its ids.
func (adapter *Adapter) InsertAll(query grimoire.Query, fields []string, allchanges []map[string]interface{}, loggers ...grimoire.Logger) ([]interface{}, error) {
	statement, args := NewBuilder(adapter.Config).InsertAll(query.Collection, fields, allchanges)
	id, _, err := adapter.Exec(statement, args, loggers...)
	if err != nil {
		return nil, err
	}

	ids := []interface{}{id}
	inc := 1

	if adapter.Config.IncrementFunc != nil {
		inc = adapter.Config.IncrementFunc(*adapter)
	}

	for i := 1; i < len(allchanges); i++ {
		ids = append(ids, id+int64(inc*i))
	}

	return ids, nil
}

// Update updates a record in database.
func (adapter *Adapter) Update(query grimoire.Query, changes map[string]interface{}, loggers ...grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Config).Update(query.Collection, changes, query.Condition)
	_, _, err := adapter.Exec(statement, args, loggers...)
	return err
}

// Delete deletes all results that match the query.
func (adapter *Adapter) Delete(query grimoire.Query, loggers ...grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Config).Delete(query.Collection, query.Condition)
	_, _, err := adapter.Exec(statement, args, loggers...)
	return err
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin() (grimoire.Adapter, error) {
	var tx *sql.Tx
	var savepoint int
	var err error

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
	var rows *sql.Rows
	var err error

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
	var res sql.Result
	var err error

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

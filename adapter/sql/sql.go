// Package sql is general sql adapter that wraps database/sql.
package sql

import (
	"database/sql"
	"time"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/errors"
)

// Adapter definition for mysql database.
type Adapter struct {
	Placeholder         string
	Ordinal             bool
	InsertDefaultValues bool
	ErrorFunc           func(error) error
	IncrementFunc       func(Adapter) int
	DB                  *sql.DB
	Tx                  *sql.Tx
}

var _ grimoire.Adapter = (*Adapter)(nil)

// Close mysql connection.
func (adapter *Adapter) Close() error {
	return adapter.DB.Close()
}

// Count retrieves count of record that match the query.
func (adapter *Adapter) Count(query grimoire.Query, loggers ...grimoire.Logger) (int, error) {
	var doc struct {
		Count int
	}

	query.Fields = []string{"COUNT(*) AS count"}
	statement, args := NewBuilder(adapter.Placeholder, adapter.Ordinal, adapter.InsertDefaultValues).Find(query)
	_, err := adapter.Query(&doc, statement, args, loggers...)
	return doc.Count, err
}

// All retrieves all record that match the query.
func (adapter *Adapter) All(query grimoire.Query, doc interface{}, loggers ...grimoire.Logger) (int, error) {
	statement, args := NewBuilder(adapter.Placeholder, adapter.Ordinal, adapter.InsertDefaultValues).Find(query)
	count, err := adapter.Query(doc, statement, args, loggers...)
	return int(count), err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(query grimoire.Query, changes map[string]interface{}, loggers ...grimoire.Logger) (interface{}, error) {
	statement, args := NewBuilder(adapter.Placeholder, adapter.Ordinal, adapter.InsertDefaultValues).Insert(query.Collection, changes)
	id, _, err := adapter.Exec(statement, args, loggers...)
	return id, err
}

// InsertAll inserts all record to database and returns its ids.
func (adapter *Adapter) InsertAll(query grimoire.Query, fields []string, allchanges []map[string]interface{}, loggers ...grimoire.Logger) ([]interface{}, error) {
	statement, args := NewBuilder(adapter.Placeholder, adapter.Ordinal, adapter.InsertDefaultValues).InsertAll(query.Collection, fields, allchanges)
	id, _, err := adapter.Exec(statement, args, loggers...)
	if err != nil {
		return nil, err
	}

	ids := []interface{}{id}
	inc := 1

	if adapter.IncrementFunc != nil {
		inc = adapter.IncrementFunc(*adapter)
	}

	for i := 1; i < len(allchanges); i++ {
		ids = append(ids, id+int64(inc*i))
	}

	return ids, nil
}

// Update updates a record in database.
func (adapter *Adapter) Update(query grimoire.Query, changes map[string]interface{}, loggers ...grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Placeholder, adapter.Ordinal, adapter.InsertDefaultValues).Update(query.Collection, changes, query.Condition)
	_, _, err := adapter.Exec(statement, args, loggers...)
	return err
}

// Delete deletes all results that match the query.
func (adapter *Adapter) Delete(query grimoire.Query, loggers ...grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Placeholder, adapter.Ordinal, adapter.InsertDefaultValues).Delete(query.Collection, query.Condition)
	_, _, err := adapter.Exec(statement, args, loggers...)
	return err
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin() (grimoire.Adapter, error) {
	Tx, err := adapter.DB.Begin()

	return &Adapter{
		Placeholder:   adapter.Placeholder,
		Ordinal:       adapter.Ordinal,
		IncrementFunc: adapter.IncrementFunc,
		ErrorFunc:     adapter.ErrorFunc,
		Tx:            Tx,
	}, err
}

// Commit commits current transaction.
func (adapter *Adapter) Commit() error {
	if adapter.Tx == nil {
		return errors.NewUnexpected("not in transaction")
	}

	err := adapter.Tx.Commit()
	return adapter.ErrorFunc(err)
}

// Rollback revert current transaction.
func (adapter *Adapter) Rollback() error {
	if adapter.Tx == nil {
		return errors.NewUnexpected("not in transaction")
	}

	err := adapter.Tx.Rollback()
	return adapter.ErrorFunc(err)
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
		return 0, adapter.ErrorFunc(err)
	}

	defer rows.Close()
	count, err := Scan(out, rows)
	return count, adapter.ErrorFunc(err)
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
		return 0, 0, adapter.ErrorFunc(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

// New initialize adapter without db.
func New(errfn func(error) error, incfn func(Adapter) int, configs ...Config) *Adapter {
	adapter := &Adapter{
		Placeholder:   "?",
		ErrorFunc:     errfn,
		IncrementFunc: incfn,
	}

	for _, config := range configs {
		config(adapter)
	}

	return adapter
}

// Config for adapter initialization.
type Config func(*Adapter)

// Placeholder set placeholder string, default "?".
func Placeholder(placeholder string) Config {
	return func(adapter *Adapter) {
		adapter.Placeholder = placeholder
	}
}

// Ordinal set the placeholder to use ordinal format.
func Ordinal(ordinal bool) Config {
	return func(adapter *Adapter) {
		adapter.Ordinal = ordinal
	}
}

// InsertDefaultValues enable "insert into collection default values" when insertinv record without any value.
func InsertDefaultValues(flag bool) Config {
	return func(adapter *Adapter) {
		adapter.InsertDefaultValues = flag
	}
}

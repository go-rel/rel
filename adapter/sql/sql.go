package sql

import (
	"database/sql"
	"time"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/errors"
)

// Adapter definition for mysql database.
type Adapter struct {
	Placeholder   string
	IsOrdinal     bool
	IncrementFunc func(Adapter) int
	ErrorFunc     func(error) error
	DB            *sql.DB
	Tx            *sql.Tx
}

var _ grimoire.Adapter = (*Adapter)(nil)

// Close mysql connection.
func (adapter *Adapter) Close() error {
	return adapter.DB.Close()
}

// All retrieves all record that match the query.
func (adapter *Adapter) All(query grimoire.Query, doc interface{}, logger grimoire.Logger) (int, error) {
	statement, args := NewBuilder(adapter.Placeholder, adapter.IsOrdinal).Find(query)
	count, err := adapter.Query(doc, statement, args, logger)
	return int(count), err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(query grimoire.Query, changes map[string]interface{}, logger grimoire.Logger) (interface{}, error) {
	statement, args := NewBuilder(adapter.Placeholder, adapter.IsOrdinal).Insert(query.Collection, changes)
	id, _, err := adapter.Exec(statement, args, logger)
	return id, err
}

// InsertAll inserts all record to database and returns its ids.
func (adapter *Adapter) InsertAll(query grimoire.Query, fields []string, allchanges []map[string]interface{}, logger grimoire.Logger) ([]interface{}, error) {
	statement, args := NewBuilder(adapter.Placeholder, adapter.IsOrdinal).InsertAll(query.Collection, fields, allchanges)
	id, _, err := adapter.Exec(statement, args, logger)
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
func (adapter *Adapter) Update(query grimoire.Query, changes map[string]interface{}, logger grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Placeholder, adapter.IsOrdinal).Update(query.Collection, changes, query.Condition)
	_, _, err := adapter.Exec(statement, args, logger)
	return err
}

// Delete deletes all results that match the query.
func (adapter *Adapter) Delete(query grimoire.Query, logger grimoire.Logger) error {
	statement, args := NewBuilder(adapter.Placeholder, adapter.IsOrdinal).Delete(query.Collection, query.Condition)
	_, _, err := adapter.Exec(statement, args, logger)
	return err
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin() (grimoire.Adapter, error) {
	Tx, err := adapter.DB.Begin()

	return &Adapter{
		Placeholder:   adapter.Placeholder,
		IsOrdinal:     adapter.IsOrdinal,
		IncrementFunc: adapter.IncrementFunc,
		ErrorFunc:     adapter.ErrorFunc,
		Tx:            Tx,
	}, err
}

// Commit commits current transaction.
func (adapter *Adapter) Commit() error {
	if adapter.Tx == nil {
		return errors.UnexpectedError("not in transaction")
	}

	err := adapter.Tx.Commit()
	return adapter.ErrorFunc(err)
}

// Rollback revert current transaction.
func (adapter *Adapter) Rollback() error {
	if adapter.Tx == nil {
		return errors.UnexpectedError("not in transaction")
	}

	err := adapter.Tx.Rollback()
	return adapter.ErrorFunc(err)
}

// Query performs query operation.
func (adapter *Adapter) Query(out interface{}, statement string, args []interface{}, logger grimoire.Logger) (int64, error) {
	var rows *sql.Rows
	var err error

	start := time.Now()
	if adapter.Tx != nil {
		rows, err = adapter.Tx.Query(statement, args...)
	} else {
		rows, err = adapter.DB.Query(statement, args...)
	}
	logger(statement, time.Since(start), err)

	if err != nil {
		return 0, adapter.ErrorFunc(err)
	}

	defer rows.Close()
	count, err := Scan(out, rows)
	return count, adapter.ErrorFunc(err)
}

// Exec performs exec operation.
func (adapter *Adapter) Exec(statement string, args []interface{}, logger grimoire.Logger) (int64, int64, error) {
	var res sql.Result
	var err error

	start := time.Now()
	if adapter.Tx != nil {
		res, err = adapter.Tx.Exec(statement, args...)
	} else {
		res, err = adapter.DB.Exec(statement, args...)
	}
	logger(statement, time.Since(start), err)

	if err != nil {
		return 0, 0, adapter.ErrorFunc(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

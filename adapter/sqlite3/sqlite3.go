package sqlite3

import (
	"context"
	"database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/errors"
	_ "github.com/mattn/go-sqlite3"
)

// Adapter definition for mysql database.
type Adapter struct {
	db *sql.DB
	tx *sql.Tx
}

var _ grimoire.Adapter = (*Adapter)(nil)

// Open mysql connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var err error
	adapter := &Adapter{}
	adapter.db, err = sql.Open("sqlite3", dsn)
	return adapter, err
}

// Close mysql connection.
func (adapter *Adapter) Close() error {
	return adapter.db.Close()
}

// All retrieves all record that match the query.
func (adapter *Adapter) All(query grimoire.Query, doc interface{}) (int, error) {
	statement, args := sqlutil.NewBuilder("?", false).Find(query)
	count, err := adapter.Query(doc, statement, args)
	return int(count), err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(query grimoire.Query, changes map[string]interface{}) (interface{}, error) {
	statement, args := sqlutil.NewBuilder("?", false).Insert(query.Collection, changes)
	id, _, err := adapter.Exec(statement, args)
	return id, err
}

// InsertAll inserts all record to database and returns its ids.
func (adapter *Adapter) InsertAll(query grimoire.Query, fields []string, allchanges []map[string]interface{}) ([]interface{}, error) {
	statement, args := sqlutil.NewBuilder("?", false).InsertAll(query.Collection, fields, allchanges)
	id, _, err := adapter.Exec(statement, args)
	if err != nil {
		return nil, err
	}

	inc := adapter.getIncrement()
	ids := []interface{}{id}

	for i := 1; i < len(allchanges); i++ {
		ids = append(ids, id+int64(inc*i))
	}

	return ids, nil
}

// Update updates a record in database.
func (adapter *Adapter) Update(query grimoire.Query, changes map[string]interface{}) error {
	statement, args := sqlutil.NewBuilder("?", false).Update(query.Collection, changes, query.Condition)
	_, _, err := adapter.Exec(statement, args)
	return err
}

// Delete deletes all results that match the query.
func (adapter *Adapter) Delete(query grimoire.Query) error {
	statement, args := sqlutil.NewBuilder("?", false).Delete(query.Collection, query.Condition)
	_, _, err := adapter.Exec(statement, args)
	return err
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin() (grimoire.Adapter, error) {
	tx, err := adapter.db.BeginTx(context.Background(), nil)
	return &Adapter{tx: tx}, err
}

// Commit commits current transaction.
func (adapter *Adapter) Commit() error {
	if adapter.tx == nil {
		return errors.UnexpectedError("not in transaction")
	}

	err := adapter.tx.Commit()
	return adapter.Error(err)
}

// Rollback revert current transaction.
func (adapter *Adapter) Rollback() error {
	if adapter.tx == nil {
		return errors.UnexpectedError("not in transaction")
	}

	err := adapter.tx.Rollback()
	return adapter.Error(err)
}

// Query performs query operation.
func (adapter *Adapter) Query(out interface{}, statement string, args []interface{}) (int64, error) {
	var rows *sql.Rows
	var err error

	if adapter.tx != nil {
		rows, err = adapter.tx.Query(statement, args...)
	} else {
		rows, err = adapter.db.Query(statement, args...)
	}

	if err != nil {
		return 0, adapter.Error(err)
	}

	defer rows.Close()
	count, err := sqlutil.Scan(out, rows)
	return count, adapter.Error(err)
}

// Exec performs exec operation.
func (adapter *Adapter) Exec(statement string, args []interface{}) (int64, int64, error) {
	var res sql.Result
	var err error

	if adapter.tx != nil {
		res, err = adapter.tx.Exec(statement, args...)
	} else {
		res, err = adapter.db.Exec(statement, args...)
	}

	if err != nil {
		return 0, 0, adapter.Error(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

func (adapter *Adapter) getIncrement() int {
	return 1
}

// Error transform adapter error to grimoire error.
func (adapter *Adapter) Error(err error) error {
	if err == nil {
		return nil
		// } else if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		// 	return errors.DuplicateError(e.Message, "")
	}

	return err
}

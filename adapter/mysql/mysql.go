package mysql

import (
	"context"
	"database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/errors"
	"github.com/go-sql-driver/mysql"
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
	adapter.db, err = sql.Open("mysql", dsn)
	return adapter, err
}

// Close mysql connection.
func (adapter *Adapter) Close() error {
	return adapter.db.Close()
}

// Find generates sql query and arguments query operation.
func (adapter *Adapter) Find(query grimoire.Query) (string, []interface{}) {
	return sqlutil.NewBuilder("?", false).Find(query)
}

// Insert generates sql query and arguments for insert operation.
func (adapter *Adapter) Insert(query grimoire.Query, changes map[string]interface{}) (string, []interface{}) {
	return sqlutil.NewBuilder("?", false).Insert(query.Collection, changes)
}

// Update generates sql query and arguments for update operation.
func (adapter *Adapter) Update(query grimoire.Query, changes map[string]interface{}) (string, []interface{}) {
	return sqlutil.NewBuilder("?", false).Update(query.Collection, changes, query.Condition)
}

// Delete generates sql query and argumetns for delete opration.
func (adapter *Adapter) Delete(query grimoire.Query) (string, []interface{}) {
	return sqlutil.NewBuilder("?", false).Delete(query.Collection, query.Condition)
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin() error {
	tx, err := adapter.db.BeginTx(context.Background(), nil)
	adapter.tx = tx
	return err
}

// Commit commits current transaction.
func (adapter *Adapter) Commit() error {
	if adapter.tx == nil {
		return errors.UnexpectedError("not in transaction")
	}

	err := adapter.tx.Commit()
	adapter.tx = nil
	return adapter.Error(err)
}

// Rollback revert current transaction.
func (adapter *Adapter) Rollback() error {
	if adapter.tx == nil {
		return errors.UnexpectedError("not in transaction")
	}

	err := adapter.tx.Rollback()
	adapter.tx = nil
	return adapter.Error(err)
}

// Query performs query operation.
func (adapter *Adapter) Query(out interface{}, qs string, args []interface{}) (int64, error) {
	var rows *sql.Rows
	var err error

	if adapter.tx != nil {
		rows, err = adapter.tx.Query(qs, args...)
	} else {
		rows, err = adapter.db.Query(qs, args...)
	}

	if err != nil {
		return 0, adapter.Error(err)
	}

	defer rows.Close()
	count, err := sqlutil.Scan(out, rows)
	return count, adapter.Error(err)
}

// Exec performs exec operation.
func (adapter *Adapter) Exec(qs string, args []interface{}) (int64, int64, error) {
	var res sql.Result
	var err error

	if adapter.tx != nil {
		res, err = adapter.tx.Exec(qs, args...)
	} else {
		res, err = adapter.db.Exec(qs, args...)
	}

	if err != nil {
		return 0, 0, adapter.Error(err)
	}

	lastID, _ := res.LastInsertId()
	rowCount, _ := res.RowsAffected()

	return lastID, rowCount, nil
}

// Error transform adapter error to grimoire error.
func (adapter *Adapter) Error(err error) error {
	if err == nil {
		return nil
	} else if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return errors.DuplicateError(e.Message, "")
	}

	return err
}

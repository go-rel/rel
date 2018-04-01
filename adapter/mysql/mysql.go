package mysql

import (
	"context"
	"database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/errors"
	"github.com/go-sql-driver/mysql"
)

type Adapter struct {
	db *sql.DB
	tx *sql.Tx
}

var _ grimoire.Adapter = (*Adapter)(nil)

func Open(dsn string) (*Adapter, error) {
	var err error
	adapter := &Adapter{}
	adapter.db, err = sql.Open("mysql", dsn)
	return adapter, err
}

func (adapter *Adapter) Close() error {
	return adapter.db.Close()
}

func (adapter *Adapter) Find(query grimoire.Query) (string, []interface{}) {
	return sqlutil.NewBuilder("?", false).Find(query)
}

func (adapter *Adapter) Insert(query grimoire.Query, changes map[string]interface{}) (string, []interface{}) {
	return sqlutil.NewBuilder("?", false).Insert(query.Collection, changes)
}

func (adapter *Adapter) Update(query grimoire.Query, changes map[string]interface{}) (string, []interface{}) {
	return sqlutil.NewBuilder("?", false).Update(query.Collection, changes, query.Condition)
}

func (adapter *Adapter) Delete(query grimoire.Query) (string, []interface{}) {
	return sqlutil.NewBuilder("?", false).Delete(query.Collection, query.Condition)
}

func (adapter *Adapter) Begin() error {
	tx, err := adapter.db.BeginTx(context.Background(), nil)
	adapter.tx = tx
	return err
}

func (adapter *Adapter) Commit() error {
	if adapter.tx == nil {
		return errors.UnexpectedError("not in transaction")
	}

	err := adapter.tx.Commit()
	adapter.tx = nil
	return adapter.Error(err)
}

func (adapter *Adapter) Rollback() error {
	if adapter.tx == nil {
		return errors.UnexpectedError("not in transaction")
	}

	err := adapter.tx.Rollback()
	adapter.tx = nil
	return adapter.Error(err)
}

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

func (adapter *Adapter) Error(err error) error {
	if err == nil {
		return nil
	} else if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return errors.DuplicateError(e.Message, "")
	}

	return err
}

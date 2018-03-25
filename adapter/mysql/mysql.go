package mysql

import (
	"context"
	"database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	_ "github.com/go-sql-driver/mysql"
)

type Adapter struct {
	db      *sql.DB
	tx      *sql.Tx
	builder sqlutil.Builder
}

var _ grimoire.Adapter = (*Adapter)(nil)

func (adapter *Adapter) Open(dsn string) error {
	var err error
	adapter.db, err = sql.Open("mysql", dsn)
	return err
}

func (adapter *Adapter) Close() error {
	return adapter.db.Close()
}

func (adapter Adapter) Find(query grimoire.Query) (string, []interface{}) {
	return adapter.builder.Find(query)
}

func (adapter Adapter) Insert(query grimoire.Query, ch changeset.Changeset) (string, []interface{}) {
	return adapter.builder.Insert(query.Collection, ch.Changes())
}

func (adapter Adapter) Update(query grimoire.Query, ch changeset.Changeset) (string, []interface{}) {
	return adapter.builder.Update(query.Collection, ch.Changes(), query.Condition)
}

func (adapter Adapter) Delete(query grimoire.Query) (string, []interface{}) {
	return adapter.builder.Delete(query.Collection, query.Condition)
}

func (adapter Adapter) Begin() error {
	tx, err := adapter.db.BeginTx(context.Background(), nil)
	adapter.tx = tx
	return err
}

func (adapter Adapter) Commit() error {
	if adapter.tx == nil {
		return nil
	}

	err := adapter.tx.Commit()
	adapter.tx = nil
	return err
}

func (adapter Adapter) Rollback() error {
	if adapter.tx == nil {
		return nil
	}

	err := adapter.tx.Rollback()
	adapter.tx = nil
	return err
}

func (adapter Adapter) Query(out interface{}, qs string, args []interface{}) errors.Error {
	var rows *sql.Rows
	var err error

	if adapter.tx != nil {
		rows, err = adapter.tx.Query(qs, args...)
	} else {
		rows, err = adapter.db.Query(qs, args...)
	}

	if err != nil {
		return errors.UnexpectedError(err.Error())
	}

	defer rows.Close()
	err = sqlutil.Scan(out, rows)
	if err == sql.ErrNoRows {
		return errors.NotFoundError(err.Error())
	}

	return errors.UnexpectedError(err.Error())
}

func (adapter Adapter) Exec(qs string, args []interface{}) (int64, int64, error) {
	var res sql.Result
	var err error

	if adapter.tx != nil {
		res, err = adapter.tx.Exec(qs, args...)
	} else {
		res, err = adapter.db.Exec(qs, args...)
	}

	if err != nil {
		return 0, 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	return lastId, rowCount, nil
}

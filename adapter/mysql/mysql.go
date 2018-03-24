package mysql

import (
	"database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	_ "github.com/go-sql-driver/mysql"
)

type Adapter struct {
	db *sql.DB
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

func (adapter Adapter) Query(out interface{}, qs string, args []interface{}) error {
	rows, err := adapter.db.Query(qs, args...)
	if err != nil {
		println(err.Error())
		return err
	}

	defer rows.Close()
	return sqlutil.Scan(out, rows)
}

func (adapter Adapter) Exec(qs string, args []interface{}) (int64, int64, error) {
	res, err := adapter.db.Exec(qs, args...)
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

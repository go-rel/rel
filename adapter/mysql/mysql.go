package mysql

import (
	"database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/query"
	_ "github.com/go-sql-driver/mysql"
)

type Adapter struct {
	db *sql.DB
}

func (adapter *Adapter) Open(dsn string) error {
	var err error
	adapter.db, err = sql.Open("mysql", dsn)
	return err
}

func (adapter *Adapter) Close() error {
	return adapter.db.Close()
}

func (adapter Adapter) All(q query.Query) (string, []interface{}) {
	return sqlutil.Builder{}.All(q)
}

func (adapter Adapter) Insert(ch *grimoire.Changeset) (string, []interface{}) {
	return sqlutil.Builder{}.Insert(ch.Collection, ch.Changes)
}

func (adapter Adapter) Update(ch *grimoire.Changeset, cond query.Condition) (string, []interface{}) {
	return sqlutil.Builder{}.Update(ch.Collection, ch.Changes, cond)
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

package mysql

import (
	"bytes"
	"database/sql"
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/query"
	_ "github.com/go-sql-driver/mysql"
	"strings"
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
	builder := sqlutil.Builder{}

	var buffer bytes.Buffer
	var args []interface{}

	if s := builder.Select(q.AsDistinct, q.Fields...); s != "" {
		buffer.WriteString(s)
	}

	if s := builder.From(q.Collection); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s, arg := builder.Join(q.JoinClause...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s, arg := builder.Where(q.Condition); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s := builder.GroupBy(q.GroupFields...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s, arg := builder.Having(q.HavingCondition); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s := builder.OrderBy(q.OrderClause...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s := builder.Offset(q.OffsetResult); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s := builder.Limit(q.LimitResult); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

func (adapter Adapter) Insert(ch *grimoire.Changeset) (string, []interface{}) {
	length := len(ch.Changes)

	var buffer bytes.Buffer
	var args = make([]interface{}, length)

	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(ch.Collection)
	buffer.WriteString(" (")

	curr := 0
	for field, value := range ch.Changes {
		if curr < length-1 {
			buffer.WriteString(",")
		}
		buffer.WriteString(field)
		args = append(args, value)

		curr++
	}
	buffer.WriteString(") VALUES ")
	buffer.WriteString("(?")
	buffer.WriteString(strings.Repeat(",?", length))
	buffer.WriteString(");")

	return buffer.String(), args
}

func (adapter Adapter) Update(ch *grimoire.Changeset, cond query.Condition) (string, []interface{}) {
	builder := sqlutil.Builder{}
	length := len(ch.Changes)

	var buffer bytes.Buffer
	var args = make([]interface{}, length)

	buffer.WriteString("UPDATE ")
	buffer.WriteString(ch.Collection)
	buffer.WriteString(" SET ")

	curr := 0
	for field, value := range ch.Changes {
		if curr < length-1 {
			buffer.WriteString(",")
		}
		buffer.WriteString(field)
		buffer.WriteString("=?")
		args = append(args, value)

		curr++
	}

	if s, arg := builder.Where(cond); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	return buffer.String(), args
}

func (adapter Adapter) Delete(query.Condition) (string, []interface{}) {
	builder := sqlutil.Builder{}
	length := len(ch.Changes)

	var buffer bytes.Buffer
	var args = make([]interface{}, length)

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(ch.Collection)
	buffer.WriteString(" ")

	if s, arg := builder.Where(cond); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	return buffer.String(), args
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

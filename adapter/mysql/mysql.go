package mysql

import (
	"database/sql"
	"fmt"
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

func (adapter *Adapter) Close(string) error {
	return adapter.db.Close()
}

func (adapter Adapter) All(q query.Query) (string, []interface{}) {
	builder := sqlutil.Builder{}

	var str string
	var args []interface{}

	if s := builder.Select(q.AsDistinct, q.Fields...); s != "" {
		str += s
	}

	if s := builder.From(q.Collection); s != "" {
		str += " " + s
	}

	if s, arg := builder.Join(q.JoinClause...); s != "" {
		str += " " + s
		args = append(args, arg...)
	}

	if s, arg := builder.Where(q.Condition); s != "" {
		str += " " + s
		args = append(args, arg...)
	}

	if s := builder.GroupBy(q.GroupFields...); s != "" {
		str += " " + s
	}

	if s, arg := builder.Having(q.HavingCondition); s != "" {
		str += " " + s
		args = append(args, arg...)
	}

	if s := builder.OrderBy(q.OrderClause...); s != "" {
		str += " " + s
	}

	if s := builder.Offset(q.OffsetResult); s != "" {
		str += " " + s
	}

	if s := builder.Limit(q.LimitResult); s != "" {
		str += " " + s
	}

	return str + ";", args
}

func (adapter Adapter) Query(qs string, args []interface{}) ([]interface{}, error) {
	rows, err := adapter.db.Query(qs, args...)
	if err != nil {
		println(err.Error())
	}

	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		println(err.Error())
	}

	fmt.Printf("%+v\n", cols)

	return nil, nil
}

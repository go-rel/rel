package mysql

import (
	"github.com/Fs02/grimoire/adapter/sql"
	"github.com/Fs02/grimoire/query"
)

type Adapter struct{}

func (a Adapter) All(q query.Query) (string, []interface{}) {
	builder := sql.Builder{}

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

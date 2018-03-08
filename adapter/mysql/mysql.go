package sql

import (
	"github.com/Fs02/grimoire/adapter/sql"
	"github.com/Fs02/grimoire/query"
	"strings"
)

type Sql struct{}

func (s Sql) All(q query.Query) (string, []interface{}) {
	builder := sql.Builder{}

	selects := builder.Select(q.AsDistinct, q.Fields...)
	from := builder.From(q.Collection)
	join, joinArgs := builder.Join(q.JoinClause...)
	where, whereArgs := builder.Where(q.Condition)
	group := builder.GroupBy(q.GroupFields...)
	having, havingArgs := builder.Having(q.HavingCondition)
	order := builder.OrderBy(q.OrderClause...)
	offset := builder.Offset(q.OffsetResult)
	limit := builder.Limit(q.LimitResult)

	args := append(joinArgs, whereArgs...)
	args = append(args, havingArgs...)

	return strings.Join([]string{selects, from, join, where, group, having, order, offset, limit}, " "), args
}

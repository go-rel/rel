package sql

import (
	"github.com/Fs02/grimoire/query"
	"strings"
)

type QueryBuilder struct{}

func (q QueryBuilder) Select(fields ...string) string {
	return "SELECT " + strings.Join(fields, ", ")
}

func (q QueryBuilder) From(collection string) string {
	return "FROM " + collection
}

func (q QueryBuilder) Join(join []query.JoinQuery) string {
	return ""
}

func (q QueryBuilder) Where(condition query.Condition) string {
	return ""
}

func (q QueryBuilder) GroupBy(fields ...string) string {
	return "GROUP BY " + strings.Join(fields, ", ")
}

func (q QueryBuilder) Having(condition query.Condition) string {
	return ""
}

func (q QueryBuilder) OrderBy(OrderBy []query.OrderQuery) string {
	return ""
}

func (q QueryBuilder) Offset(n int) string {
	return "OFFSET " + string(n)
}

func (q QueryBuilder) Limit(n int) string {
	return "LIMIT " + string(n)
}

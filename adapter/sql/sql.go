package sql

import (
	"github.com/Fs02/grimoire/query"
	"strconv"
	"strings"
)

type QueryBuilder struct{}

func (q QueryBuilder) Select(distinct bool, fields ...string) string {
	if distinct {
		return "SELECT DISTINCT " + strings.Join(fields, ", ")
	}

	return "SELECT " + strings.Join(fields, ", ")
}

func (q QueryBuilder) From(collection string) string {
	return "FROM " + collection
}

func (q QueryBuilder) Join(join []query.JoinQuery) string {
	return ""
}

func (q QueryBuilder) Where(condition query.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := q.Condition(condition)
	return "WHERE " + qs, args
}

func (q QueryBuilder) GroupBy(fields ...string) string {
	if len(fields) > 0 {
		return "GROUP BY " + strings.Join(fields, ", ")
	}

	return ""
}

func (q QueryBuilder) Having(condition query.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := q.Condition(condition)
	return "HAVING " + qs, args
}

func (q QueryBuilder) OrderBy(orders ...query.OrderQuery) string {
	length := len(orders)
	if length == 0 {
		return ""
	}

	qs := "ORDER BY "
	for i, o := range orders {
		if o.Asc() {
			qs += o.Field + " ASC"
		} else {
			qs += o.Field + " DESC"
		}

		if i < length-1 {
			qs += ", "
		}
	}

	return qs
}

func (q QueryBuilder) Offset(n int) string {
	if n > 0 {
		return "OFFSET " + strconv.Itoa(n)
	}

	return ""
}

func (q QueryBuilder) Limit(n int) string {
	if n > 0 {
		return "LIMIT " + strconv.Itoa(n)
	}

	return ""
}

func (q QueryBuilder) Condition(c query.Condition) (string, []interface{}) {
	build := func(op string, inner []query.Condition) (string, []interface{}) {
		length := len(inner)
		var qstring string
		var args []interface{}

		if length > 1 {
			qstring += "("
		}

		for i, c := range inner {
			cQstring, cArgs := q.Condition(c)
			qstring += cQstring
			args = append(args, cArgs...)

			if i < length-1 {
				qstring += " " + op + " "
			}
		}

		if length > 1 {
			qstring += ")"
		}

		return qstring, args
	}

	switch c.Type {
	case query.ConditionAnd:
		return build("AND", c.Inner)
	case query.ConditionOr:
		return build("OR", c.Inner)
	case query.ConditionXor:
		return build("XOR", c.Inner)
	case query.ConditionNot:
		qs, args := build("AND", c.Inner)
		return "NOT " + qs, args
	case query.ConditionEq:
		return c.Column + " = ?", c.Args
	case query.ConditionNe:
		return c.Column + " <> ?", c.Args
	case query.ConditionLt:
		return c.Column + " < ?", c.Args
	case query.ConditionLte:
		return c.Column + " <= ?", c.Args
	case query.ConditionGt:
		return c.Column + " > ?", c.Args
	case query.ConditionGte:
		return c.Column + " >= ?", c.Args
	case query.ConditionNil:
		return c.Column + " IS NULL", c.Args
	case query.ConditionNotNil:
		return c.Column + " IS NOT NULL", c.Args
	case query.ConditionIn:
		return c.Column + " IN (?" + strings.Repeat(",?", len(c.Args)-1) + ")", c.Args
	case query.ConditionNin:
		return c.Column + " NOT IN (?" + strings.Repeat(",?", len(c.Args)-1) + ")", c.Args
	case query.ConditionLike:
		return c.Column + " LIKE \"" + c.Expr + "\"", c.Args
	case query.ConditionNotLike:
		return c.Column + " NOT LIKE \"" + c.Expr + "\"", c.Args
	case query.ConditionFragment:
		return c.Expr, c.Args
	}

	return "", []interface{}{}
}

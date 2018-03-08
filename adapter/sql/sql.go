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

func (q QueryBuilder) Join(join ...query.JoinClause) (string, []interface{}) {
	if len(join) == 0 {
		return "", nil
	}

	var qs string
	var args []interface{}
	for i, j := range join {
		cs, jargs := q.Condition(j.Condition)
		qs += j.Mode + " " + j.Collection + " ON " + cs
		args = append(args, jargs...)

		if i < len(join)-1 {
			qs += " "
		}
	}

	return qs, args
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

func (q QueryBuilder) OrderBy(orders ...query.OrderClause) string {
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
	switch c.Type {
	case query.ConditionAnd:
		return q.build("AND", c.Inner)
	case query.ConditionOr:
		return q.build("OR", c.Inner)
	case query.ConditionXor:
		return q.build("XOR", c.Inner)
	case query.ConditionNot:
		qs, args := q.build("AND", c.Inner)
		return "NOT " + qs, args
	case query.ConditionEq:
		return q.buildComparison("=", c.Left, c.Right)
	case query.ConditionNe:
		return q.buildComparison("<>", c.Left, c.Right)
	case query.ConditionLt:
		return q.buildComparison("<", c.Left, c.Right)
	case query.ConditionLte:
		return q.buildComparison("<=", c.Left, c.Right)
	case query.ConditionGt:
		return q.buildComparison(">", c.Left, c.Right)
	case query.ConditionGte:
		return q.buildComparison(">=", c.Left, c.Right)
	case query.ConditionNil:
		return string(c.Left.Column) + " IS NULL", c.Right.Values
	case query.ConditionNotNil:
		return string(c.Left.Column) + " IS NOT NULL", c.Right.Values
	case query.ConditionIn:
		return string(c.Left.Column) + " IN (?" + strings.Repeat(",?", len(c.Right.Values)-1) + ")", c.Right.Values
	case query.ConditionNin:
		return string(c.Left.Column) + " NOT IN (?" + strings.Repeat(",?", len(c.Right.Values)-1) + ")", c.Right.Values
	case query.ConditionLike:
		return string(c.Left.Column) + " LIKE ?", c.Right.Values
	case query.ConditionNotLike:
		return string(c.Left.Column) + " NOT LIKE ?", c.Right.Values
	case query.ConditionFragment:
		return string(c.Left.Column), c.Right.Values
	}

	return "", []interface{}{}
}

func (q QueryBuilder) build(op string, inner []query.Condition) (string, []interface{}) {
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

func (q QueryBuilder) buildComparison(op string, left, right query.Operand) (string, []interface{}) {
	var cs string
	if left.Column != "" {
		cs = string(left.Column) + op
	} else {
		cs = "?" + op
	}

	if right.Column != "" {
		cs += string(right.Column)
	} else {
		cs += "?"
	}

	return cs, append(left.Values, right.Values...)
}

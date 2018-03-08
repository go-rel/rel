package sql

import (
	"github.com/Fs02/grimoire/query"
	"strconv"
	"strings"
)

type Builder struct{}

func (b Builder) Select(distinct bool, fields ...string) string {
	if distinct {
		return "SELECT DISTINCT " + strings.Join(fields, ", ")
	}

	return "SELECT " + strings.Join(fields, ", ")
}

func (b Builder) From(collection string) string {
	return "FROM " + collection
}

func (b Builder) Join(join ...query.JoinClause) (string, []interface{}) {
	if len(join) == 0 {
		return "", nil
	}

	var qs string
	var args []interface{}
	for i, j := range join {
		cs, jargs := b.Condition(j.Condition)
		qs += j.Mode + " " + j.Collection + " ON " + cs
		args = append(args, jargs...)

		if i < len(join)-1 {
			qs += " "
		}
	}

	return qs, args
}

func (b Builder) Where(condition query.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := b.Condition(condition)
	return "WHERE " + qs, args
}

func (b Builder) GroupBy(fields ...string) string {
	if len(fields) > 0 {
		return "GROUP BY " + strings.Join(fields, ", ")
	}

	return ""
}

func (b Builder) Having(condition query.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := b.Condition(condition)
	return "HAVING " + qs, args
}

func (b Builder) OrderBy(orders ...query.OrderClause) string {
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

func (b Builder) Offset(n int) string {
	if n > 0 {
		return "OFFSET " + strconv.Itoa(n)
	}

	return ""
}

func (b Builder) Limit(n int) string {
	if n > 0 {
		return "LIMIT " + strconv.Itoa(n)
	}

	return ""
}

func (b Builder) Condition(c query.Condition) (string, []interface{}) {
	switch c.Type {
	case query.ConditionAnd:
		return b.build("AND", c.Inner)
	case query.ConditionOr:
		return b.build("OR", c.Inner)
	case query.ConditionXor:
		return b.build("XOR", c.Inner)
	case query.ConditionNot:
		qs, args := b.build("AND", c.Inner)
		return "NOT " + qs, args
	case query.ConditionEq:
		return b.buildComparison("=", c.Left, c.Right)
	case query.ConditionNe:
		return b.buildComparison("<>", c.Left, c.Right)
	case query.ConditionLt:
		return b.buildComparison("<", c.Left, c.Right)
	case query.ConditionLte:
		return b.buildComparison("<=", c.Left, c.Right)
	case query.ConditionGt:
		return b.buildComparison(">", c.Left, c.Right)
	case query.ConditionGte:
		return b.buildComparison(">=", c.Left, c.Right)
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

func (b Builder) build(op string, inner []query.Condition) (string, []interface{}) {
	length := len(inner)
	var qstring string
	var args []interface{}

	if length > 1 {
		qstring += "("
	}

	for i, c := range inner {
		cQstring, cArgs := b.Condition(c)
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

func (b Builder) buildComparison(op string, left, right query.Operand) (string, []interface{}) {
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

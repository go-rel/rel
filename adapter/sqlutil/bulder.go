package sqlutil

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/Fs02/grimoire/query"
)

type Builder struct{}

func (builder Builder) All(q query.Query) (string, []interface{}) {
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

func (builder Builder) Insert(collection string, changes map[string]interface{}) (string, []interface{}) {
	length := len(changes)

	var buffer bytes.Buffer
	var args = make([]interface{}, length)

	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(collection)
	buffer.WriteString(" (")

	curr := 0
	for field, value := range changes {
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

func (builder Builder) Update(collection string, changes map[string]interface{}, cond query.Condition) (string, []interface{}) {
	length := len(changes)

	var buffer bytes.Buffer
	var args = make([]interface{}, length)

	buffer.WriteString("UPDATE ")
	buffer.WriteString(collection)
	buffer.WriteString(" SET ")

	curr := 0
	for field, value := range changes {
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

func (builder Builder) Delete(collection string, cond query.Condition) (string, []interface{}) {
	var buffer bytes.Buffer
	var args []interface{}

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(collection)
	buffer.WriteString(" ")

	if s, arg := builder.Where(cond); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	return buffer.String(), args
}

func (builder Builder) Select(distinct bool, fields ...string) string {
	if distinct {
		return "SELECT DISTINCT " + strings.Join(fields, ", ")
	}

	return "SELECT " + strings.Join(fields, ", ")
}

func (builder Builder) From(collection string) string {
	return "FROM " + collection
}

func (builder Builder) Join(join ...query.JoinClause) (string, []interface{}) {
	if len(join) == 0 {
		return "", nil
	}

	var qs string
	var args []interface{}
	for i, j := range join {
		cs, jargs := builder.Condition(j.Condition)
		qs += j.Mode + " " + j.Collection + " ON " + cs
		args = append(args, jargs...)

		if i < len(join)-1 {
			qs += " "
		}
	}

	return qs, args
}

func (builder Builder) Where(condition query.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := builder.Condition(condition)
	return "WHERE " + qs, args
}

func (builder Builder) GroupBy(fields ...string) string {
	if len(fields) > 0 {
		return "GROUP BY " + strings.Join(fields, ", ")
	}

	return ""
}

func (builder Builder) Having(condition query.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := builder.Condition(condition)
	return "HAVING " + qs, args
}

func (builder Builder) OrderBy(orders ...query.OrderClause) string {
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

func (builder Builder) Offset(n int) string {
	if n > 0 {
		return "OFFSET " + strconv.Itoa(n)
	}

	return ""
}

func (builder Builder) Limit(n int) string {
	if n > 0 {
		return "LIMIT " + strconv.Itoa(n)
	}

	return ""
}

func (builder Builder) Condition(c query.Condition) (string, []interface{}) {
	switch c.Type {
	case query.ConditionAnd:
		return builder.build("AND", c.Inner)
	case query.ConditionOr:
		return builder.build("OR", c.Inner)
	case query.ConditionXor:
		return builder.build("XOR", c.Inner)
	case query.ConditionNot:
		qs, args := builder.build("AND", c.Inner)
		return "NOT " + qs, args
	case query.ConditionEq:
		return builder.buildComparison("=", c.Left, c.Right)
	case query.ConditionNe:
		return builder.buildComparison("<>", c.Left, c.Right)
	case query.ConditionLt:
		return builder.buildComparison("<", c.Left, c.Right)
	case query.ConditionLte:
		return builder.buildComparison("<=", c.Left, c.Right)
	case query.ConditionGt:
		return builder.buildComparison(">", c.Left, c.Right)
	case query.ConditionGte:
		return builder.buildComparison(">=", c.Left, c.Right)
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

func (builder Builder) build(op string, inner []query.Condition) (string, []interface{}) {
	length := len(inner)
	var qstring string
	var args []interface{}

	if length > 1 {
		qstring += "("
	}

	for i, c := range inner {
		cQstring, cArgs := builder.Condition(c)
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

func (builder Builder) buildComparison(op string, left, right query.Operand) (string, []interface{}) {
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

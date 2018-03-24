package sqlutil

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
)

type Builder struct{}

func (builder Builder) Find(q grimoire.Query) (string, []interface{}) {
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

func (builder Builder) Update(collection string, changes map[string]interface{}, cond c.Condition) (string, []interface{}) {
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

func (builder Builder) Delete(collection string, cond c.Condition) (string, []interface{}) {
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

func (builder Builder) Join(join ...grimoire.JoinClause) (string, []interface{}) {
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

func (builder Builder) Where(condition c.Condition) (string, []interface{}) {
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

func (builder Builder) Having(condition c.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := builder.Condition(condition)
	return "HAVING " + qs, args
}

func (builder Builder) OrderBy(orders ...grimoire.OrderClause) string {
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

func (builder Builder) Condition(cond c.Condition) (string, []interface{}) {
	switch cond.Type {
	case c.ConditionAnd:
		return builder.build("AND", cond.Inner)
	case c.ConditionOr:
		return builder.build("OR", cond.Inner)
	case c.ConditionXor:
		return builder.build("XOR", cond.Inner)
	case c.ConditionNot:
		qs, args := builder.build("AND", cond.Inner)
		return "NOT " + qs, args
	case c.ConditionEq:
		return builder.buildComparison("=", cond.Left, cond.Right)
	case c.ConditionNe:
		return builder.buildComparison("<>", cond.Left, cond.Right)
	case c.ConditionLt:
		return builder.buildComparison("<", cond.Left, cond.Right)
	case c.ConditionLte:
		return builder.buildComparison("<=", cond.Left, cond.Right)
	case c.ConditionGt:
		return builder.buildComparison(">", cond.Left, cond.Right)
	case c.ConditionGte:
		return builder.buildComparison(">=", cond.Left, cond.Right)
	case c.ConditionNil:
		return string(cond.Left.Column) + " IS NULL", cond.Right.Values
	case c.ConditionNotNil:
		return string(cond.Left.Column) + " IS NOT NULL", cond.Right.Values
	case c.ConditionIn:
		return string(cond.Left.Column) + " IN (?" + strings.Repeat(",?", len(cond.Right.Values)-1) + ")", cond.Right.Values
	case c.ConditionNin:
		return string(cond.Left.Column) + " NOT IN (?" + strings.Repeat(",?", len(cond.Right.Values)-1) + ")", cond.Right.Values
	case c.ConditionLike:
		return string(cond.Left.Column) + " LIKE ?", cond.Right.Values
	case c.ConditionNotLike:
		return string(cond.Left.Column) + " NOT LIKE ?", cond.Right.Values
	case c.ConditionFragment:
		return string(cond.Left.Column), cond.Right.Values
	}

	return "", []interface{}{}
}

func (builder Builder) build(op string, inner []c.Condition) (string, []interface{}) {
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

func (builder Builder) buildComparison(op string, left, right c.Operand) (string, []interface{}) {
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

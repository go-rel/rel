package sql

import (
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
)

var fieldCache sync.Map

// Builder defines information of query builder.
type Builder struct {
	config      *Config
	returnField string
	count       int
}

// Find generates query for select.
func (builder *Builder) Find(q grimoire.Query) (string, []interface{}) {
	qs, args := builder.query(q)
	return builder.fields(q.AsDistinct, q.Fields...) + qs, args
}

// Aggregate generates query for aggregation.
func (builder *Builder) Aggregate(q grimoire.Query) (string, []interface{}) {
	qs, args := builder.query(q)

	if q.AggregateMode == "count" && q.AggregateField == "*" {
		return builder.fields(false, "count(*) AS count") + qs, args
	}

	field := q.AggregateMode + "(" + q.AggregateField + ")" + " AS " + q.AggregateMode
	return builder.fields(false, q.AggregateField, field) + qs, args
}

func (builder *Builder) query(q grimoire.Query) (string, []interface{}) {
	var buffer bytes.Buffer
	var args []interface{}

	if s := builder.from(q.Collection); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s, arg := builder.join(q.JoinClause...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s, arg := builder.where(q.Condition); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s := builder.groupBy(q.GroupFields...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)

		if s, arg := builder.having(q.HavingCondition); s != "" {
			buffer.WriteString(" ")
			buffer.WriteString(s)
			args = append(args, arg...)
		}
	}

	if s := builder.orderBy(q.OrderClause...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s := builder.limitOffset(q.LimitResult, q.OffsetResult); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Insert generates query for insert.
func (builder *Builder) Insert(collection string, changes map[string]interface{}) (string, []interface{}) {
	length := len(changes)

	var buffer bytes.Buffer
	var args = make([]interface{}, 0, length)

	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(builder.config.EscapeChar)

	if len(changes) == 0 && builder.config.InsertDefaultValues {
		buffer.WriteString(" DEFAULT VALUES")
	} else {
		buffer.WriteString(" (")

		curr := 0
		for field, value := range changes {
			buffer.WriteString(builder.config.EscapeChar)
			buffer.WriteString(field)
			buffer.WriteString(builder.config.EscapeChar)

			args = append(args, value)

			if curr < length-1 {
				buffer.WriteString(",")
			}

			curr++
		}
		buffer.WriteString(") VALUES ")

		buffer.WriteString("(")
		for i := 0; i < length; i++ {
			buffer.WriteString(builder.ph())

			if i < length-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(")")
	}

	if builder.returnField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteString(builder.config.EscapeChar)
		buffer.WriteString(builder.returnField)
		buffer.WriteString(builder.config.EscapeChar)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// InsertAll generates query for multiple insert.
func (builder *Builder) InsertAll(collection string, fields []string, allchanges []map[string]interface{}) (string, []interface{}) {
	var buffer bytes.Buffer
	var args = make([]interface{}, 0, len(fields)*len(allchanges))

	buffer.WriteString("INSERT INTO ")

	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(builder.config.EscapeChar)

	sep := builder.config.EscapeChar + "," + builder.config.EscapeChar

	buffer.WriteString(" (")
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(strings.Join(fields, sep))
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(") VALUES ")

	for i, changes := range allchanges {
		buffer.WriteString("(")

		for j, field := range fields {
			if val, exist := changes[field]; exist {
				buffer.WriteString(builder.ph())
				args = append(args, val)
			} else {
				buffer.WriteString("DEFAULT")
			}

			if j < len(fields)-1 {
				buffer.WriteString(",")
			}
		}

		if i < len(allchanges)-1 {
			buffer.WriteString("),")
		} else {
			buffer.WriteString(")")
		}
	}

	if builder.returnField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteString(builder.config.EscapeChar)
		buffer.WriteString(builder.returnField)
		buffer.WriteString(builder.config.EscapeChar)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Update generates query for update.
func (builder *Builder) Update(collection string, changes map[string]interface{}, cond c.Condition) (string, []interface{}) {
	length := len(changes)

	var buffer bytes.Buffer
	var args = make([]interface{}, 0, length)

	buffer.WriteString("UPDATE ")
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(" SET ")

	curr := 0
	for field, value := range changes {
		buffer.WriteString(builder.config.EscapeChar)
		buffer.WriteString(field)
		buffer.WriteString(builder.config.EscapeChar)
		buffer.WriteString("=")
		buffer.WriteString(builder.ph())
		args = append(args, value)

		if curr < length-1 {
			buffer.WriteString(",")
		}

		curr++
	}

	if s, arg := builder.where(cond); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Delete generates query for delete.
func (builder *Builder) Delete(collection string, cond c.Condition) (string, []interface{}) {
	var buffer bytes.Buffer
	var args []interface{}

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(builder.config.EscapeChar)

	if s, arg := builder.where(cond); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

func (builder *Builder) fields(distinct bool, fields ...string) string {
	if len(fields) == 0 {
		return "SELECT *"
	}

	var buffer bytes.Buffer

	buffer.WriteString("SELECT ")

	if distinct {
		buffer.WriteString("DISTINCT ")
	}

	l := len(fields) - 1
	for i, f := range fields {
		buffer.WriteString(builder.escape(f))

		if i < l {
			buffer.WriteString(",")
		}
	}

	return buffer.String()
}

func (builder *Builder) from(collection string) string {
	return "FROM " + builder.config.EscapeChar + collection + builder.config.EscapeChar
}

func (builder *Builder) join(join ...c.Join) (string, []interface{}) {
	if len(join) == 0 {
		return "", nil
	}

	var qs string
	var args []interface{}
	for i, j := range join {
		cs, jargs := builder.condition(j.Condition)
		qs += j.Mode + " " + builder.config.EscapeChar + j.Collection + builder.config.EscapeChar + " ON " + cs
		args = append(args, jargs...)

		if i < len(join)-1 {
			qs += " "
		}
	}

	return qs, args
}

func (builder *Builder) where(condition c.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := builder.condition(condition)
	return "WHERE " + qs, args
}

func (builder *Builder) groupBy(fields ...string) string {
	if len(fields) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	buffer.WriteString("GROUP BY ")

	l := len(fields) - 1
	for i, f := range fields {
		buffer.WriteString(builder.escape(f))

		if i < l {
			buffer.WriteString(",")
		}
	}

	return buffer.String()
}

func (builder *Builder) having(condition c.Condition) (string, []interface{}) {
	if condition.None() {
		return "", nil
	}

	qs, args := builder.condition(condition)
	return "HAVING " + qs, args
}

func (builder *Builder) orderBy(orders ...c.Order) string {
	length := len(orders)
	if length == 0 {
		return ""
	}

	qs := "ORDER BY "
	for i, o := range orders {
		if o.Asc() {
			qs += builder.escape(string(o.Field)) + " ASC"
		} else {
			qs += builder.escape(string(o.Field)) + " DESC"
		}

		if i < length-1 {
			qs += ", "
		}
	}

	return qs
}

func (builder *Builder) limitOffset(limit int, offset int) string {
	str := ""

	if limit > 0 {
		str = "LIMIT " + strconv.Itoa(limit)

		if offset > 0 {
			str += " OFFSET " + strconv.Itoa(offset)
		}
	}

	return str
}

func (builder *Builder) condition(cond c.Condition) (string, []interface{}) {
	switch cond.Type {
	case c.ConditionAnd:
		return builder.build("AND", cond.Inner)
	case c.ConditionOr:
		return builder.build("OR", cond.Inner)
	case c.ConditionNot:
		qs, args := builder.build("AND", cond.Inner)
		return "NOT " + qs, args
	case c.ConditionEq,
		c.ConditionNe,
		c.ConditionLt,
		c.ConditionLte,
		c.ConditionGt,
		c.ConditionGte:
		return builder.buildComparison(cond)
	case c.ConditionNil:
		return builder.escape(string(cond.Left.Column)) + " IS NULL", cond.Right.Values
	case c.ConditionNotNil:
		return builder.escape(string(cond.Left.Column)) + " IS NOT NULL", cond.Right.Values
	case c.ConditionIn,
		c.ConditionNin:
		return builder.buildInclusion(cond)
	case c.ConditionLike:
		return builder.escape(string(cond.Left.Column)) + " LIKE " + builder.ph(), cond.Right.Values
	case c.ConditionNotLike:
		return builder.escape(string(cond.Left.Column)) + " NOT LIKE " + builder.ph(), cond.Right.Values
	case c.ConditionFragment:
		return string(cond.Left.Column), cond.Right.Values
	}

	return "", nil
}

func (builder *Builder) build(op string, inner []c.Condition) (string, []interface{}) {
	length := len(inner)
	var qstring string
	var args []interface{}

	if length > 1 {
		qstring += "("
	}

	for i, c := range inner {
		cQstring, cArgs := builder.condition(c)
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

func (builder *Builder) buildComparison(cond c.Condition) (string, []interface{}) {
	var cs string
	var op string

	switch cond.Type {
	case c.ConditionEq:
		op = "="
	case c.ConditionNe:
		op = "<>"
	case c.ConditionLt:
		op = "<"
	case c.ConditionLte:
		op = "<="
	case c.ConditionGt:
		op = ">"
	case c.ConditionGte:
		op = ">="
	}

	if cond.Left.Column != "" {
		cs = builder.escape(string(cond.Left.Column)) + op
	} else {
		cs = builder.ph() + op
	}

	if cond.Right.Column != "" {
		cs += builder.escape(string(cond.Right.Column))
	} else {
		cs += builder.ph()
	}

	return cs, append(cond.Left.Values, cond.Right.Values...)
}

func (builder *Builder) buildInclusion(cond c.Condition) (string, []interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(builder.escape(string(cond.Left.Column)))

	if cond.Type == c.ConditionIn {
		buffer.WriteString(" IN (")
	} else {
		buffer.WriteString(" NOT IN (")
	}

	buffer.WriteString(builder.ph())
	for i := 1; i <= len(cond.Right.Values)-1; i++ {
		buffer.WriteString(",")
		buffer.WriteString(builder.ph())
	}
	buffer.WriteString(")")

	return buffer.String(), cond.Right.Values
}

func (builder *Builder) ph() string {
	if builder.config.Ordinal {
		builder.count++
		return builder.config.Placeholder + strconv.Itoa(builder.count)
	}

	return builder.config.Placeholder
}

func (builder *Builder) escape(field string) string {
	if builder.config.EscapeChar == "" || field == "*" {
		return field
	}

	key := field + builder.config.EscapeChar
	escapedField, ok := fieldCache.Load(key)
	if ok {
		return escapedField.(string)
	}

	start := strings.IndexRune(field, '(')
	end := strings.IndexRune(field, ')')

	if start >= 0 && end >= 0 && end > start {
		escapedField = field[:start+1] + builder.escape(field[start+1:end]) + field[end:]
	} else if strings.HasSuffix(field, "*") {
		escapedField = builder.config.EscapeChar + strings.Replace(field, ".", builder.config.EscapeChar+".", 1)
	} else {
		escapedField = builder.config.EscapeChar +
			strings.Replace(field, ".", builder.config.EscapeChar+"."+builder.config.EscapeChar, 1) +
			builder.config.EscapeChar
	}

	fieldCache.Store(key, escapedField)
	return escapedField.(string)
}

// Returning append returning to insert query.
func (builder *Builder) Returning(field string) *Builder {
	builder.returnField = field
	return builder
}

// NewBuilder create new SQL builder.
func NewBuilder(config *Config) *Builder {
	return &Builder{
		config: config,
	}
}

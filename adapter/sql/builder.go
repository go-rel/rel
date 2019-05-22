package sql

import (
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/Fs02/grimoire/query"
)

// UnescapeCharacter disable field escaping when it starts with this character.
var UnescapeCharacter byte = '^'

var fieldCache sync.Map

// Builder defines information of query builder.
type Builder struct {
	config      *Config
	returnField string
	count       int
}

// Find generates query for select.
func (builder *Builder) Find(q query.Query) (string, []interface{}) {
	qs, args := builder.query(q)
	return builder.fields(q.SelectClause.OnlyDistinct, q.SelectClause.Fields) + qs, args
}

// Aggregate generates query for aggregation.
func (builder *Builder) Aggregate(q query.Query) (string, []interface{}) {
	qs, args := builder.query(q)
	field := "" //q.AggregateMode + "(" + q.AggregateField + ") AS " + q.AggregateMode

	return builder.fields(false, append(q.GroupClause.Fields, field)) + qs, args
}

func (builder *Builder) query(q query.Query) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		args   []interface{}
	)

	if s := builder.from(q.Collection); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s, arg := builder.join(q.JoinClause...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s, arg := builder.where(q.WhereClause); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s := builder.groupBy(q.GroupClause.Fields); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)

		if s, arg := builder.having(q.GroupClause.Filter); s != "" {
			buffer.WriteString(" ")
			buffer.WriteString(s)
			args = append(args, arg...)
		}
	}

	if s := builder.orderBy(q.SortClause); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s := builder.limitOffset(q.LimitClause, q.OffsetClause); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if q.LockClause != "" {
		buffer.WriteString(" ")
		buffer.WriteString(string(q.LockClause))
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Insert generates query for insert.
func (builder *Builder) Insert(collection string, changes map[string]interface{}) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		length = len(changes)
		args   = make([]interface{}, 0, length)
	)

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
	var (
		buffer bytes.Buffer
		args   = make([]interface{}, 0, len(fields)*len(allchanges))
	)

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
func (builder *Builder) Update(collection string, changes map[string]interface{}, filter query.FilterClause) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		length = len(changes)
		args   = make([]interface{}, 0, length)
	)

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

	if s, arg := builder.where(filter); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Delete generates query for delete.
func (builder *Builder) Delete(collection string, filter query.FilterClause) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		args   []interface{}
	)

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(builder.config.EscapeChar)

	if s, arg := builder.where(filter); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

func (builder *Builder) fields(distinct bool, fields []string) string {
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

func (builder *Builder) join(joins ...query.JoinClause) (string, []interface{}) {
	if len(joins) == 0 {
		return "", nil
	}

	var (
		buffer bytes.Buffer
		args   []interface{}
	)

	for i, join := range joins {
		buffer.WriteString(join.Mode)
		buffer.WriteString(builder.config.EscapeChar)
		buffer.WriteString(join.Collection)
		buffer.WriteString(builder.config.EscapeChar)
		buffer.WriteString(" ON ")
		buffer.WriteString(join.From)
		buffer.WriteString(join.To)

		args = append(args, join.Arguments...)

		if i < len(joins)-1 {
			buffer.WriteString(" ")
		}
	}

	return buffer.String(), args
}

func (builder *Builder) where(filter query.FilterClause) (string, []interface{}) {
	if filter.None() {
		return "", nil
	}

	qs, args := builder.filter(filter)
	return "WHERE " + qs, args
}

func (builder *Builder) groupBy(fields []string) string {
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

func (builder *Builder) having(filter query.FilterClause) (string, []interface{}) {
	if filter.None() {
		return "", nil
	}

	qs, args := builder.filter(filter)
	return "HAVING " + qs, args
}

func (builder *Builder) orderBy(orders []query.SortClause) string {
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

func (builder *Builder) limitOffset(limit query.Limit, offset query.Offset) string {
	str := ""

	if limit > 0 {
		str = "LIMIT " + strconv.Itoa(int(limit))

		if offset > 0 {
			str += " OFFSET " + strconv.Itoa(int(offset))
		}
	}

	return str
}

func (builder *Builder) filter(filter query.FilterClause) (string, []interface{}) {
	switch filter.Type {
	case query.AndOp:
		return builder.build("AND", filter.Inner)
	case query.OrOp:
		return builder.build("OR", filter.Inner)
	case query.NotOp:
		qs, args := builder.build("AND", filter.Inner)
		return "NOT " + qs, args
	case query.EqOp,
		query.NeOp,
		query.LtOp,
		query.LteOp,
		query.GtOp,
		query.GteOp:
		return builder.buildComparison(filter)
	case query.NilOp:
		return builder.escape(filter.Field) + " IS NULL", filter.Values
	case query.NotNilOp:
		return builder.escape(filter.Field) + " IS NOT NULL", filter.Values
	case query.InOp,
		query.NinOp:
		return builder.buildInclusion(filter)
	case query.LikeOp:
		return builder.escape(filter.Field) + " LIKE " + builder.ph(), filter.Values
	case query.NotLikeOp:
		return builder.escape(filter.Field) + " NOT LIKE " + builder.ph(), filter.Values
	case query.FragmentOp:
		return filter.Field, filter.Values
	}

	return "", nil
}

func (builder *Builder) build(op string, inner []query.FilterClause) (string, []interface{}) {
	var (
		qstring string
		length  = len(inner)
		args    []interface{}
	)

	if length > 1 {
		qstring += "("
	}

	for i, c := range inner {
		cQstring, cArgs := builder.filter(c)
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

func (builder *Builder) buildComparison(filter query.FilterClause) (string, []interface{}) {
	var (
		cs string
		op string
	)

	switch filter.Type {
	case query.EqOp:
		op = "="
	case query.NeOp:
		op = "<>"
	case query.LtOp:
		op = "<"
	case query.LteOp:
		op = "<="
	case query.GtOp:
		op = ">"
	case query.GteOp:
		op = ">="
	}

	cs = filter.Field + op + builder.ph()

	return cs, filter.Values
}

func (builder *Builder) buildInclusion(filter query.FilterClause) (string, []interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(builder.escape(filter.Field))

	if filter.Type == query.InOp {
		buffer.WriteString(" IN (")
	} else {
		buffer.WriteString(" NOT IN (")
	}

	buffer.WriteString(builder.ph())
	for i := 1; i <= len(filter.Values)-1; i++ {
		buffer.WriteString(",")
		buffer.WriteString(builder.ph())
	}
	buffer.WriteString(")")

	return buffer.String(), filter.Values
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

	if len(field) > 0 && field[0] == UnescapeCharacter {
		escapedField = field[1:]
	} else if start, end := strings.IndexRune(field, '('), strings.IndexRune(field, ')'); start >= 0 && end >= 0 && end > start {
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

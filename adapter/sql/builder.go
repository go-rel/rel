package sql

import (
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/change"
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
func (builder *Builder) Find(query grimoire.Query) (string, []interface{}) {
	qs, args := builder.query(query)
	return builder.fields(query.SelectQuery.OnlyDistinct, query.SelectQuery.Fields) + qs, args
}

// Aggregate generates query for aggregation.
func (builder *Builder) Aggregate(query grimoire.Query) (string, []interface{}) {
	qs, args := builder.query(query)
	field := "" //query.AggregateMode + "(" + query.AggregateField + ") AS " + query.AggregateMode

	return builder.fields(false, append(query.GroupQuery.Fields, field)) + qs, args
}

func (builder *Builder) query(query grimoire.Query) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		args   []interface{}
	)

	if s := builder.from(query.Collection); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s, arg := builder.join(query.JoinQuery...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s, arg := builder.where(query.WhereClause); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s := builder.groupBy(query.GroupQuery.Fields); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)

		if s, arg := builder.having(query.GroupQuery.Filter); s != "" {
			buffer.WriteString(" ")
			buffer.WriteString(s)
			args = append(args, arg...)
		}
	}

	if s := builder.orderBy(query.SortQuery); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s := builder.limitOffset(query.LimitClause, query.OffsetClause); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if query.LockClause != "" {
		buffer.WriteString(" ")
		buffer.WriteString(string(query.LockClause))
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Insert generates query for insert.
func (builder *Builder) Insert(collection string, changes change.Changes) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		length = len(changes.Changes)
		args   = make([]interface{}, 0, length)
	)

	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(builder.escape(collection))

	if length == 0 && builder.config.InsertDefaultValues {
		buffer.WriteString(" DEFAULT VALUES")
	} else {
		buffer.WriteString(" (")

		for i, ch := range changes.Changes {
			switch ch.Type {
			case change.SetOp:
				buffer.WriteString(builder.config.EscapeChar)
				buffer.WriteString(ch.Field)
				buffer.WriteString(builder.config.EscapeChar)
				args = append(args, ch.Value)
			case change.FragmentOp:
				buffer.WriteString(ch.Field)
				args = append(args, ch.Value.([]interface{})...)
			case change.IncOp, change.DecOp:
				continue
			}

			if i < length-1 {
				buffer.WriteString(",")
			}
		}

		buffer.WriteString(") VALUES ")

		buffer.WriteString("(")
		for i := 0; i < len(args); i++ {
			buffer.WriteString(builder.ph())

			if i < len(args)-1 {
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
func (builder *Builder) InsertAll(collection string, fields []string, allchanges []change.Changes) (string, []interface{}) {
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

	// for i, changes := range allchanges {
	// 	buffer.WriteString("(")

	// 	for j, field := range fields {
	// 		if ch, ok := changes.Get(field); ok && ch.Type == change.SetOp {
	// 			buffer.WriteString(builder.ph())
	// 			args = append(args, val)
	// 		} else {
	// 			buffer.WriteString("DEFAULT")
	// 		}

	// 		if j < len(fields)-1 {
	// 			buffer.WriteString(",")
	// 		}
	// 	}

	// 	if i < len(allchanges)-1 {
	// 		buffer.WriteString("),")
	// 	} else {
	// 		buffer.WriteString(")")
	// 	}
	// }

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
func (builder *Builder) Update(collection string, changes change.Changes, filter query.FilterQuery) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		length = len(changes.Changes)
		args   = make([]interface{}, 0, length)
	)

	buffer.WriteString("UPDATE ")
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(builder.config.EscapeChar)
	buffer.WriteString(" SET ")

	for i, ch := range changes.Changes {
		switch ch.Type {
		case change.SetOp:
			buffer.WriteString(builder.escape(ch.Field))
			buffer.WriteString("=")
			buffer.WriteString(builder.ph())
			args = append(args, ch.Value)
		case change.IncOp:
			buffer.WriteString(builder.escape(ch.Field))
			buffer.WriteString("=")
			buffer.WriteString(builder.escape(ch.Field))
			buffer.WriteString("+")
			buffer.WriteString(builder.ph())
			args = append(args, ch.Value)
		case change.DecOp:
			buffer.WriteString(builder.escape(ch.Field))
			buffer.WriteString("=")
			buffer.WriteString(builder.escape(ch.Field))
			buffer.WriteString("-")
			buffer.WriteString(builder.ph())
			args = append(args, ch.Value)
		case change.FragmentOp:
			buffer.WriteString(ch.Field)
			args = append(args, ch.Value.([]interface{})...)
		}

		if i < length-1 {
			buffer.WriteString(",")
		}
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
func (builder *Builder) Delete(collection string, filter query.FilterQuery) (string, []interface{}) {
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
		if distinct {
			return "SELECT DISTINCT *"
		}
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

func (builder *Builder) join(joins ...query.JoinQuery) (string, []interface{}) {
	if len(joins) == 0 {
		return "", nil
	}

	var (
		buffer bytes.Buffer
		args   []interface{}
	)

	for i, join := range joins {
		buffer.WriteString(join.Mode)
		buffer.WriteString(" ")
		buffer.WriteString(builder.config.EscapeChar)
		buffer.WriteString(join.Collection)
		buffer.WriteString(builder.config.EscapeChar)
		buffer.WriteString(" ON ")
		buffer.WriteString(builder.escape(join.From))
		buffer.WriteString("=")
		buffer.WriteString(builder.escape(join.To))

		args = append(args, join.Arguments...)

		if i < len(joins)-1 {
			buffer.WriteString(" ")
		}
	}

	return buffer.String(), args
}

func (builder *Builder) where(filter query.FilterQuery) (string, []interface{}) {
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

func (builder *Builder) having(filter query.FilterQuery) (string, []interface{}) {
	if filter.None() {
		return "", nil
	}

	qs, args := builder.filter(filter)
	return "HAVING " + qs, args
}

func (builder *Builder) orderBy(orders []query.SortQuery) string {
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

func (builder *Builder) filter(filter query.FilterQuery) (string, []interface{}) {
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

func (builder *Builder) build(op string, inner []query.FilterQuery) (string, []interface{}) {
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

func (builder *Builder) buildComparison(filter query.FilterQuery) (string, []interface{}) {
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

	cs = builder.escape(filter.Field) + op + builder.ph()

	return cs, filter.Values
}

func (builder *Builder) buildInclusion(filter query.FilterQuery) (string, []interface{}) {
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

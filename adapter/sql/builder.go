package sql

import (
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/Fs02/grimoire"
)

// UnescapeCharacter disable field escaping when it starts with this character.
var UnescapeCharacter byte = '^'

var fieldCache sync.Map

// Builder defines information of query b.
type Builder struct {
	config      *Config
	returnField string
	count       int
}

// Find generates query for select.
func (b *Builder) Find(query grimoire.Query) (string, []interface{}) {
	var (
		qs, args = b.query(query)
	)

	return b.fields(query.SelectQuery.OnlyDistinct, query.SelectQuery.Fields) + qs, args
}

// Aggregate generates query for aggregation.
func (b *Builder) Aggregate(query grimoire.Query, mode string, field string) (string, []interface{}) {
	var (
		qs, args    = b.query(query)
		selectfield = mode + "(" + field + ") AS " + mode
	)

	return b.fields(false, append(query.GroupQuery.Fields, selectfield)) + qs, args
}

func (b *Builder) query(query grimoire.Query) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		args   []interface{}
	)

	if s := b.from(query.Collection); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s, arg := b.join(query.JoinQuery...); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s, arg := b.where(query.WhereQuery); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	if s := b.groupBy(query.GroupQuery.Fields); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)

		if s, arg := b.having(query.GroupQuery.Filter); s != "" {
			buffer.WriteString(" ")
			buffer.WriteString(s)
			args = append(args, arg...)
		}
	}

	if s := b.orderBy(query.SortQuery); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if s := b.limitOffset(query.LimitQuery, query.OffsetQuery); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
	}

	if query.LockQuery != "" {
		buffer.WriteString(" ")
		buffer.WriteString(string(query.LockQuery))
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Insert generates query for insert.
func (b *Builder) Insert(collection string, changes grimoire.Changes) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		length = len(changes.Changes)
		args   = make([]interface{}, 0, length)
	)

	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(b.escape(collection))

	if length == 0 && b.config.InsertDefaultValues {
		buffer.WriteString(" DEFAULT VALUES")
	} else {
		buffer.WriteString(" (")

		for i, ch := range changes.Changes {
			switch ch.Type {
			case grimoire.ChangeSetOp:
				buffer.WriteString(b.config.EscapeChar)
				buffer.WriteString(ch.Field)
				buffer.WriteString(b.config.EscapeChar)
				args = append(args, ch.Value)
			case grimoire.ChangeFragmentOp:
				buffer.WriteString(ch.Field)
				args = append(args, ch.Value.([]interface{})...)
			case grimoire.ChangeIncOp, grimoire.ChangeDecOp:
				continue
			}

			if i < length-1 {
				buffer.WriteString(",")
			}
		}

		buffer.WriteString(") VALUES ")

		buffer.WriteString("(")
		for i := 0; i < len(args); i++ {
			buffer.WriteString(b.ph())

			if i < len(args)-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(")")
	}

	if b.returnField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(b.returnField)
		buffer.WriteString(b.config.EscapeChar)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// InsertAll generates query for multiple insert.
func (b *Builder) InsertAll(collection string, fields []string, allchanges []grimoire.Changes) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		args   = make([]interface{}, 0, len(fields)*len(allchanges))
	)

	buffer.WriteString("INSERT INTO ")

	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(b.config.EscapeChar)

	sep := b.config.EscapeChar + "," + b.config.EscapeChar

	buffer.WriteString(" (")
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(strings.Join(fields, sep))
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(") VALUES ")

	for i, changes := range allchanges {
		buffer.WriteString("(")

		for j, field := range fields {
			if ch, ok := changes.Get(field); ok && ch.Type == grimoire.ChangeSetOp {
				buffer.WriteString(b.ph())
				args = append(args, ch.Value)
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

	if b.returnField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(b.returnField)
		buffer.WriteString(b.config.EscapeChar)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Update generates query for update.
func (b *Builder) Update(collection string, changes grimoire.Changes, filter grimoire.FilterQuery) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		length = len(changes.Changes)
		args   = make([]interface{}, 0, length)
	)

	buffer.WriteString("UPDATE ")
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(" SET ")

	for i, ch := range changes.Changes {
		switch ch.Type {
		case grimoire.ChangeSetOp:
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("=")
			buffer.WriteString(b.ph())
			args = append(args, ch.Value)
		case grimoire.ChangeIncOp:
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("=")
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("+")
			buffer.WriteString(b.ph())
			args = append(args, ch.Value)
		case grimoire.ChangeDecOp:
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("=")
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("-")
			buffer.WriteString(b.ph())
			args = append(args, ch.Value)
		case grimoire.ChangeFragmentOp:
			buffer.WriteString(ch.Field)
			args = append(args, ch.Value.([]interface{})...)
		}

		if i < length-1 {
			buffer.WriteString(",")
		}
	}

	if s, arg := b.where(filter); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

// Delete generates query for delete.
func (b *Builder) Delete(collection string, filter grimoire.FilterQuery) (string, []interface{}) {
	var (
		buffer bytes.Buffer
		args   []interface{}
	)

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(b.config.EscapeChar)

	if s, arg := b.where(filter); s != "" {
		buffer.WriteString(" ")
		buffer.WriteString(s)
		args = append(args, arg...)
	}

	buffer.WriteString(";")

	return buffer.String(), args
}

func (b *Builder) fields(distinct bool, fields []string) string {
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
		buffer.WriteString(b.escape(f))

		if i < l {
			buffer.WriteString(",")
		}
	}

	return buffer.String()
}

func (b *Builder) from(collection string) string {
	return "FROM " + b.config.EscapeChar + collection + b.config.EscapeChar
}

func (b *Builder) join(joins ...grimoire.JoinQuery) (string, []interface{}) {
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
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(join.Collection)
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(" ON ")
		buffer.WriteString(b.escape(join.From))
		buffer.WriteString("=")
		buffer.WriteString(b.escape(join.To))

		args = append(args, join.Arguments...)

		if i < len(joins)-1 {
			buffer.WriteString(" ")
		}
	}

	return buffer.String(), args
}

func (b *Builder) where(filter grimoire.FilterQuery) (string, []interface{}) {
	if filter.None() {
		return "", nil
	}

	qs, args := b.filter(filter)
	return "WHERE " + qs, args
}

func (b *Builder) groupBy(fields []string) string {
	if len(fields) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	buffer.WriteString("GROUP BY ")

	l := len(fields) - 1
	for i, f := range fields {
		buffer.WriteString(b.escape(f))

		if i < l {
			buffer.WriteString(",")
		}
	}

	return buffer.String()
}

func (b *Builder) having(filter grimoire.FilterQuery) (string, []interface{}) {
	if filter.None() {
		return "", nil
	}

	qs, args := b.filter(filter)
	return "HAVING " + qs, args
}

func (b *Builder) orderBy(orders []grimoire.SortQuery) string {
	length := len(orders)
	if length == 0 {
		return ""
	}

	qs := "ORDER BY "
	for i, o := range orders {
		if o.Asc() {
			qs += b.escape(string(o.Field)) + " ASC"
		} else {
			qs += b.escape(string(o.Field)) + " DESC"
		}

		if i < length-1 {
			qs += ", "
		}
	}

	return qs
}

func (b *Builder) limitOffset(limit grimoire.Limit, offset grimoire.Offset) string {
	str := ""

	if limit > 0 {
		str = "LIMIT " + strconv.Itoa(int(limit))

		if offset > 0 {
			str += " OFFSET " + strconv.Itoa(int(offset))
		}
	}

	return str
}

func (b *Builder) filter(filter grimoire.FilterQuery) (string, []interface{}) {
	switch filter.Type {
	case grimoire.FilterAndOp:
		return b.build("AND", filter.Inner)
	case grimoire.FilterOrOp:
		return b.build("OR", filter.Inner)
	case grimoire.FilterNotOp:
		qs, args := b.build("AND", filter.Inner)
		return "NOT " + qs, args
	case grimoire.FilterEqOp,
		grimoire.FilterNeOp,
		grimoire.FilterLtOp,
		grimoire.FilterLteOp,
		grimoire.FilterGtOp,
		grimoire.FilterGteOp:
		return b.buildComparison(filter)
	case grimoire.FilterNilOp:
		return b.escape(filter.Field) + " IS NULL", filter.Values
	case grimoire.FilterNotNilOp:
		return b.escape(filter.Field) + " IS NOT NULL", filter.Values
	case grimoire.FilterInOp,
		grimoire.FilterNinOp:
		return b.buildInclusion(filter)
	case grimoire.FilterLikeOp:
		return b.escape(filter.Field) + " LIKE " + b.ph(), filter.Values
	case grimoire.FilterNotLikeOp:
		return b.escape(filter.Field) + " NOT LIKE " + b.ph(), filter.Values
	case grimoire.FilterFragmentOp:
		return filter.Field, filter.Values
	}

	return "", nil
}

func (b *Builder) build(op string, inner []grimoire.FilterQuery) (string, []interface{}) {
	var (
		qstring string
		length  = len(inner)
		args    []interface{}
	)

	if length > 1 {
		qstring += "("
	}

	for i, c := range inner {
		cQstring, cArgs := b.filter(c)
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

func (b *Builder) buildComparison(filter grimoire.FilterQuery) (string, []interface{}) {
	var (
		cs string
		op string
	)

	switch filter.Type {
	case grimoire.FilterEqOp:
		op = "="
	case grimoire.FilterNeOp:
		op = "<>"
	case grimoire.FilterLtOp:
		op = "<"
	case grimoire.FilterLteOp:
		op = "<="
	case grimoire.FilterGtOp:
		op = ">"
	case grimoire.FilterGteOp:
		op = ">="
	}

	cs = b.escape(filter.Field) + op + b.ph()

	return cs, filter.Values
}

func (b *Builder) buildInclusion(filter grimoire.FilterQuery) (string, []interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(b.escape(filter.Field))

	if filter.Type == grimoire.FilterInOp {
		buffer.WriteString(" IN (")
	} else {
		buffer.WriteString(" NOT IN (")
	}

	buffer.WriteString(b.ph())
	for i := 1; i <= len(filter.Values)-1; i++ {
		buffer.WriteString(",")
		buffer.WriteString(b.ph())
	}
	buffer.WriteString(")")

	return buffer.String(), filter.Values
}

func (b *Builder) ph() string {
	if b.config.Ordinal {
		b.count++
		return b.config.Placeholder + strconv.Itoa(b.count)
	}

	return b.config.Placeholder
}

func (b *Builder) escape(field string) string {
	if b.config.EscapeChar == "" || field == "*" {
		return field
	}

	key := field + b.config.EscapeChar
	escapedField, ok := fieldCache.Load(key)
	if ok {
		return escapedField.(string)
	}

	if len(field) > 0 && field[0] == UnescapeCharacter {
		escapedField = field[1:]
	} else if start, end := strings.IndexRune(field, '('), strings.IndexRune(field, ')'); start >= 0 && end >= 0 && end > start {
		escapedField = field[:start+1] + b.escape(field[start+1:end]) + field[end:]
	} else if strings.HasSuffix(field, "*") {
		escapedField = b.config.EscapeChar + strings.Replace(field, ".", b.config.EscapeChar+".", 1)
	} else {
		escapedField = b.config.EscapeChar +
			strings.Replace(field, ".", b.config.EscapeChar+"."+b.config.EscapeChar, 1) +
			b.config.EscapeChar
	}

	fieldCache.Store(key, escapedField)
	return escapedField.(string)
}

// Returning append returning to insert grimoire.
func (b *Builder) Returning(field string) *Builder {
	b.returnField = field
	return b
}

// NewBuilder create new SQL builder.
func NewBuilder(config *Config) *Builder {
	return &Builder{
		config: config,
	}
}

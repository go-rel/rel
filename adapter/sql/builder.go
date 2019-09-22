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
		buffer Buffer
	)

	b.fields(&buffer, query.SelectQuery.OnlyDistinct, query.SelectQuery.Fields)
	b.query(&buffer, query)

	return buffer.String(), buffer.Arguments
}

// Aggregate generates query for aggregation.
func (b *Builder) Aggregate(query grimoire.Query, mode string, field string) (string, []interface{}) {
	var (
		buffer      Buffer
		selectfield = mode + "(" + field + ") AS " + mode
	)

	b.fields(&buffer, false, append(query.GroupQuery.Fields, selectfield))
	b.query(&buffer, query)

	return buffer.String(), buffer.Arguments
}

func (b *Builder) query(buffer *Buffer, query grimoire.Query) {
	b.from(buffer, query.Collection)
	b.join(buffer, query.JoinQuery)
	b.where(buffer, query.WhereQuery)

	if len(query.GroupQuery.Fields) > 0 {
		b.groupBy(buffer, query.GroupQuery.Fields)
		b.having(buffer, query.GroupQuery.Filter)
	}

	b.orderBy(buffer, query.SortQuery)
	b.limitOffset(buffer, query.LimitQuery, query.OffsetQuery)

	if query.LockQuery != "" {
		buffer.WriteString(" ")
		buffer.WriteString(string(query.LockQuery))
	}

	buffer.WriteString(";")
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
		buffer Buffer
		length = len(changes.Changes)
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
			buffer.Append(ch.Value)
		case grimoire.ChangeIncOp:
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("=")
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("+")
			buffer.WriteString(b.ph())
			buffer.Append(ch.Value)
		case grimoire.ChangeDecOp:
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("=")
			buffer.WriteString(b.escape(ch.Field))
			buffer.WriteString("-")
			buffer.WriteString(b.ph())
			buffer.Append(ch.Value)
		case grimoire.ChangeFragmentOp:
			buffer.WriteString(ch.Field)
			buffer.Append(ch.Value.([]interface{})...)
		}

		if i < length-1 {
			buffer.WriteString(",")
		}
	}

	b.where(&buffer, filter)

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments
}

// Delete generates query for delete.
func (b *Builder) Delete(collection string, filter grimoire.FilterQuery) (string, []interface{}) {
	var (
		buffer Buffer
	)

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(b.config.EscapeChar)

	b.where(&buffer, filter)

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments
}

func (b *Builder) fields(buffer *Buffer, distinct bool, fields []string) {
	if len(fields) == 0 {
		if distinct {
			buffer.WriteString("SELECT DISTINCT *")
			return
		}
		buffer.WriteString("SELECT *")
		return
	}

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
}

func (b *Builder) from(buffer *Buffer, collection string) {
	buffer.WriteString(" FROM ")
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(collection)
	buffer.WriteString(b.config.EscapeChar)
}

func (b *Builder) join(buffer *Buffer, joins []grimoire.JoinQuery) {
	if len(joins) == 0 {
		return
	}

	for _, join := range joins {
		buffer.WriteString(" ")
		buffer.WriteString(join.Mode)
		buffer.WriteString(" ")
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(join.Collection)
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(" ON ")
		buffer.WriteString(b.escape(join.From))
		buffer.WriteString("=")
		buffer.WriteString(b.escape(join.To))

		buffer.Append(join.Arguments...)
	}
}

func (b *Builder) where(buffer *Buffer, filter grimoire.FilterQuery) {
	if filter.None() {
		return
	}

	buffer.WriteString(" WHERE ")
	b.filter(buffer, filter)
}

func (b *Builder) groupBy(buffer *Buffer, fields []string) {
	buffer.WriteString(" GROUP BY ")

	l := len(fields) - 1
	for i, f := range fields {
		buffer.WriteString(b.escape(f))

		if i < l {
			buffer.WriteString(",")
		}
	}
}

func (b *Builder) having(buffer *Buffer, filter grimoire.FilterQuery) {
	if filter.None() {
		return
	}

	buffer.WriteString(" HAVING ")
	b.filter(buffer, filter)
}

func (b *Builder) orderBy(buffer *Buffer, orders []grimoire.SortQuery) {
	var (
		length = len(orders)
	)

	if length == 0 {
		return
	}

	buffer.WriteString(" ORDER BY")
	for i, order := range orders {
		buffer.WriteByte(' ')
		buffer.WriteString(b.escape(order.Field))

		if order.Asc() {
			buffer.WriteString(" ASC")
		} else {
			buffer.WriteString(" DESC")
		}

		if i < length-1 {
			buffer.WriteByte(',')
		}
	}
}

func (b *Builder) limitOffset(buffer *Buffer, limit grimoire.Limit, offset grimoire.Offset) {
	if limit > 0 {
		buffer.WriteString(" LIMIT ")
		buffer.WriteString(strconv.Itoa(int(limit)))

		if offset > 0 {
			buffer.WriteString(" OFFSET ")
			buffer.WriteString(strconv.Itoa(int(offset)))
		}
	}
}

func (b *Builder) filter(buffer *Buffer, filter grimoire.FilterQuery) {
	switch filter.Type {
	case grimoire.FilterAndOp:
		b.build(buffer, "AND", filter.Inner)
	case grimoire.FilterOrOp:
		b.build(buffer, "OR", filter.Inner)
	case grimoire.FilterNotOp:
		buffer.WriteString("NOT ")
		b.build(buffer, "AND", filter.Inner)
	case grimoire.FilterEqOp,
		grimoire.FilterNeOp,
		grimoire.FilterLtOp,
		grimoire.FilterLteOp,
		grimoire.FilterGtOp,
		grimoire.FilterGteOp:
		b.buildComparison(buffer, filter)
	case grimoire.FilterNilOp:
		buffer.WriteString(b.escape(filter.Field))
		buffer.WriteString(" IS NULL")
	case grimoire.FilterNotNilOp:
		buffer.WriteString(b.escape(filter.Field))
		buffer.WriteString(" IS NOT NULL")
	case grimoire.FilterInOp,
		grimoire.FilterNinOp:
		b.buildInclusion(buffer, filter)
	case grimoire.FilterLikeOp:
		buffer.WriteString(b.escape(filter.Field))
		buffer.WriteString(" LIKE ")
		buffer.WriteString(b.ph())
		buffer.Append(filter.Value)
	case grimoire.FilterNotLikeOp:
		buffer.WriteString(b.escape(filter.Field))
		buffer.WriteString(" NOT LIKE ")
		buffer.WriteString(b.ph())
		buffer.Append(filter.Value)
	case grimoire.FilterFragmentOp:
		buffer.WriteString(filter.Field)
		buffer.Append(filter.Value.([]interface{})...)
	}
}

func (b *Builder) build(buffer *Buffer, op string, inner []grimoire.FilterQuery) {
	var (
		length = len(inner)
	)

	if length > 1 {
		buffer.WriteByte('(')
	}

	for i, c := range inner {
		b.filter(buffer, c)

		if i < length-1 {
			buffer.WriteByte(' ')
			buffer.WriteString(op)
			buffer.WriteByte(' ')
		}
	}

	if length > 1 {
		buffer.WriteByte(')')
	}
}

func (b *Builder) buildComparison(buffer *Buffer, filter grimoire.FilterQuery) {
	buffer.WriteString(b.escape(filter.Field))

	switch filter.Type {
	case grimoire.FilterEqOp:
		buffer.WriteByte('=')
	case grimoire.FilterNeOp:
		buffer.WriteString("<>")
	case grimoire.FilterLtOp:
		buffer.WriteByte('<')
	case grimoire.FilterLteOp:
		buffer.WriteString("<=")
	case grimoire.FilterGtOp:
		buffer.WriteByte('>')
	case grimoire.FilterGteOp:
		buffer.WriteString(">=")
	}

	buffer.WriteString(b.ph())
	buffer.Append(filter.Value)
}

func (b *Builder) buildInclusion(buffer *Buffer, filter grimoire.FilterQuery) {
	var (
		values = filter.Value.([]interface{})
	)

	buffer.WriteString(b.escape(filter.Field))

	if filter.Type == grimoire.FilterInOp {
		buffer.WriteString(" IN (")
	} else {
		buffer.WriteString(" NOT IN (")
	}

	buffer.WriteString(b.ph())
	for i := 1; i <= len(values)-1; i++ {
		buffer.WriteString(",")
		buffer.WriteString(b.ph())
	}
	buffer.WriteString(")")
	buffer.Append(values...)
}

func (b *Builder) ph() string {
	if b.config.Ordinal {
		b.count++
		return b.config.Placeholder + strconv.Itoa(b.count)
	}

	return b.config.Placeholder
}

type fieldCacheKey struct {
	field  string
	escape string
}

func (b *Builder) escape(field string) string {
	if b.config.EscapeChar == "" || field == "*" {
		return field
	}

	key := fieldCacheKey{field: field, escape: b.config.EscapeChar}
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

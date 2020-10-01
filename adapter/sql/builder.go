package sql

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/go-rel/rel"
)

// UnescapeCharacter disable field escaping when it starts with this character.
var UnescapeCharacter byte = '^'

var fieldCache sync.Map

// Builder defines information of query b.
type Builder struct {
	config      Config
	returnField string
	count       int
}

// Table generates query for table creation and modification.
func (b *Builder) Table(table rel.Table) string {
	var buffer Buffer

	switch table.Op {
	case rel.SchemaCreate:
		b.createTable(&buffer, table)
	case rel.SchemaAlter:
		b.alterTable(&buffer, table)
	case rel.SchemaRename:
		buffer.WriteString("ALTER TABLE ")
		buffer.WriteString(Escape(b.config, table.Name))
		buffer.WriteString(" RENAME TO ")
		buffer.WriteString(Escape(b.config, table.Rename))
		buffer.WriteByte(';')
	case rel.SchemaDrop:
		buffer.WriteString("DROP TABLE ")

		if table.Optional {
			buffer.WriteString("IF EXISTS ")
		}

		buffer.WriteString(Escape(b.config, table.Name))
		buffer.WriteByte(';')
	}

	return buffer.String()
}

func (b *Builder) createTable(buffer *Buffer, table rel.Table) {
	buffer.WriteString("CREATE TABLE ")

	if table.Optional {
		buffer.WriteString("IF NOT EXISTS ")
	}

	buffer.WriteString(Escape(b.config, table.Name))
	buffer.WriteString(" (")

	for i, def := range table.Definitions {
		if i > 0 {
			buffer.WriteString(", ")
		}
		switch v := def.(type) {
		case rel.Column:
			b.column(buffer, v)
		case rel.Key:
			b.key(buffer, v)
		case rel.Raw:
			buffer.WriteString(string(v))
		}
	}

	buffer.WriteByte(')')
	b.options(buffer, table.Options)
	buffer.WriteByte(';')
}

func (b *Builder) alterTable(buffer *Buffer, table rel.Table) {
	for _, def := range table.Definitions {
		buffer.WriteString("ALTER TABLE ")
		buffer.WriteString(Escape(b.config, table.Name))
		buffer.WriteByte(' ')

		switch v := def.(type) {
		case rel.Column:
			switch v.Op {
			case rel.SchemaCreate:
				buffer.WriteString("ADD COLUMN ")
				b.column(buffer, v)
			case rel.SchemaRename:
				// Add Change
				buffer.WriteString("RENAME COLUMN ")
				buffer.WriteString(Escape(b.config, v.Name))
				buffer.WriteString(" TO ")
				buffer.WriteString(Escape(b.config, v.Rename))
			case rel.SchemaDrop:
				buffer.WriteString("DROP COLUMN ")
				buffer.WriteString(Escape(b.config, v.Name))
			}
		case rel.Key:
			// TODO: Rename and Drop, PR welcomed.
			switch v.Op {
			case rel.SchemaCreate:
				buffer.WriteString("ADD ")
				b.key(buffer, v)
			}
		}

		b.options(buffer, table.Options)
		buffer.WriteByte(';')
	}
}

func (b *Builder) column(buffer *Buffer, column rel.Column) {
	var (
		typ, m, n = b.config.MapColumnFunc(&column)
	)

	buffer.WriteString(Escape(b.config, column.Name))
	buffer.WriteByte(' ')
	buffer.WriteString(typ)

	if m != 0 {
		buffer.WriteByte('(')
		buffer.WriteString(strconv.Itoa(m))

		if n != 0 {
			buffer.WriteByte(',')
			buffer.WriteString(strconv.Itoa(n))
		}

		buffer.WriteByte(')')
	}

	if column.Unsigned {
		buffer.WriteString(" UNSIGNED")
	}

	if column.Unique {
		buffer.WriteString(" UNIQUE")
	}

	if column.Required {
		buffer.WriteString(" NOT NULL")
	}

	if column.Default != nil {
		buffer.WriteString(" DEFAULT ")
		switch v := column.Default.(type) {
		case string:
			// TODO: single quote only required by postgres.
			buffer.WriteByte('\'')
			buffer.WriteString(v)
			buffer.WriteByte('\'')
		default:
			// TODO: improve
			bytes, _ := json.Marshal(column.Default)
			buffer.Write(bytes)
		}
	}

	b.options(buffer, column.Options)
}

func (b *Builder) key(buffer *Buffer, key rel.Key) {
	var (
		typ = string(key.Type)
	)

	buffer.WriteString(typ)

	if key.Name != "" {
		buffer.WriteByte(' ')
		buffer.WriteString(Escape(b.config, key.Name))
	}

	buffer.WriteString(" (")
	for i, col := range key.Columns {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(Escape(b.config, col))
	}
	buffer.WriteString(")")

	if key.Type == rel.ForeignKey {
		buffer.WriteString(" REFERENCES ")
		buffer.WriteString(Escape(b.config, key.Reference.Table))

		buffer.WriteString(" (")
		for i, col := range key.Reference.Columns {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(Escape(b.config, col))
		}
		buffer.WriteString(")")

		if onDelete := key.Reference.OnDelete; onDelete != "" {
			buffer.WriteString(" ON DELETE ")
			buffer.WriteString(onDelete)
		}

		if onUpdate := key.Reference.OnUpdate; onUpdate != "" {
			buffer.WriteString(" ON UPDATE ")
			buffer.WriteString(onUpdate)
		}
	}

	b.options(buffer, key.Options)
}

// Index generates query for index.
func (b *Builder) Index(index rel.Index) string {
	var buffer Buffer

	switch index.Op {
	case rel.SchemaCreate:
		buffer.WriteString("CREATE ")
		if index.Unique {
			buffer.WriteString("UNIQUE ")
		}
		buffer.WriteString("INDEX ")

		if index.Optional {
			buffer.WriteString("IF NOT EXISTS ")
		}

		buffer.WriteString(Escape(b.config, index.Name))
		buffer.WriteString(" ON ")
		buffer.WriteString(Escape(b.config, index.Table))

		buffer.WriteString(" (")
		for i, col := range index.Columns {
			if i > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(Escape(b.config, col))
		}
		buffer.WriteString(")")
	case rel.SchemaDrop:
		buffer.WriteString("DROP INDEX ")

		if index.Optional {
			buffer.WriteString("IF EXISTS ")
		}

		buffer.WriteString(Escape(b.config, index.Name))

		if b.config.DropIndexOnTable {
			buffer.WriteString(" ON ")
			buffer.WriteString(Escape(b.config, index.Table))
		}
	}

	b.options(&buffer, index.Options)
	buffer.WriteByte(';')

	return buffer.String()
}

func (b *Builder) options(buffer *Buffer, options string) {
	if options == "" {
		return
	}

	buffer.WriteByte(' ')
	buffer.WriteString(options)
}

// Find generates query for select.
func (b *Builder) Find(query rel.Query) (string, []interface{}) {
	if query.SQLQuery.Statement != "" {
		return query.SQLQuery.Statement, query.SQLQuery.Values
	}

	var (
		buffer Buffer
	)

	// TODO: calculate arguments size and if possible buffer size

	b.fields(&buffer, query.SelectQuery.OnlyDistinct, query.SelectQuery.Fields)
	b.query(&buffer, query)

	return buffer.String(), buffer.Arguments
}

// Aggregate generates query for aggregation.
func (b *Builder) Aggregate(query rel.Query, mode string, field string) (string, []interface{}) {
	var (
		buffer Buffer
	)

	buffer.WriteString("SELECT ")
	buffer.WriteString(mode)
	buffer.WriteByte('(')
	buffer.WriteString(Escape(b.config, field))
	buffer.WriteString(") AS ")
	buffer.WriteString(mode)

	for _, f := range query.GroupQuery.Fields {
		buffer.WriteByte(',')
		buffer.WriteString(Escape(b.config, f))
	}

	b.query(&buffer, query)

	return buffer.String(), buffer.Arguments
}

func (b *Builder) query(buffer *Buffer, query rel.Query) {
	b.from(buffer, query.Table)
	b.join(buffer, query.Table, query.JoinQuery)
	b.where(buffer, query.WhereQuery)

	if len(query.GroupQuery.Fields) > 0 {
		b.groupBy(buffer, query.GroupQuery.Fields)
		b.having(buffer, query.GroupQuery.Filter)
	}

	b.orderBy(buffer, query.SortQuery)
	b.limitOffset(buffer, query.LimitQuery, query.OffsetQuery)

	if query.LockQuery != "" {
		buffer.WriteByte(' ')
		buffer.WriteString(string(query.LockQuery))
	}

	buffer.WriteString(";")
}

// Insert generates query for insert.
func (b *Builder) Insert(table string, mutates map[string]rel.Mutate) (string, []interface{}) {
	var (
		buffer Buffer
		count  = len(mutates)
	)

	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(Escape(b.config, table))

	if count == 0 && b.config.InsertDefaultValues {
		buffer.WriteString(" DEFAULT VALUES")
	} else {
		buffer.Arguments = make([]interface{}, count)
		buffer.WriteString(" (")

		i := 0
		for field, mut := range mutates {
			if mut.Type == rel.ChangeSetOp {
				buffer.WriteString(b.config.EscapeChar)
				buffer.WriteString(field)
				buffer.WriteString(b.config.EscapeChar)
				buffer.Arguments[i] = mut.Value
			}

			if i < count-1 {
				buffer.WriteByte(',')
			}
			i++
		}

		buffer.WriteString(") VALUES ")

		buffer.WriteByte('(')
		for i := 0; i < len(buffer.Arguments); i++ {
			buffer.WriteString(b.ph())

			if i < len(buffer.Arguments)-1 {
				buffer.WriteByte(',')
			}
		}
		buffer.WriteByte(')')
	}

	if b.returnField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(b.returnField)
		buffer.WriteString(b.config.EscapeChar)
	}

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments
}

// InsertAll generates query for multiple insert.
func (b *Builder) InsertAll(table string, fields []string, bulkMutates []map[string]rel.Mutate) (string, []interface{}) {
	var (
		buffer       Buffer
		fieldsCount  = len(fields)
		mutatesCount = len(bulkMutates)
	)

	buffer.Arguments = make([]interface{}, 0, fieldsCount*mutatesCount)

	buffer.WriteString("INSERT INTO ")

	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(table)
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(" (")

	for i := range fields {
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(fields[i])
		buffer.WriteString(b.config.EscapeChar)

		if i < fieldsCount-1 {
			buffer.WriteByte(',')
		}
	}

	buffer.WriteString(") VALUES ")

	for i, mutates := range bulkMutates {
		buffer.WriteByte('(')

		for j, field := range fields {
			if mut, ok := mutates[field]; ok && mut.Type == rel.ChangeSetOp {
				buffer.WriteString(b.ph())
				buffer.Append(mut.Value)
			} else {
				buffer.WriteString("DEFAULT")
			}

			if j < fieldsCount-1 {
				buffer.WriteByte(',')
			}
		}

		if i < mutatesCount-1 {
			buffer.WriteString("),")
		} else {
			buffer.WriteByte(')')
		}
	}

	if b.returnField != "" {
		buffer.WriteString(" RETURNING ")
		buffer.WriteString(b.config.EscapeChar)
		buffer.WriteString(b.returnField)
		buffer.WriteString(b.config.EscapeChar)
	}

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments
}

// Update generates query for update.
func (b *Builder) Update(table string, mutates map[string]rel.Mutate, filter rel.FilterQuery) (string, []interface{}) {
	var (
		buffer Buffer
		count  = len(mutates)
	)

	buffer.WriteString("UPDATE ")
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(table)
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(" SET ")

	i := 0
	for field, mut := range mutates {
		switch mut.Type {
		case rel.ChangeSetOp:
			buffer.WriteString(Escape(b.config, field))
			buffer.WriteByte('=')
			buffer.WriteString(b.ph())
			buffer.Append(mut.Value)
		case rel.ChangeIncOp:
			buffer.WriteString(Escape(b.config, field))
			buffer.WriteByte('=')
			buffer.WriteString(Escape(b.config, field))
			buffer.WriteByte('+')
			buffer.WriteString(b.ph())
			buffer.Append(mut.Value)
		case rel.ChangeFragmentOp:
			buffer.WriteString(field)
			buffer.Append(mut.Value.([]interface{})...)
		}

		if i < count-1 {
			buffer.WriteByte(',')
		}
		i++
	}

	b.where(&buffer, filter)

	buffer.WriteString(";")

	return buffer.String(), buffer.Arguments
}

// Delete generates query for delete.
func (b *Builder) Delete(table string, filter rel.FilterQuery) (string, []interface{}) {
	var (
		buffer Buffer
	)

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(table)
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
		buffer.WriteString(Escape(b.config, f))

		if i < l {
			buffer.WriteByte(',')
		}
	}
}

func (b *Builder) from(buffer *Buffer, table string) {
	buffer.WriteString(" FROM ")
	buffer.WriteString(b.config.EscapeChar)
	buffer.WriteString(table)
	buffer.WriteString(b.config.EscapeChar)
}

func (b *Builder) join(buffer *Buffer, table string, joins []rel.JoinQuery) {
	if len(joins) == 0 {
		return
	}

	for _, join := range joins {
		var (
			from = Escape(b.config, join.From)
			to   = Escape(b.config, join.To)
		)

		// TODO: move this to core functionality, and infer join condition using assoc data.
		if join.Arguments == nil && (join.From == "" || join.To == "") {
			from = b.config.EscapeChar + table + b.config.EscapeChar +
				"." + b.config.EscapeChar + strings.TrimSuffix(join.Table, "s") + "_id" + b.config.EscapeChar
			to = b.config.EscapeChar + join.Table + b.config.EscapeChar +
				"." + b.config.EscapeChar + "id" + b.config.EscapeChar
		}

		buffer.WriteByte(' ')
		buffer.WriteString(join.Mode)
		buffer.WriteByte(' ')

		if join.Table != "" {
			buffer.WriteString(b.config.EscapeChar)
			buffer.WriteString(join.Table)
			buffer.WriteString(b.config.EscapeChar)
			buffer.WriteString(" ON ")
			buffer.WriteString(from)
			buffer.WriteString("=")
			buffer.WriteString(to)
		}

		buffer.Append(join.Arguments...)
	}
}

func (b *Builder) where(buffer *Buffer, filter rel.FilterQuery) {
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
		buffer.WriteString(Escape(b.config, f))

		if i < l {
			buffer.WriteByte(',')
		}
	}
}

func (b *Builder) having(buffer *Buffer, filter rel.FilterQuery) {
	if filter.None() {
		return
	}

	buffer.WriteString(" HAVING ")
	b.filter(buffer, filter)
}

func (b *Builder) orderBy(buffer *Buffer, orders []rel.SortQuery) {
	var (
		length = len(orders)
	)

	if length == 0 {
		return
	}

	buffer.WriteString(" ORDER BY")
	for i, order := range orders {
		buffer.WriteByte(' ')
		buffer.WriteString(Escape(b.config, order.Field))

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

func (b *Builder) limitOffset(buffer *Buffer, limit rel.Limit, offset rel.Offset) {
	if limit > 0 {
		buffer.WriteString(" LIMIT ")
		buffer.WriteString(strconv.Itoa(int(limit)))

		if offset > 0 {
			buffer.WriteString(" OFFSET ")
			buffer.WriteString(strconv.Itoa(int(offset)))
		}
	}
}

func (b *Builder) filter(buffer *Buffer, filter rel.FilterQuery) {
	switch filter.Type {
	case rel.FilterAndOp:
		b.build(buffer, "AND", filter.Inner)
	case rel.FilterOrOp:
		b.build(buffer, "OR", filter.Inner)
	case rel.FilterNotOp:
		buffer.WriteString("NOT ")
		b.build(buffer, "AND", filter.Inner)
	case rel.FilterEqOp,
		rel.FilterNeOp,
		rel.FilterLtOp,
		rel.FilterLteOp,
		rel.FilterGtOp,
		rel.FilterGteOp:
		b.buildComparison(buffer, filter)
	case rel.FilterNilOp:
		buffer.WriteString(Escape(b.config, filter.Field))
		buffer.WriteString(" IS NULL")
	case rel.FilterNotNilOp:
		buffer.WriteString(Escape(b.config, filter.Field))
		buffer.WriteString(" IS NOT NULL")
	case rel.FilterInOp,
		rel.FilterNinOp:
		b.buildInclusion(buffer, filter)
	case rel.FilterLikeOp:
		buffer.WriteString(Escape(b.config, filter.Field))
		buffer.WriteString(" LIKE ")
		buffer.WriteString(b.ph())
		buffer.Append(filter.Value)
	case rel.FilterNotLikeOp:
		buffer.WriteString(Escape(b.config, filter.Field))
		buffer.WriteString(" NOT LIKE ")
		buffer.WriteString(b.ph())
		buffer.Append(filter.Value)
	case rel.FilterFragmentOp:
		buffer.WriteString(filter.Field)
		buffer.Append(filter.Value.([]interface{})...)
	}
}

func (b *Builder) build(buffer *Buffer, op string, inner []rel.FilterQuery) {
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

func (b *Builder) buildComparison(buffer *Buffer, filter rel.FilterQuery) {
	buffer.WriteString(Escape(b.config, filter.Field))

	switch filter.Type {
	case rel.FilterEqOp:
		buffer.WriteByte('=')
	case rel.FilterNeOp:
		buffer.WriteString("<>")
	case rel.FilterLtOp:
		buffer.WriteByte('<')
	case rel.FilterLteOp:
		buffer.WriteString("<=")
	case rel.FilterGtOp:
		buffer.WriteByte('>')
	case rel.FilterGteOp:
		buffer.WriteString(">=")
	}

	buffer.WriteString(b.ph())
	buffer.Append(filter.Value)
}

func (b *Builder) buildInclusion(buffer *Buffer, filter rel.FilterQuery) {
	var (
		values = filter.Value.([]interface{})
	)

	buffer.WriteString(Escape(b.config, filter.Field))

	if filter.Type == rel.FilterInOp {
		buffer.WriteString(" IN (")
	} else {
		buffer.WriteString(" NOT IN (")
	}

	buffer.WriteString(b.ph())
	for i := 1; i <= len(values)-1; i++ {
		buffer.WriteByte(',')
		buffer.WriteString(b.ph())
	}
	buffer.WriteByte(')')
	buffer.Append(values...)
}

func (b *Builder) ph() string {
	if b.config.Ordinal {
		b.count++
		return b.config.Placeholder + strconv.Itoa(b.count)
	}

	return b.config.Placeholder
}

// Returning append returning to insert rel.
func (b *Builder) Returning(field string) *Builder {
	b.returnField = field
	return b
}

// NewBuilder create new SQL builder.
func NewBuilder(config Config) *Builder {
	return &Builder{
		config: config,
	}
}

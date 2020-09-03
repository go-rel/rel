package rel

// TableOption interface.
// Available options are: Comment, Options.
type TableOption interface {
	applyTable(table *Table)
}

func applyTableOptions(table *Table, options []TableOption) {
	for i := range options {
		options[i].applyTable(table)
	}
}

// ColumnOption interface.
// Available options are: Nil, Unsigned, Limit, Precision, Scale, Default, Comment, Options.
type ColumnOption interface {
	applyColumn(column *Column)
}

func applyColumnOptions(column *Column, options []ColumnOption) {
	for i := range options {
		options[i].applyColumn(column)
	}
}

// KeyOption interface.
// Available options are: Comment, Options.
type KeyOption interface {
	applyKey(key *Key)
}

func applyKeyOptions(key *Key, options []KeyOption) {
	for i := range options {
		options[i].applyKey(key)
	}
}

// Unique set column as unique.
type Unique bool

func (r Unique) applyColumn(column *Column) {
	column.Unique = bool(r)
}

func (r Unique) applyIndex(index *Index) {
	index.Unique = bool(r)
}

// Required disallows nil values in the column.
type Required bool

func (r Required) applyColumn(column *Column) {
	column.Required = bool(r)
}

// Unsigned sets integer column to be unsigned.
type Unsigned bool

func (u Unsigned) applyColumn(column *Column) {
	column.Unsigned = bool(u)
}

// Precision defines the precision for the decimal fields, representing the total number of digits in the number.
type Precision int

func (p Precision) applyColumn(column *Column) {
	column.Precision = int(p)
}

// Scale Defines the scale for the decimal fields, representing the number of digits after the decimal point.
type Scale int

func (s Scale) applyColumn(column *Column) {
	column.Scale = int(s)
}

type defaultValue struct {
	value interface{}
}

func (d defaultValue) applyColumn(column *Column) {
	column.Default = d.value
}

// Default allows to set a default value on the column.).
func Default(def interface{}) ColumnOption {
	return defaultValue{value: def}
}

// OnDelete option for foreign key.
type OnDelete string

func (od OnDelete) applyKey(key *Key) {
	key.Reference.OnDelete = string(od)
}

// OnUpdate option for foreign key.
type OnUpdate string

func (ou OnUpdate) applyKey(key *Key) {
	key.Reference.OnUpdate = string(ou)
}

// Options options for table, column and index.
type Options string

func (o Options) applyTable(table *Table) {
	table.Options = string(o)
}

func (o Options) applyColumn(column *Column) {
	column.Options = string(o)
}

func (o Options) applyIndex(index *Index) {
	index.Options = string(o)
}

func (o Options) applyKey(key *Key) {
	key.Options = string(o)
}

// Optional option.
// when used with create table, will create table only if it's not exists.
// when used with drop table, will drop table only if it's exists.
type Optional bool

func (o Optional) applyTable(table *Table) {
	table.Optional = bool(o)
}

func (o Optional) applyIndex(index *Index) {
	index.Optional = bool(o)
}

// Raw string
type Raw string

func (r Raw) internalMigration()       {}
func (r Raw) internalTableDefinition() {}

// Do used internally for schema migration.
type Do func(Repository) error

func (d Do) internalMigration() {}

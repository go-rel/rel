package schema

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

// IndexOption interface.
// Available options are: Comment, Options.
type IndexOption interface {
	applyIndex(index *Index)
}

func applyIndexOptions(index *Index, options []IndexOption) {
	for i := range options {
		options[i].applyIndex(index)
	}
}

// Comment options for table, column and index.
type Comment string

func (c Comment) applyTable(table *Table) {
	table.Comment = string(c)
}

func (c Comment) applyColumn(column *Column) {
	column.Comment = string(c)
}

func (c Comment) applyIndex(index *Index) {
	index.Comment = string(c)
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

// Limit sets the maximum size of the string/text/binary/integer columns.
type Limit int

func (l Limit) applyColumn(column *Column) {
	column.Limit = int(l)
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

// Name option for defining custom index name.
type Name string

func (n Name) applyIndex(index *Index) {
	index.Name = string(n)
}

// OnDelete option for foreign key index.
type OnDelete string

func (od OnDelete) applyIndex(index *Index) {
	index.Reference.OnDelete = string(od)
}

// OnUpdate option for foreign key index.
type OnUpdate string

func (ou OnUpdate) applyIndex(index *Index) {
	index.Reference.OnUpdate = string(ou)
}

package schema

// ColumnOp definition.
type ColumnOp uint8

const (
	// AddColumnOp operation.
	AddColumnOp ColumnOp = iota
	// AlterColumnOp operation.
	AlterColumnOp
	// RenameColumnOp operation.
	RenameColumnOp
	// DropColumnOp operation.
	DropColumnOp
)

// ColumnType definition.
type ColumnType string

const (
	// Boolean ColumnType.
	Boolean ColumnType = "boolean"
	// Integer ColumnType.
	Integer ColumnType = "integer"
	// BigInt ColumnType.
	BigInt ColumnType = "bigint"
	// Float ColumnType.
	Float ColumnType = "float"
	// Decimal ColumnType.
	Decimal ColumnType = "decimal"
	// String ColumnType.
	String ColumnType = "string"
	// Text ColumnType.
	Text ColumnType = "text"
	// Binary ColumnType.
	Binary ColumnType = "binary"
	// Date ColumnType.
	Date ColumnType = "date"
	// DateTime ColumnType.
	DateTime ColumnType = "datetime"
	// Time ColumnType.
	Time ColumnType = "time"
	// Timestamp ColumnType.
	Timestamp ColumnType = "timestamp"
)

// Column definition.
type Column struct {
	Op        ColumnOp
	Name      string
	Type      ColumnType
	NewName   string
	Nil       bool
	Limit     int
	Precision int
	Scale     int
	Comment   string
	Default   interface{}
}

func addColumn(name string, typ ColumnType, options []ColumnOption) Column {
	column := Column{
		Op:   AddColumnOp,
		Name: name,
		Type: typ,
	}

	applyColumnOptions(&column, options)
	return column
}

func alterColumn(name string, typ ColumnType, options []ColumnOption) Column {
	column := Column{
		Op:   AlterColumnOp,
		Name: name,
		Type: typ,
	}

	applyColumnOptions(&column, options)
	return column
}

func renameColumn(name string, newName string, options []ColumnOption) Column {
	column := Column{
		Op:      RenameColumnOp,
		Name:    name,
		NewName: newName,
	}

	applyColumnOptions(&column, options)
	return column
}

func dropColumn(name string, options []ColumnOption) Column {
	column := Column{
		Op:   DropColumnOp,
		Name: name,
	}

	applyColumnOptions(&column, options)
	return column
}

// ColumnOption functor.
type ColumnOption func(column *Column)

// Nil allows or disallows nil values in the column.
func Nil(allow bool) ColumnOption {
	return func(column *Column) {
		column.Nil = allow
	}
}

// Limit sets the maximum size of the string/text/binary/integer columns.
func Limit(limit int) ColumnOption {
	return func(column *Column) {
		column.Limit = limit
	}
}

// Precision defines the precision for the decimal fields, representing the total number of digits in the number.
func Precision(precision int) ColumnOption {
	return func(column *Column) {
		column.Precision = precision
	}
}

// Scale Defines the scale for the decimal fields, representing the number of digits after the decimal point.
func Scale(scale int) ColumnOption {
	return func(column *Column) {
		column.Scale = scale
	}
}

// Comment adds a comment for the column.
func Comment(comment string) ColumnOption {
	return func(column *Column) {
		column.Comment = comment
	}
}

// Default allows to set a default value on the column.).
func Default(def interface{}) ColumnOption {
	return func(column *Column) {
		column.Default = def
	}
}

func applyColumnOptions(table *Column, options []ColumnOption) {
	for i := range options {
		options[i](table)
	}
}

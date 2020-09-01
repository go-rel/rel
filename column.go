package rel

// ColumnType definition.
type ColumnType string

const (
	// ID ColumnType.
	ID ColumnType = "ID"
	// Bool ColumnType.
	Bool ColumnType = "BOOL"
	// Int ColumnType.
	Int ColumnType = "INT"
	// BigInt ColumnType.
	BigInt ColumnType = "BIGINT"
	// Float ColumnType.
	Float ColumnType = "FLOAT"
	// Decimal ColumnType.
	Decimal ColumnType = "DECIMAL"
	// String ColumnType.
	String ColumnType = "STRING"
	// Text ColumnType.
	Text ColumnType = "TEXT"
	// Date ColumnType.
	Date ColumnType = "DATE"
	// DateTime ColumnType.
	DateTime ColumnType = "DATETIME"
	// Time ColumnType.
	Time ColumnType = "TIME"
	// Timestamp ColumnType.
	Timestamp ColumnType = "TIMESTAMP"
)

// Column definition.
type Column struct {
	Op        SchemaOp
	Name      string
	Type      ColumnType
	Rename    string
	Unique    bool
	Required  bool
	Unsigned  bool
	Limit     int
	Precision int
	Scale     int
	Default   interface{}
	Options   string
}

func (Column) internalTableDefinition() {}

func createColumn(name string, typ ColumnType, options []ColumnOption) Column {
	column := Column{
		Op:   SchemaCreate,
		Name: name,
		Type: typ,
	}

	applyColumnOptions(&column, options)
	return column
}

func renameColumn(name string, newName string, options []ColumnOption) Column {
	column := Column{
		Op:     SchemaRename,
		Name:   name,
		Rename: newName,
	}

	applyColumnOptions(&column, options)
	return column
}

func dropColumn(name string, options []ColumnOption) Column {
	column := Column{
		Op:   SchemaDrop,
		Name: name,
	}

	applyColumnOptions(&column, options)
	return column
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

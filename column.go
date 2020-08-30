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
	NewName   string
	Required  bool
	Unsigned  bool
	Limit     int
	Precision int
	Scale     int
	Default   interface{}
	Options   string
}

func createColumn(name string, typ ColumnType, options []ColumnOption) Column {
	column := Column{
		Op:   SchemaCreate,
		Name: name,
		Type: typ,
	}

	applyColumnOptions(&column, options)
	return column
}

func alterColumn(name string, typ ColumnType, options []ColumnOption) Column {
	column := Column{
		Op:   SchemaAlter,
		Name: name,
		Type: typ,
	}

	applyColumnOptions(&column, options)
	return column
}

func renameColumn(name string, newName string, options []ColumnOption) Column {
	column := Column{
		Op:      SchemaRename,
		Name:    name,
		NewName: newName,
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

package schema

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
	Op        Op
	Name      string
	Type      ColumnType
	NewName   string
	Required  bool
	Unsigned  bool
	Limit     int
	Precision int
	Scale     int
	Default   interface{}
	Comment   string
	Options   string
}

func addColumn(name string, typ ColumnType, options []ColumnOption) Column {
	column := Column{
		Op:   Add,
		Name: name,
		Type: typ,
	}

	applyColumnOptions(&column, options)
	return column
}

func alterColumn(name string, typ ColumnType, options []ColumnOption) Column {
	column := Column{
		Op:   Alter,
		Name: name,
		Type: typ,
	}

	applyColumnOptions(&column, options)
	return column
}

func renameColumn(name string, newName string, options []ColumnOption) Column {
	column := Column{
		Op:      Rename,
		Name:    name,
		NewName: newName,
	}

	applyColumnOptions(&column, options)
	return column
}

func dropColumn(name string, options []ColumnOption) Column {
	column := Column{
		Op:   Drop,
		Name: name,
	}

	applyColumnOptions(&column, options)
	return column
}

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
	Limit     int
	Default   interface{}
	Null      bool
	Precision int
	Scale     int
	Comment   string
}

func addColumn(name string, typ ColumnType, options []Option) Column {
	return Column{
		Op:   AddColumnOp,
		Name: name,
		Type: typ,
	}
}

func alterColumn(name string, typ ColumnType, options []Option) Column {
	return Column{
		Op:   AlterColumnOp,
		Name: name,
		Type: typ,
	}
}

func renameColumn(name string, newName string, options []Option) Column {
	return Column{
		Op:      RenameColumnOp,
		Name:    name,
		NewName: newName,
	}
}

func dropColumn(name string, options []Option) Column {
	return Column{
		Op:   DropColumnOp,
		Name: name,
	}
}

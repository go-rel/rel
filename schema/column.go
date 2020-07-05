package schema

// ColumnOp definition.
type ColumnOp uint8

const (
	// AddColumn operation.
	AddColumn ColumnOp = iota
	// AlterColumn operation.
	AlterColumn
	// RenameColumn operation.
	RenameColumn
	// DropColumn operation.
	DropColumn
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
	Op      ColumnOp
	Name    string
	Type    ColumnType
	NewName string
}

func addColumn(name string, typ ColumnType, options ...Option) Column {
	return Column{
		Op:   AddColumn,
		Name: name,
		Type: typ,
	}
}

func alterColumn(name string, typ ColumnType, options ...Option) Column {
	return Column{
		Op:   AlterColumn,
		Name: name,
		Type: typ,
	}
}

func renameColumn(name string, newName string, options ...Option) Column {
	return Column{
		Op:      RenameColumn,
		Name:    name,
		NewName: newName,
	}
}

func dropColumn(name string, options ...Option) Column {
	return Column{
		Op:   DropColumn,
		Name: name,
	}
}

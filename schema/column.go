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
	Name string
	Type ColumnType
}

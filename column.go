package rel

// ColumnType definition.
type ColumnType string

const (
	// ID ColumnType.
	ID ColumnType = "ID"
	// BigID ColumnType.
	BigID ColumnType = "BigID"
	// Bool ColumnType.
	Bool ColumnType = "BOOL"
	// SmallInt ColumnType.
	SmallInt ColumnType = "SMALLINT"
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
	// JSON ColumnType that will fallback to Text ColumnType if adapter does not support it.
	JSON ColumnType = "JSON"
	// Date ColumnType.
	Date ColumnType = "DATE"
	// DateTime ColumnType.
	DateTime ColumnType = "DATETIME"
	// Time ColumnType.
	Time ColumnType = "TIME"
)

// Column definition.
type Column struct {
	Op        SchemaOp
	Name      string
	Type      ColumnType
	Rename    string
	Primary   bool
	Unique    bool
	Required  bool
	Unsigned  bool
	Limit     int
	Precision int
	Scale     int
	Default   any
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

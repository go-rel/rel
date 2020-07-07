package schema

// TableOp definition.
type TableOp uint8

const (
	// CreateTableOp operation.
	CreateTableOp TableOp = iota
	// AlterTableOp operation.
	AlterTableOp
	// RenameTableOp operation.
	RenameTableOp
	// DropTableOp operation.
	DropTableOp
)

// Table definition.
type Table struct {
	Op      TableOp
	Name    string
	NewName string
	Columns []Column
	Indices []Index
	Options string
}

// Column defines a column with name and type.
func (t *Table) Column(name string, typ ColumnType, options ...ColumnOption) {
	t.Columns = append(t.Columns, addColumn(name, typ, options))
}

// Index defines an index for columns.
func (t *Table) Index(columns []string, options ...ColumnOption) {
	index := Index{Columns: columns}
	t.Indices = append(t.Indices, index)
}

// Boolean defines a column with name and Boolean type.
func (t *Table) Boolean(name string, options ...ColumnOption) {
	t.Column(name, Boolean, options...)
}

// Integer defines a column with name and Integer type.
func (t *Table) Integer(name string, options ...ColumnOption) {
	t.Column(name, Integer, options...)
}

// BigInt defines a column with name and BigInt type.
func (t *Table) BigInt(name string, options ...ColumnOption) {
	t.Column(name, BigInt, options...)
}

// Float defines a column with name and Float type.
func (t *Table) Float(name string, options ...ColumnOption) {
	t.Column(name, Float, options...)
}

// Decimal defines a column with name and Decimal type.
func (t *Table) Decimal(name string, options ...ColumnOption) {
	t.Column(name, Decimal, options...)
}

// String defines a column with name and String type.
func (t *Table) String(name string, options ...ColumnOption) {
	t.Column(name, String, options...)
}

// Text defines a column with name and Text type.
func (t *Table) Text(name string, options ...ColumnOption) {
	t.Column(name, Text, options...)
}

// Binary defines a column with name and Binary type.
func (t *Table) Binary(name string, options ...ColumnOption) {
	t.Column(name, Binary, options...)
}

// Date defines a column with name and Date type.
func (t *Table) Date(name string, options ...ColumnOption) {
	t.Column(name, Date, options...)
}

// DateTime defines a column with name and DateTime type.
func (t *Table) DateTime(name string, options ...ColumnOption) {
	t.Column(name, DateTime, options...)
}

// Time defines a column with name and Time type.
func (t *Table) Time(name string, options ...ColumnOption) {
	t.Column(name, Time, options...)
}

// Timestamp defines a column with name and Timestamp type.
func (t *Table) Timestamp(name string, options ...ColumnOption) {
	t.Column(name, Timestamp, options...)
}

func (t Table) migrate() {}

// AlterTable Migrator.
type AlterTable struct {
	Table
}

// RenameColumn to a new name.
func (at *AlterTable) RenameColumn(name string, newName string, options ...ColumnOption) {
	at.Columns = append(at.Columns, renameColumn(name, newName, options))
}

// AlterColumn from this table.
func (at *AlterTable) AlterColumn(name string, typ ColumnType, options ...ColumnOption) {
	at.Columns = append(at.Columns, alterColumn(name, typ, options))
}

// DropColumn from this table.
func (at *AlterTable) DropColumn(name string, options ...ColumnOption) {
	at.Columns = append(at.Columns, dropColumn(name, options))
}

func createTable(name string, options []TableOption) Table {
	table := Table{
		Op:   CreateTableOp,
		Name: name,
	}

	applyTableOptions(&table, options)
	return table
}

func alterTable(name string, options []TableOption) AlterTable {
	table := Table{
		Op:   AlterTableOp,
		Name: name,
	}

	applyTableOptions(&table, options)
	return AlterTable{Table: table}
}

func renameTable(name string, newName string, options []TableOption) Table {
	table := Table{
		Op:      RenameTableOp,
		Name:    name,
		NewName: newName,
	}

	applyTableOptions(&table, options)
	return table
}

func dropTable(name string, options []TableOption) Table {
	table := Table{
		Op:   DropTableOp,
		Name: name,
	}

	applyTableOptions(&table, options)
	return table
}

// TableOption functor.
type TableOption func(table *Table)

// Options allow additional SQL fragment to be used when creating a table.
func Options(options string) TableOption {
	return func(table *Table) {
		table.Options = options
	}
}

func applyTableOptions(table *Table, options []TableOption) {
	for i := range options {
		options[i](table)
	}
}

package schema

// Table definition.
type Table struct {
	Name    string
	Columns []Column
	Indices []Index
}

// Column defines a column with name and type.
func (t *Table) Column(name string, typ ColumnType, options ...Option) {
	t.Columns = append(t.Columns, addColumn(name, typ, options...))
}

// Index defines an index for columns.
func (t *Table) Index(columns []string, options ...Option) {
	index := Index{Columns: columns}
	t.Indices = append(t.Indices, index)
}

// Boolean defines a column with name and Boolean type.
func (t *Table) Boolean(name string, options ...Option) {
	t.Column(name, Boolean, options...)
}

// Integer defines a column with name and Integer type.
func (t *Table) Integer(name string, options ...Option) {
	t.Column(name, Integer, options...)
}

// BigInt defines a column with name and BigInt type.
func (t *Table) BigInt(name string, options ...Option) {
	t.Column(name, BigInt, options...)
}

// Float defines a column with name and Float type.
func (t *Table) Float(name string, options ...Option) {
	t.Column(name, Float, options...)
}

// Decimal defines a column with name and Decimal type.
func (t *Table) Decimal(name string, options ...Option) {
	t.Column(name, Decimal, options...)
}

// String defines a column with name and String type.
func (t *Table) String(name string, options ...Option) {
	t.Column(name, String, options...)
}

// Text defines a column with name and Text type.
func (t *Table) Text(name string, options ...Option) {
	t.Column(name, Text, options...)
}

// Date defines a column with name and Date type.
func (t *Table) Date(name string, options ...Option) {
	t.Column(name, Date, options...)
}

// DateTime defines a column with name and DateTime type.
func (t *Table) DateTime(name string, options ...Option) {
	t.Column(name, DateTime, options...)
}

// Time defines a column with name and Time type.
func (t *Table) Time(name string, options ...Option) {
	t.Column(name, Time, options...)
}

// Timestamp defines a column with name and Timestamp type.
func (t *Table) Timestamp(name string, options ...Option) {
	t.Column(name, Timestamp, options...)
}

// CreateTable Migrator.
type CreateTable Table

func (ct CreateTable) migrate() {}

// AlterTable Migrator.
type AlterTable struct {
	Table
}

// RenameColumn to a new name.
func (at *AlterTable) RenameColumn(name string, newName string, options ...Option) {
	at.Columns = append(at.Columns, renameColumn(name, newName, options...))
}

// AlterColumn from this table.
func (at *AlterTable) AlterColumn(name string, typ ColumnType, options ...Option) {
	at.Columns = append(at.Columns, alterColumn(name, typ, options...))
}

// DropColumn from this table.
func (at *AlterTable) DropColumn(name string, options ...Option) {
	at.Columns = append(at.Columns, dropColumn(name, options...))
}

func (at AlterTable) migrate() {}

// RenameTable Migrator.
type RenameTable struct {
	Name    string
	NewName string
}

func (rt RenameTable) migrate() {}

// DropTable Migrator.
type DropTable struct {
	Name string
}

func (dt DropTable) migrate() {}

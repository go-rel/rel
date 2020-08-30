package rel

// Table definition.
type Table struct {
	Op          SchemaOp
	Name        string
	NewName     string
	Definitions []interface{}
	Optional    bool
	Options     string
}

// Column defines a column with name and type.
func (t *Table) Column(name string, typ ColumnType, options ...ColumnOption) {
	t.Definitions = append(t.Definitions, createColumn(name, typ, options))
}

// ID defines a column with name and ID type.
// the resulting database type will depends on database.
func (t *Table) ID(name string, options ...ColumnOption) {
	t.Column(name, ID, options...)
}

// Bool defines a column with name and Bool type.
func (t *Table) Bool(name string, options ...ColumnOption) {
	t.Column(name, Bool, options...)
}

// Int defines a column with name and Int type.
func (t *Table) Int(name string, options ...ColumnOption) {
	t.Column(name, Int, options...)
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

// // Index defines an index for columns.
// func (t *Table) Index(columns []string, typ IndexType, options ...IndexOption) {
// 	t.Definitions = append(t.Definitions, createIndex(columns, typ, options))
// }

// PrimaryKey defines a primary key for table.
func (t *Table) PrimaryKey(column string, options ...KeyOption) {
	t.PrimaryKeys([]string{column}, options...)
}

// PrimaryKeys defines composite primary keys for table.
func (t *Table) PrimaryKeys(columns []string, options ...KeyOption) {
	t.Definitions = append(t.Definitions, createPrimaryKeys(columns, options))
}

// ForeignKey defines foreign key index.
func (t *Table) ForeignKey(column string, refTable string, refColumn string, options ...KeyOption) {
	t.Definitions = append(t.Definitions, createForeignKey(column, refTable, refColumn, options))
}

// Unique defines an unique key for columns.
func (t *Table) Unique(columns []string, options ...KeyOption) {
	t.Definitions = append(t.Definitions, createKeys(columns, UniqueKey, options))
}

// Fragment defines anything using sql fragment.
func (t *Table) Fragment(fragment string) {
	t.Definitions = append(t.Definitions, fragment)
}

func (t Table) migrate() {}

// AlterTable Migrator.
type AlterTable struct {
	Table
}

// RenameColumn to a new name.
func (at *AlterTable) RenameColumn(name string, newName string, options ...ColumnOption) {
	at.Definitions = append(at.Definitions, renameColumn(name, newName, options))
}

// AlterColumn from this table.
func (at *AlterTable) AlterColumn(name string, typ ColumnType, options ...ColumnOption) {
	at.Definitions = append(at.Definitions, alterColumn(name, typ, options))
}

// DropColumn from this table.
func (at *AlterTable) DropColumn(name string, options ...ColumnOption) {
	at.Definitions = append(at.Definitions, dropColumn(name, options))
}

// RenameIndex to a new name.
func (at *AlterTable) RenameIndex(name string, newName string, options ...IndexOption) {
	at.Definitions = append(at.Definitions, renameIndex(name, newName, options))
}

// DropIndex from this table.
func (at *AlterTable) DropIndex(name string, options ...IndexOption) {
	at.Definitions = append(at.Definitions, dropIndex(name, options))
}

func createTable(name string, options []TableOption) Table {
	table := Table{
		Op:   SchemaCreate,
		Name: name,
	}

	applyTableOptions(&table, options)
	return table
}

func alterTable(name string, options []TableOption) AlterTable {
	table := Table{
		Op:   SchemaAlter,
		Name: name,
	}

	applyTableOptions(&table, options)
	return AlterTable{Table: table}
}

func renameTable(name string, newName string, options []TableOption) Table {
	table := Table{
		Op:      SchemaRename,
		Name:    name,
		NewName: newName,
	}

	applyTableOptions(&table, options)
	return table
}

func dropTable(name string, options []TableOption) Table {
	table := Table{
		Op:   SchemaDrop,
		Name: name,
	}

	applyTableOptions(&table, options)
	return table
}

// TableOption interface.
// Available options are: Comment, Options.
type TableOption interface {
	applyTable(table *Table)
}

func applyTableOptions(table *Table, options []TableOption) {
	for i := range options {
		options[i].applyTable(table)
	}
}

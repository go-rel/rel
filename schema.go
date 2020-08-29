package rel

// SchemaOp type.
type SchemaOp uint8

const (
	// SchemaAdd operation.
	SchemaAdd SchemaOp = iota
	// SchemaAlter operation.
	SchemaAlter
	// SchemaRename operation.
	SchemaRename
	// SchemaDrop operation.
	SchemaDrop
)

// Migrator private interface.
type Migrator interface {
	migrate()
}

// Schema builder.
type Schema struct {
	Pending []Migrator
}

func (s *Schema) add(migrator Migrator) {
	s.Pending = append(s.Pending, migrator)
}

// CreateTable with name and its definition.
func (s *Schema) CreateTable(name string, fn func(t *Table), options ...TableOption) {
	table := createTable(name, options)
	fn(&table)
	s.add(table)
}

// AlterTable with name and its definition.
func (s *Schema) AlterTable(name string, fn func(t *AlterTable), options ...TableOption) {
	table := alterTable(name, options)
	fn(&table)
	s.add(table.Table)
}

// RenameTable by name.
func (s *Schema) RenameTable(name string, newName string, options ...TableOption) {
	s.add(renameTable(name, newName, options))
}

// DropTable by name.
func (s *Schema) DropTable(name string, options ...TableOption) {
	s.add(dropTable(name, options))
}

// AddColumn with name and type.
func (s *Schema) AddColumn(table string, name string, typ ColumnType, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.Column(name, typ, options...)
	s.add(at.Table)
}

// AlterColumn by name.
func (s *Schema) AlterColumn(table string, name string, typ ColumnType, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.AlterColumn(name, typ, options...)
	s.add(at.Table)
}

// RenameColumn by name.
func (s *Schema) RenameColumn(table string, name string, newName string, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.RenameColumn(name, newName, options...)
	s.add(at.Table)
}

// DropColumn by name.
func (s *Schema) DropColumn(table string, name string, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.DropColumn(name, options...)
	s.add(at.Table)
}

// AddIndex for columns.
func (s *Schema) AddIndex(table string, column []string, typ IndexType, options ...IndexOption) {
	at := alterTable(table, nil)
	at.Index(column, typ, options...)
	s.add(at.Table)
}

// RenameIndex by name.
func (s *Schema) RenameIndex(table string, name string, newName string, options ...IndexOption) {
	at := alterTable(table, nil)
	at.RenameIndex(name, newName, options...)
	s.add(at.Table)
}

// DropIndex by name.
func (s *Schema) DropIndex(table string, name string, options ...IndexOption) {
	at := alterTable(table, nil)
	at.DropIndex(name, options...)
	s.add(at.Table)
}

// Exec queries using repo.
// Useful for data migration.
// func (s *Schema) Exec(func(repo rel.Repository) error) {
// }

// Comment options for table, column and index.
type Comment string

func (c Comment) applyTable(table *Table) {
	table.Comment = string(c)
}

func (c Comment) applyColumn(column *Column) {
	column.Comment = string(c)
}

func (c Comment) applyIndex(index *Index) {
	index.Comment = string(c)
}

// Options options for table, column and index.
type Options string

func (o Options) applyTable(table *Table) {
	table.Options = string(o)
}

func (o Options) applyColumn(column *Column) {
	column.Options = string(o)
}

func (o Options) applyIndex(index *Index) {
	index.Options = string(o)
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

// Name option for defining custom index name.
type Name string

func (n Name) applyIndex(index *Index) {
	index.Name = string(n)
}

// OnDelete option for foreign key index.
type OnDelete string

func (od OnDelete) applyIndex(index *Index) {
	index.Reference.OnDelete = string(od)
}

// OnUpdate option for foreign key index.
type OnUpdate string

func (ou OnUpdate) applyIndex(index *Index) {
	index.Reference.OnUpdate = string(ou)
}

// IfNotExists option for table creation.
type IfNotExists bool

func (ine IfNotExists) applyTable(table *Table) {
	table.IfNotExists = bool(ine)
}

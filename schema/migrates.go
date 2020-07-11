package schema

// Migrator private interface.
type Migrator interface {
	migrate()
}

// Migrates builder.
type Migrates []Migrator

func (m *Migrates) add(migrator Migrator) {
	*m = append(*m, migrator)
}

// CreateTable with name and its definition.
func (m *Migrates) CreateTable(name string, fn func(t *Table), options ...TableOption) {
	table := createTable(name, options)
	fn(&table)
	m.add(table)
}

// AlterTable with name and its definition.
func (m *Migrates) AlterTable(name string, fn func(t *AlterTable), options ...TableOption) {
	table := alterTable(name, options)
	fn(&table)
	m.add(table.Table)
}

// RenameTable by name.
func (m *Migrates) RenameTable(name string, newName string, options ...TableOption) {
	m.add(renameTable(name, newName, options))
}

// DropTable by name.
func (m *Migrates) DropTable(name string, options ...TableOption) {
	m.add(dropTable(name, options))
}

// AddColumn with name and type.
func (m *Migrates) AddColumn(table string, name string, typ ColumnType, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.Column(name, typ, options...)
	m.add(at.Table)
}

// AlterColumn by name.
func (m *Migrates) AlterColumn(table string, name string, typ ColumnType, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.AlterColumn(name, typ, options...)
	m.add(at.Table)
}

// RenameColumn by name.
func (m *Migrates) RenameColumn(table string, name string, newName string, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.RenameColumn(name, newName, options...)
	m.add(at.Table)
}

// DropColumn by name.
func (m *Migrates) DropColumn(table string, name string, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.DropColumn(name, options...)
	m.add(at.Table)
}

// AddIndex for columns.
func (m *Migrates) AddIndex(table string, column []string, typ IndexType, options ...IndexOption) {
	at := alterTable(table, nil)
	at.Index(column, typ, options...)
	m.add(at.Table)
}

// RenameIndex by name.
func (m *Migrates) RenameIndex(table string, name string, newName string, options ...IndexOption) {
	at := alterTable(table, nil)
	at.RenameIndex(name, newName, options...)
	m.add(at.Table)
}

// DropIndex by name.
func (m *Migrates) DropIndex(table string, name string, options ...IndexOption) {
	at := alterTable(table, nil)
	at.DropIndex(name, options...)
	m.add(at.Table)
}

// Exec queries using repo.
// Useful for data migration.
// func (m *Migrates) Exec(func(repo rel.Repository) error) {
// }

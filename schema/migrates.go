package schema

import "github.com/Fs02/rel"

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
func (m *Migrates) CreateTable(name string, fn func(t *Table), options ...Option) {
	table := Table{Name: name}
	fn(&table)
	m.add(CreateTable(table))
}

// AlterTable with name and its definition.
func (m *Migrates) AlterTable(name string, fn func(t *AlterTable), options ...Option) {
	table := AlterTable{Table: Table{Name: name}}
	fn(&table)
	m.add(table)
}

// RenameTable by name.
func (m *Migrates) RenameTable(name string, newName string) {
	m.add(RenameTable{Name: name, NewName: newName})
}

// DropTable by name.
func (m *Migrates) DropTable(name string) {
	m.add(DropTable{Name: name})
}

// AddColumn with name and type.
func (m *Migrates) AddColumn(table string, name string, typ ColumnType, options ...Option) {
	at := AlterTable{Table: Table{Name: name}}
	at.Column(name, typ, options...)
	m.add(at)
}

// RenameColumn by name.
func (m *Migrates) RenameColumn(table string, name string, newName string, options ...Option) {
	at := AlterTable{Table: Table{Name: name}}
	at.RenameColumn(name, newName, options...)
	m.add(at)
}

// DropColumn by name.
func (m *Migrates) DropColumn(table string, name string, options ...Option) {
	at := AlterTable{Table: Table{Name: name}}
	at.DropColumn(name, options...)
	m.add(at)
}

// AddIndex for columns.
func (m *Migrates) AddIndex(table string, column []string, options ...Option) {
}

// RenameIndex by name.
func (m *Migrates) RenameIndex(table string, name string, newName string) {
}

// DropIndex by name.
func (m *Migrates) DropIndex(table string, name string) {
}

// Exec queries using repo.
// Useful for data migration.
func (m *Migrates) Exec(func(repo rel.Repository) error) {
}

package rel

import (
	"context"
	"strings"
)

// SchemaOp type.
type SchemaOp uint8

const (
	// SchemaCreate operation.
	SchemaCreate SchemaOp = iota
	// SchemaAlter operation.
	SchemaAlter
	// SchemaRename operation.
	SchemaRename
	// SchemaDrop operation.
	SchemaDrop
)

func (s SchemaOp) String() string {
	return [...]string{"create", "alter", "rename", "drop"}[s]
}

// Migration definition.
type Migration interface {
	internalMigration()
	description() string
}

// Schema builder.
type Schema struct {
	Migrations []Migration
}

func (s *Schema) add(migration Migration) {
	s.Migrations = append(s.Migrations, migration)
}

// CreateTable with name and its definition.
func (s *Schema) CreateTable(name string, fn func(t *Table), options ...TableOption) {
	table := createTable(name, options)
	fn(&table)
	s.add(table)
}

// CreateTableIfNotExists with name and its definition.
func (s *Schema) CreateTableIfNotExists(name string, fn func(t *Table), options ...TableOption) {
	table := createTableIfNotExists(name, options)
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

// DropTableIfExists by name.
func (s *Schema) DropTableIfExists(name string, options ...TableOption) {
	s.add(dropTableIfExists(name, options))
}

// AddColumn with name and type.
func (s *Schema) AddColumn(table string, name string, typ ColumnType, options ...ColumnOption) {
	at := alterTable(table, nil)
	at.Column(name, typ, options...)
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

// CreateIndex for columns on a table.
func (s *Schema) CreateIndex(table string, name string, column []string, options ...IndexOption) {
	s.add(createIndex(table, name, column, options))
}

// CreateUniqueIndex for columns on a table.
func (s *Schema) CreateUniqueIndex(table string, name string, column []string, options ...IndexOption) {
	s.add(createUniqueIndex(table, name, column, options))
}

// DropIndex by name.
func (s *Schema) DropIndex(table string, name string, options ...IndexOption) {
	s.add(dropIndex(table, name, options))
}

// Exec queries.
func (s *Schema) Exec(raw Raw) {
	s.add(raw)
}

// Do migration using golang codes.
func (s *Schema) Do(fn Do) {
	s.add(fn)
}

// String returns schema operation.
func (s Schema) String() string {
	descs := make([]string, len(s.Migrations))
	for i := range descs {
		descs[i] = s.Migrations[i].description()
	}

	return strings.Join(descs, ", ")
}

// Raw string
type Raw string

func (r Raw) description() string {
	return "execute raw command"
}

func (r Raw) internalMigration()       {}
func (r Raw) internalTableDefinition() {}

// Do used internally for schema migration.
type Do func(context.Context, Repository) error

func (d Do) description() string {
	return "run go code"
}

func (d Do) internalMigration() {}

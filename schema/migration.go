package schema

import (
	"github.com/Fs02/rel"
)

// Migrater private interface.
type Migrater interface {
	migrate()
}

// Migrate builder.
type Migrate []Migrater

// CreateTable with name and its definition.
func (m *Migrate) CreateTable(name string, fn func(t *Table), options ...Option) {
	table := Table{Name: name} // TODO: use actual migrater
	fn(&table)
	*m = append(*m, table)
}

// RenameTable by name.
func (m *Migrate) RenameTable(name string, newname string) {
}

// DropTable by name.
func (m *Migrate) DropTable(name string) {
}

// AddColumn with name and type.
func (m *Migrate) AddColumn(table string, name string, typ ColumnType, options ...Option) {
}

// RenameColumn by name.
func (m *Migrate) RenameColumn(table string, name string, newname string) {
}

// DropColumn by name.
func (m *Migrate) DropColumn(table string, name string) {
}

// AddIndex for columns.
func (m *Migrate) AddIndex(table string, column []string, options ...Option) {
}

// RenameIndex by name.
func (m *Migrate) RenameIndex(table string, name string, newname string) {
}

// DropIndex by name.
func (m *Migrate) DropIndex(table string, name string) {
}

// Exec queries using repo.
// Useful for data migration.
func (m *Migrate) Exec(func(repo rel.Repository) error) {
}

// Migration definition.
type Migration struct {
	Version int
	Ups     Migrate
	Downs   Migrate
}

// Up migration.
func (m *Migration) Up(fn func(migrate *Migrate)) {
	fn(&m.Ups)
}

// Down migration.
func (m *Migration) Down(fn func(migrate *Migrate)) {
	fn(&m.Downs)
}

// NewMigration for schema.
func NewMigration(version int) Migration {
	return Migration{
		Version: version,
	}
}

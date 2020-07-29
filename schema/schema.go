package schema

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

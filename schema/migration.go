package schema

// Migration definition.
type Migration struct {
	Version int
	Ups     Migrates
	Downs   Migrates
}

// Up migration.
func (m *Migration) Up(fn func(migrate *Migrates)) {
	fn(&m.Ups)
}

// Down migration.
func (m *Migration) Down(fn func(migrate *Migrates)) {
	fn(&m.Downs)
}

// NewMigration for schema.
func NewMigration(version int) Migration {
	return Migration{
		Version: version,
	}
}

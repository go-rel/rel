package schema

import (
	"context"
	"sort"

	"github.com/Fs02/rel"
)

const schemaVersionTable = "rel_schema_versions"

type schemaVersion struct {
	ID      int
	Version int
}

func (schemaVersion) Table() string {
	return schemaVersionTable
}

// MigrationStep definition.
type MigrationStep struct {
	Version int
	Up      Schema
	Down    Schema
	Applied bool
}

// MigrationSteps definition.
type MigrationSteps []MigrationStep

func (ms MigrationSteps) Len() int {
	return len(ms)
}

func (ms MigrationSteps) Less(i, j int) bool {
	return ms[i].Version < ms[j].Version
}

func (ms MigrationSteps) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

// Migration manager definition.
type Migration struct {
	Adapter    Adapter
	Repository rel.Repository
	Steps      MigrationSteps
}

// AddVersion adds a migration with explicit version.
func (m *Migration) AddVersion(version int, up func(schema *Schema), down func(schema *Schema)) {
	var upSchema, downSchema Schema

	up(&upSchema)
	down(&downSchema)

	m.Steps = append(m.Steps, MigrationStep{Version: version, Up: upSchema, Down: downSchema})
}

func (m Migration) buildVersionTableDefinition() Table {
	var versionTable = createTable(schemaVersionTable, nil) // TODO: create if not exists

	versionTable.Int("id")
	versionTable.Int("version")
	versionTable.DateTime("created_at")
	versionTable.DateTime("updated_at")

	versionTable.PrimaryKey("id")
	versionTable.Unique([]string{"version"})

	return versionTable
}

func (m *Migration) sync(ctx context.Context) {
	var (
		versions []schemaVersion
		vi       int
	)

	check(m.Adapter.Apply(ctx, m.buildVersionTableDefinition()))
	m.Repository.MustFindAll(ctx, &versions, rel.NewSortAsc("version"))
	sort.Sort(m.Steps)

	for i := range m.Steps {
		if len(versions) <= vi {
			break
		}

		if m.Steps[i].Version == versions[vi].Version {
			m.Steps[i].Applied = true
			vi++
		}
	}

	if len(versions) <= vi {
		panic("rel: inconsistent schema state")
	}
}

// Up perform migration to the latest schema version.
func (m *Migration) Up(ctx context.Context) {
	m.sync(ctx)

	for _, step := range m.Steps {
		if step.Applied {
			continue
		}

		m.Repository.Transaction(ctx, func(ctx context.Context) error {
			m.Repository.MustInsert(ctx, &schemaVersion{Version: step.Version})

			adapter := m.Repository.Adapter(ctx).(Adapter)
			for _, migrator := range step.Up.Pending {
				switch v := migrator.(type) {
				case Table:
					check(adapter.Apply(ctx, v))
				}
			}

			return nil
		})
	}
}

// Down perform rollback migration 1 step.
func (m *Migration) Down(ctx context.Context) {
	m.sync(ctx)

	for i := range m.Steps {
		step := m.Steps[len(m.Steps)-i-1]
		if !step.Applied {
			continue
		}

		m.Repository.Transaction(ctx, func(ctx context.Context) error {
			m.Repository.MustInsert(ctx, &schemaVersion{Version: step.Version})

			adapter := m.Repository.Adapter(ctx).(Adapter)
			for _, migrator := range step.Down.Pending {
				switch v := migrator.(type) {
				case Table:
					check(adapter.Apply(ctx, v))
				}
			}

			return nil
		})

		return
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

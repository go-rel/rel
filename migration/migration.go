package migration

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

// Step definition.
type Step struct {
	Version int
	Up      rel.Schema
	Down    rel.Schema
	Applied bool
}

// Steps definition.
type Steps []Step

func (s Steps) Len() int {
	return len(s)
}

func (s Steps) Less(i, j int) bool {
	return s[i].Version < s[j].Version
}

func (s Steps) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Migration manager definition.
type Migration struct {
	Repository rel.Repository
	Steps      Steps
}

// AddVersion adds a migration with explicit version.
func (m *Migration) AddVersion(version int, up func(schema *rel.Schema), down func(schema *rel.Schema)) {
	var upSchema, downSchema rel.Schema

	up(&upSchema)
	down(&downSchema)

	m.Steps = append(m.Steps, Step{Version: version, Up: upSchema, Down: downSchema})
}

func (m Migration) buildVersionTableDefinition() rel.Table {
	var schema rel.Schema
	schema.CreateTable(schemaVersionTable, func(t *rel.Table) {
		t.Int("id")
		t.Int("version")
		t.DateTime("created_at")
		t.DateTime("updated_at")

		t.PrimaryKey("id")
		t.Unique([]string{"version"})
	}, rel.IfNotExists(true))

	return schema.Pending[0].(rel.Table)
}

func (m *Migration) sync(ctx context.Context) {
	var (
		versions []schemaVersion
		vi       int
		adapter  = m.Repository.Adapter(ctx).(rel.Adapter)
	)

	check(adapter.Apply(ctx, m.buildVersionTableDefinition()))
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

			adapter := m.Repository.Adapter(ctx).(rel.Adapter)
			for _, migrator := range step.Up.Pending {
				// TODO: exec script
				switch v := migrator.(type) {
				case rel.Table:
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

			adapter := m.Repository.Adapter(ctx).(rel.Adapter)
			for _, migrator := range step.Down.Pending {
				switch v := migrator.(type) {
				case rel.Table:
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

package migrator

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Fs02/rel"
)

const versionTable = "rel_schema_versions"

type version struct {
	ID        int
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time

	up      rel.Schema
	down    rel.Schema
	applied bool
}

func (version) Table() string {
	return versionTable
}

type versions []version

func (v versions) Len() int {
	return len(v)
}

func (v versions) Less(i, j int) bool {
	return v[i].Version < v[j].Version
}

func (v versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// Migrator is a migration manager that handles migration logic.
type Migrator struct {
	repo     rel.Repository
	versions versions
}

// RegisterVersion registers a migration with an explicit version.
func (m *Migrator) RegisterVersion(v int, up func(schema *rel.Schema), down func(schema *rel.Schema)) {
	var upSchema, downSchema rel.Schema

	up(&upSchema)
	down(&downSchema)

	m.versions = append(m.versions, version{Version: v, up: upSchema, down: downSchema})
}

func (m Migrator) buildVersionTableDefinition() rel.Table {
	var schema rel.Schema
	schema.CreateTable(versionTable, func(t *rel.Table) {
		t.Int("id")
		t.Int("version")
		t.DateTime("created_at")
		t.DateTime("updated_at")

		t.PrimaryKey("id")
		t.Unique([]string{"version"})
	}, rel.Optional(true))

	return schema.Pending[0].(rel.Table)
}

func (m *Migrator) sync(ctx context.Context) {
	var (
		versions versions
		vi       int
		adapter  = m.repo.Adapter(ctx).(rel.Adapter)
	)

	check(adapter.Apply(ctx, m.buildVersionTableDefinition()))
	m.repo.MustFindAll(ctx, &versions, rel.NewSortAsc("version"))
	sort.Sort(m.versions)

	for i := range m.versions {
		if vi < len(versions) && m.versions[i].Version == versions[vi].Version {
			m.versions[i].ID = versions[vi].ID
			m.versions[i].applied = true
			vi++
		} else {
			m.versions[i].applied = false
		}
	}

	if vi != len(versions) {
		panic(fmt.Sprint("rel: missing local migration: ", versions[vi].Version))
	}
}

// Migrate to the latest schema version.
func (m *Migrator) Migrate(ctx context.Context) {
	m.sync(ctx)

	for _, step := range m.versions {
		if step.applied {
			continue
		}

		err := m.repo.Transaction(ctx, func(ctx context.Context) error {
			m.repo.MustInsert(ctx, &version{Version: step.Version})

			adapter := m.repo.Adapter(ctx).(rel.Adapter)
			for _, migrator := range step.up.Pending {
				// TODO: exec script
				switch v := migrator.(type) {
				case rel.Table:
					check(adapter.Apply(ctx, v))
				}
			}

			return nil
		})

		check(err)
	}
}

// Rollback migration 1 step.
func (m *Migrator) Rollback(ctx context.Context) {
	m.sync(ctx)

	for i := range m.versions {
		v := m.versions[len(m.versions)-i-1]
		if !v.applied {
			continue
		}

		err := m.repo.Transaction(ctx, func(ctx context.Context) error {
			m.repo.MustDelete(ctx, &v)

			adapter := m.repo.Adapter(ctx).(rel.Adapter)
			for _, migrator := range v.down.Pending {
				switch v := migrator.(type) {
				case rel.Table:
					check(adapter.Apply(ctx, v))
				}
			}

			return nil
		})

		check(err)

		// only rollback one version.
		return
	}
}

// New migrationr.
func New(repo rel.Repository) Migrator {
	return Migrator{repo: repo}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

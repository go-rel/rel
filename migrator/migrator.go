package migrator

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/go-rel/rel"
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
	repo               rel.Repository
	instrumenter       rel.Instrumenter
	versions           versions
	versionTableExists bool
}

// Instrumentation function.
func (m *Migrator) Instrumentation(instrumenter rel.Instrumenter) {
	m.instrumenter = instrumenter
}

// Register a migration.
func (m *Migrator) Register(v int, up func(schema *rel.Schema), down func(schema *rel.Schema)) {
	var upSchema, downSchema rel.Schema

	up(&upSchema)
	down(&downSchema)

	m.versions = append(m.versions, version{Version: v, up: upSchema, down: downSchema})
}

func (m Migrator) buildVersionTableDefinition() rel.Table {
	var schema rel.Schema
	schema.CreateTableIfNotExists(versionTable, func(t *rel.Table) {
		t.ID("id")
		t.BigInt("version", rel.Unsigned(true), rel.Unique(true))
		t.DateTime("created_at")
		t.DateTime("updated_at")
	})

	return schema.Migrations[0].(rel.Table)
}

func (m *Migrator) sync(ctx context.Context) {
	var (
		versions versions
		vi       int
		adapter  = m.repo.Adapter(ctx).(rel.Adapter)
	)

	if !m.versionTableExists {
		check(adapter.Apply(ctx, m.buildVersionTableDefinition()))
		m.versionTableExists = true
	}

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

	for _, v := range m.versions {
		if v.applied {
			continue
		}

		finish := m.instrumenter.Observe(ctx, "migrate", strconv.Itoa(v.Version)+" "+v.up.String())

		err := m.repo.Transaction(ctx, func(ctx context.Context) error {
			m.repo.MustInsert(ctx, &version{Version: v.Version})
			m.run(ctx, v.up.Migrations)
			return nil
		})

		finish(err)
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

		finish := m.instrumenter.Observe(ctx, "rollback", strconv.Itoa(v.Version)+" "+v.down.String())

		err := m.repo.Transaction(ctx, func(ctx context.Context) error {
			m.repo.MustDelete(ctx, &v)
			m.run(ctx, v.down.Migrations)
			return nil
		})

		finish(err)
		check(err)

		// only rollback one version.
		return
	}
}

func (m *Migrator) run(ctx context.Context, migrations []rel.Migration) {
	adapter := m.repo.Adapter(ctx).(rel.Adapter)
	for _, migration := range migrations {
		if fn, ok := migration.(rel.Do); ok {
			check(fn(m.repo))
		} else {
			check(adapter.Apply(ctx, migration))
		}
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

package migrator

import (
	"context"
	"errors"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestMigrator(t *testing.T) {
	var (
		ctx      = context.TODO()
		repo     = reltest.New()
		migrator = New(repo)
	)

	t.Run("Register", func(t *testing.T) {
		migrator.Register(20200829084000,
			func(schema *rel.Schema) {
				schema.CreateTable("users", func(t *rel.Table) {
					t.ID("id")
				})
			},
			func(schema *rel.Schema) {
				schema.DropTable("users")
			},
		)

		migrator.Register(20200828100000,
			func(schema *rel.Schema) {
				schema.CreateTable("tags", func(t *rel.Table) {
					t.ID("id")
				})

				schema.Do(func(repo rel.Repository) error {
					assert.NotNil(t, repo)
					return nil
				})
			},
			func(schema *rel.Schema) {
				schema.DropTable("tags")
			},
		)

		migrator.Register(20200829115100,
			func(schema *rel.Schema) {
				schema.CreateTable("books", func(t *rel.Table) {
					t.ID("id")
				})
			},
			func(schema *rel.Schema) {
				schema.DropTable("books")
			},
		)

		assert.Len(t, migrator.versions, 3)
		assert.Equal(t, 20200829084000, migrator.versions[0].Version)
		assert.Equal(t, 20200828100000, migrator.versions[1].Version)
		assert.Equal(t, 20200829115100, migrator.versions[2].Version)
	})

	t.Run("Migrate", func(t *testing.T) {
		repo.ExpectFindAll(rel.NewSortAsc("version")).
			Result(versions{{ID: 1, Version: 20200829115100}})

		repo.ExpectTransaction(func(repo *reltest.Repository) {
			repo.ExpectInsert().For(&version{Version: 20200828100000})
		})

		repo.ExpectTransaction(func(repo *reltest.Repository) {
			repo.ExpectInsert().For(&version{Version: 20200829084000})
		})

		migrator.Migrate(ctx)
	})

	t.Run("Rollback", func(t *testing.T) {
		repo.ExpectFindAll(rel.NewSortAsc("version")).
			Result(versions{
				{ID: 1, Version: 20200828100000},
				{ID: 2, Version: 20200829084000},
			})

		assert.Equal(t, 20200829084000, migrator.versions[1].Version)

		repo.ExpectTransaction(func(repo *reltest.Repository) {
			repo.ExpectDelete().For(&migrator.versions[1])
		})

		migrator.Rollback(ctx)
	})
}

func TestMigrator_Sync(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
		nfn  = func(schema *rel.Schema) {}
	)

	tests := []struct {
		name    string
		applied versions
		synced  versions
		isPanic bool
	}{
		{
			name: "all migrated",
			applied: versions{
				{ID: 1, Version: 1},
				{ID: 2, Version: 2},
				{ID: 3, Version: 3},
			},
			synced: versions{
				{ID: 1, Version: 1, applied: true},
				{ID: 2, Version: 2, applied: true},
				{ID: 3, Version: 3, applied: true},
			},
		},
		{
			name:    "not migrated",
			applied: versions{},
			synced: versions{
				{ID: 0, Version: 1, applied: false},
				{ID: 0, Version: 2, applied: false},
				{ID: 0, Version: 3, applied: false},
			},
		},
		{
			name: "first not migrated",
			applied: versions{
				{ID: 2, Version: 2},
				{ID: 3, Version: 3},
			},
			synced: versions{
				{ID: 0, Version: 1, applied: false},
				{ID: 2, Version: 2, applied: true},
				{ID: 3, Version: 3, applied: true},
			},
		},
		{
			name: "middle not migrated",
			applied: versions{
				{ID: 1, Version: 1},
				{ID: 3, Version: 3},
			},
			synced: versions{
				{ID: 1, Version: 1, applied: true},
				{ID: 0, Version: 2, applied: false},
				{ID: 3, Version: 3, applied: true},
			},
		},
		{
			name: "last not migrated",
			applied: versions{
				{ID: 1, Version: 1},
				{ID: 2, Version: 2},
			},
			synced: versions{
				{ID: 1, Version: 1, applied: true},
				{ID: 2, Version: 2, applied: true},
				{ID: 0, Version: 3, applied: false},
			},
		},
		{
			name: "broken migration",
			applied: versions{
				{ID: 1, Version: 1},
				{ID: 2, Version: 2},
				{ID: 3, Version: 3},
				{ID: 4, Version: 4},
			},
			synced: versions{
				{ID: 1, Version: 1, applied: true},
				{ID: 2, Version: 2, applied: true},
				{ID: 3, Version: 3, applied: true},
			},
			isPanic: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			migrator := New(repo)
			migrator.Register(3, nfn, nfn)
			migrator.Register(2, nfn, nfn)
			migrator.Register(1, nfn, nfn)

			repo.ExpectFindAll(rel.NewSortAsc("version")).Result(test.applied)

			if test.isPanic {
				assert.Panics(t, func() {
					migrator.sync(ctx)
				})
			} else {
				assert.NotPanics(t, func() {
					migrator.sync(ctx)
				})

				assert.Equal(t, test.synced, migrator.versions)
			}
		})
	}
}

func TestMigrator_Instrumentation(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
		m    = New(repo)
	)

	m.Instrumentation(func(context.Context, string, string) func(error) { return nil })
	m.instrumenter.Observe(ctx, "test", "test")
}

func TestCheck(t *testing.T) {
	assert.Panics(t, func() {
		check(errors.New("error"))
	})
}

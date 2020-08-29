package migration

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestSchema_Migration(t *testing.T) {
	var (
		// ctx       = context.TODO()
		repo      = reltest.New()
		migration = Migration{
			Repository: repo,
		}
	)

	t.Run("AddVersion", func(t *testing.T) {
		migration.AddVersion(20200829114000,
			func(schema *rel.Schema) {
				schema.CreateTable("users", func(t *rel.Table) {
					t.Int("id")
					t.PrimaryKey("id")
				})
			},
			func(schema *rel.Schema) {
				schema.DropTable("users")
			},
		)

		migration.AddVersion(20200828100000,
			func(schema *rel.Schema) {
				schema.CreateTable("tags", func(t *rel.Table) {
					t.Int("id")
					t.PrimaryKey("id")
				})
			},
			func(schema *rel.Schema) {
				schema.DropTable("tags")
			},
		)

		migration.AddVersion(20200829115100,
			func(schema *rel.Schema) {
				schema.CreateTable("books", func(t *rel.Table) {
					t.Int("id")
					t.PrimaryKey("id")
				})
			},
			func(schema *rel.Schema) {
				schema.DropTable("books")
			},
		)

		assert.Len(t, migration.Steps, 3)
		assert.Equal(t, 20200829114000, migration.Steps[0].Version)
		assert.Equal(t, 20200828100000, migration.Steps[1].Version)
		assert.Equal(t, 20200829115100, migration.Steps[2].Version)
	})

	// t.Run("Up", func(t *testing.T) {
	// 	migration.Up(ctx)
	// })
}

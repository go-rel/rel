package specs

import (
	"context"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/migrator"
)

// Migrate database for specs execution.
func Migrate(ctx context.Context, adapter rel.Adapter, rollback bool) {
	m := migrator.New(rel.New(adapter))

	m.RegisterVersion(1,
		func(schema *rel.Schema) {
			schema.CreateTable("users", func(t *rel.Table) {
				t.Int("id")
				t.String("slug", rel.Limit(30))
				t.String("name", rel.Limit(30), rel.Default(""))
				t.String("gender", rel.Limit(10), rel.Default(""))
				t.Int("age", rel.Required(true), rel.Default(0))
				t.String("note", rel.Limit(50))
				t.DateTime("created_at")
				t.DateTime("updated_at")

				t.PrimaryKey("id")
				t.Unique([]string{"slug"})
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("users")
		},
	)

	m.RegisterVersion(2,
		func(schema *rel.Schema) {
			schema.CreateTable("addresses", func(t *rel.Table) {
				t.Int("id")
				t.Int("user_id")
				t.String("name", rel.Limit(60), rel.Required(true), rel.Default(""))
				t.DateTime("created_at")
				t.DateTime("updated_at")

				t.PrimaryKey("id")
				t.ForeignKey("user_id", "users", "id")
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("addresses")
		},
	)

	m.RegisterVersion(3,
		func(schema *rel.Schema) {
			schema.CreateTable("extras", func(t *rel.Table) {
				t.Int("id")
				t.Int("user_id")
				t.String("slug", rel.Limit(30))
				t.Int("score", rel.Default(0))

				t.PrimaryKey("id")
				t.ForeignKey("user_id", "users", "id")
				t.Unique([]string{"slug"})
				t.Fragment("CONSTRAINT extras_score_check CHECK (score>=0 AND score<=100)")
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("extras")
		},
	)

	m.RegisterVersion(4,
		func(schema *rel.Schema) {
			schema.CreateTable("composites", func(t *rel.Table) {
				t.Int("primary1")
				t.Int("primary2")
				t.String("data")

				t.Index([]string{"primary1", "primary2"}, rel.PrimaryKey)
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("composites")
		},
	)

	if rollback {
		for i := 0; i < 4; i++ {
			m.Rollback(ctx)
		}
	} else {
		m.Migrate(ctx)
	}
}

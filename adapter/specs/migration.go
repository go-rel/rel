package specs

import (
	"testing"
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/migrator"
)

var m migrator.Migrator

// Migrate database for specs execution.
func Migrate(t *testing.T, repo rel.Repository, rollback bool) {
	if rollback {
		for i := 0; i < 4; i++ {
			m.Rollback(ctx)
		}
		return
	}

	m = migrator.New(repo)
	m.Register(1,
		func(schema *rel.Schema) {
			schema.CreateTable("users", func(t *rel.Table) {
				t.ID("id")
				t.String("slug", rel.Limit(30))
				t.String("name", rel.Limit(30), rel.Default(""))
				t.String("gender", rel.Limit(10), rel.Default(""))
				t.Int("age", rel.Required(true), rel.Default(0))
				t.String("note", rel.Limit(50))
				t.DateTime("created_at")
				t.DateTime("updated_at")

				t.Unique([]string{"slug"})
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("users")
		},
	)

	m.Register(2,
		func(schema *rel.Schema) {
			schema.CreateTable("addresses", func(t *rel.Table) {
				t.ID("id")
				t.Int("user_id", rel.Unsigned(true))
				t.String("name", rel.Limit(60), rel.Required(true), rel.Default(""))
				t.DateTime("created_at")
				t.DateTime("updated_at")

				t.ForeignKey("user_id", "users", "id")
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("addresses")
		},
	)

	m.Register(3,
		func(schema *rel.Schema) {
			schema.CreateTable("extras", func(t *rel.Table) {
				t.ID("id")
				t.Int("user_id", rel.Unsigned(true))
				t.String("slug", rel.Limit(30))
				t.Int("score", rel.Default(0))

				t.ForeignKey("user_id", "users", "id")
				t.Unique([]string{"slug"})
				t.Fragment("CONSTRAINT extras_score_check CHECK (score>=0 AND score<=100)")
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("extras")
		},
	)

	m.Register(4,
		func(schema *rel.Schema) {
			schema.CreateTable("composites", func(t *rel.Table) {
				t.Int("primary1")
				t.Int("primary2")
				t.String("data")

				t.PrimaryKey([]string{"primary1", "primary2"})
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("composites")
		},
	)

	m.Migrate(ctx)
}

// MigrateTable specs.
func MigrateTable(t *testing.T, repo rel.Repository) {
	m.Register(5,
		func(schema *rel.Schema) {
			schema.CreateTable("dummies", func(t *rel.Table) {
				t.ID("id")
				t.Bool("bool1")
				t.Bool("bool2", rel.Default(true))
				t.Int("int1")
				t.Int("int2", rel.Default(8), rel.Unsigned(true), rel.Limit(10))
				t.BigInt("bigint1")
				t.BigInt("bigint2", rel.Default(8), rel.Unsigned(true), rel.Limit(200))
				t.Float("float1")
				t.Float("float2", rel.Default(10.00), rel.Precision(2))
				t.Decimal("decimal1")
				t.Decimal("decimal2", rel.Default(10.00), rel.Precision(6), rel.Scale(2))
				t.String("string1")
				t.String("string2", rel.Default("string"), rel.Limit(100))
				t.Text("text")
				t.Date("date1")
				t.Date("date2", rel.Default(time.Now()))
				t.DateTime("datetime1")
				t.DateTime("datetime2", rel.Default(time.Now()))
				t.Time("time1")
				t.Time("time2", rel.Default(time.Now()))
				t.Timestamp("timestamp1")
				t.Timestamp("timestamp2", rel.Default(time.Now()))

				t.Unique([]string{"int1"})
				t.Unique([]string{"bigint1", "bigint2"})
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("dummies")
		},
	)

	m.Migrate(ctx)
	m.Rollback(ctx)
}

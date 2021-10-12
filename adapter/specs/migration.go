package specs

import (
	"testing"
	"time"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/migrator"
)

var m migrator.Migrator

// Setup database for specs execution.
func Setup(t *testing.T, repo rel.Repository) func() {
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
			})

			schema.CreateUniqueIndex("users", "unique_slug", []string{"slug"})
		},
		func(schema *rel.Schema) {
			schema.DropIndex("users", "unique_slug")
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

				t.PrimaryKeys([]string{"primary1", "primary2"})
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("composites")
		},
	)

	m.Migrate(ctx)

	return func() {
		for i := 0; i < 4; i++ {
			m.Rollback(ctx)
		}
	}
}

// Migrate specs.
func Migrate(t *testing.T, repo rel.Repository, flags ...Flag) {
	m.Register(5,
		func(schema *rel.Schema) {
			schema.CreateTable("dummies", func(t *rel.Table) {
				t.BigID("id")
				t.Bool("bool1")
				t.Bool("bool2", rel.Default(true))
				t.Int("int1")
				t.Int("int2", rel.Default(8), rel.Unsigned(true), rel.Limit(10))
				t.Int("int3", rel.Unique(true))
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

				t.Unique([]string{"int2"})
				t.Unique([]string{"bigint1", "bigint2"})
			})
		},
		func(schema *rel.Schema) {
			schema.DropTable("dummies")
		},
	)
	defer m.Rollback(ctx)

	m.Register(6,
		func(schema *rel.Schema) {
			schema.AlterTable("dummies", func(t *rel.AlterTable) {
				t.Bool("new_column")
			})
			schema.AddColumn("dummies", "new_column1", rel.Int, rel.Unsigned(true))
		},
		func(schema *rel.Schema) {
			if SkipDropColumn.disabled(flags) {
				schema.AlterTable("dummies", func(t *rel.AlterTable) {
					t.DropColumn("new_column")
				})
				schema.DropColumn("dummies", "new_column1")
			}
		},
	)
	defer m.Rollback(ctx)

	if SkipRenameColumn.disabled(flags) {
		m.Register(7,
			func(schema *rel.Schema) {
				schema.AlterTable("dummies", func(t *rel.AlterTable) {
					t.RenameColumn("text", "teks")
					t.RenameColumn("date2", "date3")
				})
				schema.RenameColumn("dummies", "decimal1", "decimal0")
			},
			func(schema *rel.Schema) {
				schema.AlterTable("dummies", func(t *rel.AlterTable) {
					t.RenameColumn("teks", "text")
					t.RenameColumn("date3", "date2")
				})
				schema.RenameColumn("dummies", "decimal0", "decimal1")
			},
		)
		defer m.Rollback(ctx)
	}

	m.Register(8,
		func(schema *rel.Schema) {
			schema.CreateIndex("dummies", "int1_idx", []string{"int1"})
			schema.CreateIndex("dummies", "string1_string2_idx", []string{"string1", "string2"})
		},
		func(schema *rel.Schema) {
			schema.DropIndex("dummies", "int1_idx")
			schema.DropIndex("dummies", "string1_string2_idx")
		},
	)
	defer m.Rollback(ctx)

	m.Register(9,
		func(schema *rel.Schema) {
			schema.RenameTable("dummies", "new_dummies")
		},
		func(schema *rel.Schema) {
			schema.RenameTable("new_dummies", "dummies")
		},
	)
	defer m.Rollback(ctx)

	m.Register(10,
		func(schema *rel.Schema) {
			schema.CreateTableIfNotExists("dummies2", func(t *rel.Table) {
				t.ID("id")
			})
		},
		func(schema *rel.Schema) {
			schema.DropTableIfExists("dummies2")
		},
	)
	defer m.Rollback(ctx)

	m.Register(11,
		func(schema *rel.Schema) {
			schema.CreateTableIfNotExists("dummies2", func(t *rel.Table) {
				t.ID("id")
				t.Int("field1")
				t.Int("field2")
			})
		},
		func(schema *rel.Schema) {
			schema.DropTableIfExists("dummies2")
		},
	)
	defer m.Rollback(ctx)

	m.Register(12,
		func(schema *rel.Schema) {
			tm := time.Now()
			schema.CreateTableIfNotExists("dummies3", func(t *rel.Table) {
				t.ID("id")
				t.Int("field1")
				t.Float("field2", rel.Default(float32(1.337)))
				t.DateTime("created_at", rel.Default(tm))
				t.DateTime("updated_at", rel.Default(tm))
				t.Bool("is_active", rel.Default(true))
			})
			schema.CreateUniqueIndex("dummies3", "dummies3_field1_active_uq", []string{"field1"}, rel.Eq("is_active", true))
			schema.CreateUniqueIndex("dummies3", "dummies3_field2_uq", []string{"field2"}, rel.Gt("field1", 12))
			schema.CreateUniqueIndex("dummies3", "dummies3_field1_time_uq", []string{"field2"}, rel.Gt("created_at", &tm))
		},
		func(schema *rel.Schema) {
			schema.DropIndex("dummies3", "dummies3_field1_time_uq")
			schema.DropIndex("dummies3", "dummies3_field2_uq")
			schema.DropIndex("dummies3", "dummies3_field1_active_uq")
			schema.DropTableIfExists("dummies3")
		},
	)
	defer m.Rollback(ctx)

	m.Register(13,
		func(schema *rel.Schema) {
			schema.CreateTableIfNotExists("options", func(t *rel.Table) {
				t.ID("id")
				t.String("name")
				t.JSON("value")
			})
		},
		func(schema *rel.Schema) {
			schema.DropTableIfExists("options")
		},
	)
	defer m.Rollback(ctx)

	m.Migrate(ctx)
}

package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaOp(t *testing.T) {
	ops := map[string]SchemaOp{
		"create": SchemaCreate,
		"alter":  SchemaAlter,
		"rename": SchemaRename,
		"drop":   SchemaDrop,
	}

	for name, op := range ops {
		assert.Equal(t, name, op.String())
	}
}

func TestSchema_CreateTable(t *testing.T) {
	var schema Schema

	schema.CreateTable("products", func(t *Table) {
		t.ID("id")
		t.String("name")
		t.Text("description")
	})

	assert.Equal(t, Table{
		Op:   SchemaCreate,
		Name: "products",
		Definitions: []TableDefinition{
			Column{Name: "id", Type: ID, Primary: true},
			Column{Name: "name", Type: String},
			Column{Name: "description", Type: Text},
		},
	}, schema.Migrations[0])

	schema.CreateTableIfNotExists("wishlists", func(t *Table) {
		t.ID("id")
	})

	assert.Equal(t, Table{
		Op:       SchemaCreate,
		Name:     "wishlists",
		Optional: true,
		Definitions: []TableDefinition{
			Column{Name: "id", Type: ID, Primary: true},
		},
	}, schema.Migrations[1])

	assert.Equal(t, "create table products, create table wishlists", schema.String())
}

func TestSchema_AlterTable(t *testing.T) {
	var schema Schema

	schema.AlterTable("users", func(t *AlterTable) {
		t.Bool("verified")
		t.RenameColumn("name", "fullname")
	})

	assert.Equal(t, Table{
		Op:   SchemaAlter,
		Name: "users",
		Definitions: []TableDefinition{
			Column{Name: "verified", Type: Bool, Op: SchemaCreate},
			Column{Name: "name", Rename: "fullname", Op: SchemaRename},
		},
	}, schema.Migrations[0])
}

func TestSchema_RenameTable(t *testing.T) {
	var schema Schema

	schema.RenameTable("trxs", "transactions")

	assert.Equal(t, Table{
		Op:     SchemaRename,
		Name:   "trxs",
		Rename: "transactions",
	}, schema.Migrations[0])
}

func TestSchema_DropTable(t *testing.T) {
	var schema Schema

	schema.DropTable("logs")

	assert.Equal(t, Table{
		Op:   SchemaDrop,
		Name: "logs",
	}, schema.Migrations[0])

	schema.DropTableIfExists("logs")

	assert.Equal(t, Table{
		Op:       SchemaDrop,
		Name:     "logs",
		Optional: true,
	}, schema.Migrations[1])
}

func TestSchema_AddColumn(t *testing.T) {
	var schema Schema

	schema.AddColumn("products", "description", String)

	assert.Equal(t, Table{
		Op:   SchemaAlter,
		Name: "products",
		Definitions: []TableDefinition{
			Column{Name: "description", Type: String, Op: SchemaCreate},
		},
	}, schema.Migrations[0])
}

func TestSchema_RenameColumn(t *testing.T) {
	var schema Schema

	schema.RenameColumn("users", "name", "fullname")

	assert.Equal(t, Table{
		Op:   SchemaAlter,
		Name: "users",
		Definitions: []TableDefinition{
			Column{Name: "name", Rename: "fullname", Op: SchemaRename},
		},
	}, schema.Migrations[0])
}

func TestSchema_DropColumn(t *testing.T) {
	var schema Schema

	schema.DropColumn("users", "verified")

	assert.Equal(t, Table{
		Op:   SchemaAlter,
		Name: "users",
		Definitions: []TableDefinition{
			Column{Name: "verified", Op: SchemaDrop},
		},
	}, schema.Migrations[0])
}

func TestSchema_CreateIndex(t *testing.T) {
	var schema Schema

	schema.CreateIndex("products", "sale_idx", []string{"sale"})

	assert.Equal(t, Index{
		Table:   "products",
		Name:    "sale_idx",
		Columns: []string{"sale"},
		Op:      SchemaCreate,
	}, schema.Migrations[0])
}

func TestSchema_CreateIndex_unique(t *testing.T) {
	var schema Schema

	schema.CreateIndex("products", "sale_idx", []string{"sale"}, Unique(true))

	assert.Equal(t, Index{
		Table:   "products",
		Name:    "sale_idx",
		Unique:  true,
		Columns: []string{"sale"},
		Op:      SchemaCreate,
	}, schema.Migrations[0])
}

func TestSchema_CreateUniqueIndex(t *testing.T) {
	var schema Schema

	schema.CreateUniqueIndex("products", "sale_idx", []string{"sale"})
	assert.Equal(t, Index{
		Table:   "products",
		Name:    "sale_idx",
		Unique:  true,
		Columns: []string{"sale"},
		Op:      SchemaCreate,
	}, schema.Migrations[0])
}

func TestSchema_DropIndex(t *testing.T) {
	var schema Schema

	schema.DropIndex("products", "sale")

	assert.Equal(t, Index{
		Table: "products",
		Name:  "sale",
		Op:    SchemaDrop,
	}, schema.Migrations[0])
}

func TestRaw(t *testing.T) {
	var schema Schema

	schema.Exec("RAW SQL")
	assert.Equal(t, Raw("RAW SQL"), schema.Migrations[0])
}

func TestRaw_Description(t *testing.T) {
	assert.Equal(t, "execute raw command", Raw("").description())
}

func TestRaw_InternalMigration(t *testing.T) {
	assert.NotPanics(t, func() { Raw("").internalMigration() })
}

func TestRaw_InternalTableDefinition(t *testing.T) {
	assert.NotPanics(t, func() { Raw("").internalTableDefinition() })
}

func TestDo(t *testing.T) {
	var (
		schema Schema
	)

	schema.Do(func(repo Repository) error { return nil })
	assert.NotNil(t, schema.Migrations[0])
}

func TestDo_InternalTableDefinition(t *testing.T) {
	assert.NotPanics(t, func() { Do(nil).internalMigration() })
}

func TestDo_Description(t *testing.T) {
	assert.Equal(t, "run go code", Do(nil).description())
}

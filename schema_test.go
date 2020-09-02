package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			Column{Name: "id", Type: ID},
			Column{Name: "name", Type: String},
			Column{Name: "description", Type: Text},
		},
	}, schema.Migrations[0])

	schema.CreateTableIfNotExists("products", func(t *Table) {
		t.ID("id")
	})

	assert.Equal(t, Table{
		Op:       SchemaCreate,
		Name:     "products",
		Optional: true,
		Definitions: []TableDefinition{
			Column{Name: "id", Type: ID},
		},
	}, schema.Migrations[1])
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

func TestSchema_Exec(t *testing.T) {
	var schema Schema

	schema.Exec("RAW SQL")
	assert.Equal(t, Raw("RAW SQL"), schema.Migrations[0])
}

func TestSchema_Do(t *testing.T) {
	var (
		schema Schema
	)

	schema.Do(func(repo Repository) error { return nil })
	assert.NotNil(t, schema.Migrations[0])
}

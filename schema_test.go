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
	}, schema.Migration[0])
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
			Column{Name: "name", NewName: "fullname", Op: SchemaRename},
		},
	}, schema.Migration[0])
}

func TestSchema_RenameTable(t *testing.T) {
	var schema Schema

	schema.RenameTable("trxs", "transactions")

	assert.Equal(t, Table{
		Op:      SchemaRename,
		Name:    "trxs",
		NewName: "transactions",
	}, schema.Migration[0])
}

func TestSchema_DropTable(t *testing.T) {
	var schema Schema

	schema.DropTable("logs")

	assert.Equal(t, Table{
		Op:   SchemaDrop,
		Name: "logs",
	}, schema.Migration[0])
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
	}, schema.Migration[0])
}

func TestSchema_AlterColumn(t *testing.T) {
	var schema Schema

	schema.AlterColumn("products", "sale", Bool)

	assert.Equal(t, Table{
		Op:   SchemaAlter,
		Name: "products",
		Definitions: []TableDefinition{
			Column{Name: "sale", Type: Bool, Op: SchemaAlter},
		},
	}, schema.Migration[0])
}

func TestSchema_RenameColumn(t *testing.T) {
	var schema Schema

	schema.RenameColumn("users", "name", "fullname")

	assert.Equal(t, Table{
		Op:   SchemaAlter,
		Name: "users",
		Definitions: []TableDefinition{
			Column{Name: "name", NewName: "fullname", Op: SchemaRename},
		},
	}, schema.Migration[0])
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
	}, schema.Migration[0])
}

func TestSchema_CreateIndex(t *testing.T) {
	var schema Schema

	schema.CreateIndex("products", []string{"sale"}, SimpleIndex)

	assert.Equal(t, Index{
		Table:   "products",
		Columns: []string{"sale"},
		Type:    SimpleIndex,
		Op:      SchemaCreate,
	}, schema.Migration[0])
}

func TestSchema_DropIndex(t *testing.T) {
	var schema Schema

	schema.DropIndex("products", "sale")

	assert.Equal(t, Index{
		Table: "products",
		Name:  "sale",
		Op:    SchemaDrop,
	}, schema.Migration[0])
}

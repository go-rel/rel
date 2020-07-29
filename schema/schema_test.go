package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema_CreateTable(t *testing.T) {
	var schema Schema

	schema.CreateTable("products", func(t *Table) {
		t.Int("id")
		t.String("name")
		t.Text("description")

		t.PrimaryKey("id")
	})

	assert.Equal(t, Table{
		Op:   Add,
		Name: "products",
		Definitions: []interface{}{
			Column{Name: "id", Type: Int},
			Column{Name: "name", Type: String},
			Column{Name: "description", Type: Text},
			Index{Columns: []string{"id"}, Type: PrimaryKey},
		},
	}, schema.Pending[0])
}

func TestSchema_AlterTable(t *testing.T) {
	var schema Schema

	schema.AlterTable("users", func(t *AlterTable) {
		t.Bool("verified")
		t.RenameColumn("name", "fullname")
	})

	assert.Equal(t, Table{
		Op:   Alter,
		Name: "users",
		Definitions: []interface{}{
			Column{Name: "verified", Type: Bool, Op: Add},
			Column{Name: "name", NewName: "fullname", Op: Rename},
		},
	}, schema.Pending[0])
}

func TestSchema_RenameTable(t *testing.T) {
	var schema Schema

	schema.RenameTable("trxs", "transactions")

	assert.Equal(t, Table{
		Op:      Rename,
		Name:    "trxs",
		NewName: "transactions",
	}, schema.Pending[0])
}

func TestSchema_DropTable(t *testing.T) {
	var schema Schema

	schema.DropTable("logs")

	assert.Equal(t, Table{
		Op:   Drop,
		Name: "logs",
	}, schema.Pending[0])
}

func TestSchema_AddColumn(t *testing.T) {
	var schema Schema

	schema.AddColumn("products", "description", String)

	assert.Equal(t, Table{
		Op:   Alter,
		Name: "products",
		Definitions: []interface{}{
			Column{Name: "description", Type: String, Op: Add},
		},
	}, schema.Pending[0])
}

func TestSchema_AlterColumn(t *testing.T) {
	var schema Schema

	schema.AlterColumn("products", "sale", Bool)

	assert.Equal(t, Table{
		Op:   Alter,
		Name: "products",
		Definitions: []interface{}{
			Column{Name: "sale", Type: Bool, Op: Alter},
		},
	}, schema.Pending[0])
}

func TestSchema_RenameColumn(t *testing.T) {
	var schema Schema

	schema.RenameColumn("users", "name", "fullname")

	assert.Equal(t, Table{
		Op:   Alter,
		Name: "users",
		Definitions: []interface{}{
			Column{Name: "name", NewName: "fullname", Op: Rename},
		},
	}, schema.Pending[0])
}

func TestSchema_DropColumn(t *testing.T) {
	var schema Schema

	schema.DropColumn("users", "verified")

	assert.Equal(t, Table{
		Op:   Alter,
		Name: "users",
		Definitions: []interface{}{
			Column{Name: "verified", Op: Drop},
		},
	}, schema.Pending[0])
}

func TestSchema_AddIndex(t *testing.T) {
	var schema Schema

	schema.AddIndex("products", []string{"sale"}, SimpleIndex)

	assert.Equal(t, Table{
		Op:   Alter,
		Name: "products",
		Definitions: []interface{}{
			Index{Columns: []string{"sale"}, Type: SimpleIndex, Op: Add},
		},
	}, schema.Pending[0])
}

func TestSchema_RenameIndex(t *testing.T) {
	var schema Schema

	schema.RenameIndex("products", "store_id", "fk_store_id")

	assert.Equal(t, Table{
		Op:   Alter,
		Name: "products",
		Definitions: []interface{}{
			Index{Name: "store_id", NewName: "fk_store_id", Op: Rename},
		},
	}, schema.Pending[0])
}

func TestSchema_DropIndex(t *testing.T) {
	var schema Schema

	schema.DropIndex("products", "sale")

	assert.Equal(t, Table{
		Op:   Alter,
		Name: "products",
		Definitions: []interface{}{
			Index{Name: "sale", Op: Drop},
		},
	}, schema.Pending[0])
}

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigration_tables(t *testing.T) {
	var (
		migration = NewMigration(20200705164100)
	)

	migration.Up(func(m *Migrates) {
		m.CreateTable("products", func(t *Table) {
			t.Integer("id")
			t.String("name")
			t.Text("description")

			t.PrimaryKey("id")
		})

		m.AlterTable("users", func(t *AlterTable) {
			t.Boolean("verified")
			t.RenameColumn("name", "fullname")
		})

		m.RenameTable("trxs", "transactions")

		m.DropTable("logs")
	})

	migration.Down(func(m *Migrates) {
		m.CreateTable("logs", func(t *Table) {
			t.Integer("id")
			t.String("value")

			t.PrimaryKey("id")
		})

		m.RenameTable("transactions", "trxs")

		m.AlterTable("users", func(t *AlterTable) {
			t.DropColumn("verified")
			t.RenameColumn("fullname", "name")
		})

		m.DropTable("products")
	})

	assert.Equal(t, Migration{
		Version: 20200705164100,
		Ups: Migrates{
			Table{
				Op:   Add,
				Name: "products",
				Definitions: []interface{}{
					Column{Name: "id", Type: Integer},
					Column{Name: "name", Type: String},
					Column{Name: "description", Type: Text},
					Index{Columns: []string{"id"}, Type: PrimaryKey},
				},
			},
			Table{
				Op:   Alter,
				Name: "users",
				Definitions: []interface{}{
					Column{Name: "verified", Type: Boolean, Op: Add},
					Column{Name: "name", NewName: "fullname", Op: Rename},
				},
			},
			Table{
				Op:      Rename,
				Name:    "trxs",
				NewName: "transactions",
			},
			Table{
				Op:   Drop,
				Name: "logs",
			},
		},
		Downs: Migrates{
			Table{
				Op:   Add,
				Name: "logs",
				Definitions: []interface{}{
					Column{Name: "id", Type: Integer},
					Column{Name: "value", Type: String},
					Index{Columns: []string{"id"}, Type: PrimaryKey},
				},
			},
			Table{
				Op:      Rename,
				Name:    "transactions",
				NewName: "trxs",
			},
			Table{
				Op:   Alter,
				Name: "users",
				Definitions: []interface{}{
					Column{Name: "verified", Op: Drop},
					Column{Name: "fullname", NewName: "name", Op: Rename},
				},
			},
			Table{
				Op:   Drop,
				Name: "products",
			},
		},
	}, migration)
}

func TestMigration_columns(t *testing.T) {
	var (
		migration = NewMigration(20200805165500)
	)

	migration.Up(func(m *Migrates) {
		m.AddColumn("products", "description", String)
		m.AlterColumn("products", "sale", Boolean)
		m.RenameColumn("users", "name", "fullname")
		m.DropColumn("users", "verified")
		m.AddIndex("products", []string{"sale"}, SimpleIndex)
		m.RenameIndex("products", "store_id", "fk_store_id")
	})

	migration.Down(func(m *Migrates) {
		m.AddColumn("users", "verified", Boolean)
		m.RenameColumn("users", "fullname", "name")
		m.AlterColumn("products", "sale", Integer)
		m.DropColumn("products", "description")
		m.DropIndex("products", "sale")
		m.RenameIndex("products", "fk_store_id", "store_id")
	})

	assert.Equal(t, Migration{
		Version: 20200805165500,
		Ups: Migrates{
			Table{
				Op:   Alter,
				Name: "products",
				Definitions: []interface{}{
					Column{Name: "description", Type: String, Op: Add},
				},
			},
			Table{
				Op:   Alter,
				Name: "products",
				Definitions: []interface{}{
					Column{Name: "sale", Type: Boolean, Op: Alter},
				},
			},
			Table{
				Op:   Alter,
				Name: "users",
				Definitions: []interface{}{
					Column{Name: "name", NewName: "fullname", Op: Rename},
				},
			},
			Table{
				Op:   Alter,
				Name: "users",
				Definitions: []interface{}{
					Column{Name: "verified", Op: Drop},
				},
			},
			Table{
				Op:   Alter,
				Name: "products",
				Definitions: []interface{}{
					Index{Columns: []string{"sale"}, Type: SimpleIndex, Op: Add},
				},
			},
			Table{
				Op:   Alter,
				Name: "products",
				Definitions: []interface{}{
					Index{Name: "store_id", NewName: "fk_store_id", Op: Rename},
				},
			},
		},
		Downs: Migrates{
			Table{
				Op:   Alter,
				Name: "users",
				Definitions: []interface{}{
					Column{Name: "verified", Type: Boolean, Op: Add},
				},
			},
			Table{
				Op:   Alter,
				Name: "users",
				Definitions: []interface{}{
					Column{Name: "fullname", NewName: "name", Op: Rename},
				},
			},
			Table{
				Op:   Alter,
				Name: "products",
				Definitions: []interface{}{
					Column{Name: "sale", Type: Integer, Op: Alter},
				},
			},
			Table{
				Op:   Alter,
				Name: "products",
				Definitions: []interface{}{
					Column{Name: "description", Op: Drop},
				},
			},
			Table{
				Op:   Alter,
				Name: "products",
				Definitions: []interface{}{
					Index{Name: "sale", Op: Drop},
				},
			},
			Table{
				Op:   Alter,
				Name: "products",
				Definitions: []interface{}{
					Index{Name: "fk_store_id", NewName: "store_id", Op: Rename},
				},
			},
		},
	}, migration)
}

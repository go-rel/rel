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
			CreateTable{
				Name: "products",
				Columns: []Column{
					{Name: "id", Type: Integer},
					{Name: "name", Type: String},
					{Name: "description", Type: Text},
				},
			},
			AlterTable{
				Table: Table{
					Name: "users",
					Columns: []Column{
						{Name: "verified", Type: Boolean},
						{Name: "name", NewName: "fullname", Op: RenameColumn},
					},
				},
			},
			RenameTable{
				Name:    "trxs",
				NewName: "transactions",
			},
			DropTable{Name: "logs"},
		},
		Downs: Migrates{
			CreateTable{
				Name: "logs",
				Columns: []Column{
					{Name: "id", Type: Integer},
					{Name: "value", Type: String},
				},
			},
			RenameTable{
				Name:    "transactions",
				NewName: "trxs",
			},
			AlterTable{
				Table: Table{
					Name: "users",
					Columns: []Column{
						{Name: "verified", Op: DropColumn},
						{Name: "fullname", NewName: "name", Op: RenameColumn},
					},
				},
			},
			DropTable{Name: "products"},
		},
	}, migration)
}

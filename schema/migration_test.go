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
			Table{
				Op:   CreateTableOp,
				Name: "products",
				Columns: []Column{
					{Name: "id", Type: Integer},
					{Name: "name", Type: String},
					{Name: "description", Type: Text},
				},
			},
			Table{
				Op:   AlterTableOp,
				Name: "users",
				Columns: []Column{
					{Name: "verified", Type: Boolean, Op: AddColumnOp},
					{Name: "name", NewName: "fullname", Op: RenameColumnOp},
				},
			},
			Table{
				Op:      RenameTableOp,
				Name:    "trxs",
				NewName: "transactions",
			},
			Table{
				Op:   DropTableOp,
				Name: "logs",
			},
		},
		Downs: Migrates{
			Table{
				Op:   CreateTableOp,
				Name: "logs",
				Columns: []Column{
					{Name: "id", Type: Integer},
					{Name: "value", Type: String},
				},
			},
			Table{
				Op:      RenameTableOp,
				Name:    "transactions",
				NewName: "trxs",
			},
			Table{
				Op:   AlterTableOp,
				Name: "users",
				Columns: []Column{
					{Name: "verified", Op: DropColumnOp},
					{Name: "fullname", NewName: "name", Op: RenameColumnOp},
				},
			},
			Table{
				Op:   DropTableOp,
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
	})

	migration.Down(func(m *Migrates) {
		m.AddColumn("users", "verified", Boolean)
		m.RenameColumn("users", "fullname", "name")
		m.AlterColumn("products", "sale", Integer)
		m.DropColumn("products", "description")
	})

	assert.Equal(t, Migration{
		Version: 20200805165500,
		Ups: Migrates{
			Table{
				Op:   AlterTableOp,
				Name: "products",
				Columns: []Column{
					{Name: "description", Type: String, Op: AddColumnOp},
				},
			},
			Table{
				Op:   AlterTableOp,
				Name: "products",
				Columns: []Column{
					{Name: "sale", Type: Boolean, Op: AlterColumnOp},
				},
			},
			Table{
				Op:   AlterTableOp,
				Name: "users",
				Columns: []Column{
					{Name: "name", NewName: "fullname", Op: RenameColumnOp},
				},
			},
			Table{
				Op:   AlterTableOp,
				Name: "users",
				Columns: []Column{
					{Name: "verified", Op: DropColumnOp},
				},
			},
		},
		Downs: Migrates{
			Table{
				Op:   AlterTableOp,
				Name: "users",
				Columns: []Column{
					{Name: "verified", Type: Boolean, Op: AddColumnOp},
				},
			},
			Table{
				Op:   AlterTableOp,
				Name: "users",
				Columns: []Column{
					{Name: "fullname", NewName: "name", Op: RenameColumnOp},
				},
			},
			Table{
				Op:   AlterTableOp,
				Name: "products",
				Columns: []Column{
					{Name: "sale", Type: Integer, Op: AlterColumnOp},
				},
			},
			Table{
				Op:   AlterTableOp,
				Name: "products",
				Columns: []Column{
					{Name: "description", Op: DropColumnOp},
				},
			},
		},
	}, migration)
}

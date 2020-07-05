package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	var (
		migration = NewMigration(20200705164100)
	)

	migration.Up(func(m *Migrate) {
		m.CreateTable("products", func(t *Table) {
			t.Integer("id")
			t.String("name")
			t.Text("description")
		})
	})

	migration.Down(func(m *Migrate) {
		m.DropTable("products")
	})

	assert.Equal(t, Migration{
		Version: 20200705164100,
		Ups: Migrate{
			Table{
				Name: "products",
				Columns: []Column{
					{Name: "id", Type: Integer},
					{Name: "name", Type: String},
					{Name: "description", Type: Text},
				},
			},
		},
	}, migration)
}

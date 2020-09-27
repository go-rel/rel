// 20202806225100_create_todos.go

package migrations

import (
	"github.com/go-rel/rel"
)

// MigrateCreateTodos definition
func MigrateCreateTodos(schema *rel.Schema) {
	schema.CreateTable("todos", func(t *rel.Table) {
		t.ID("id")
		t.DateTime("created_at")
		t.DateTime("updated_at")
		t.String("title")
		t.Bool("completed")
		t.Int("order")
	})

	schema.CreateIndex("todos", "order", []string{"order"})
}

// RollbackCreateTodos definition
func RollbackCreateTodos(schema *rel.Schema) {
	schema.DropTable("todos")
}

package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateColumn(t *testing.T) {
	var (
		options = []ColumnOption{
			Unique(true),
			Required(true),
			Unsigned(true),
			Limit(1000),
			Precision(5),
			Scale(2),
			Default(0),
			Options("options"),
		}
		column = createColumn("add", Decimal, options)
	)

	assert.Equal(t, Column{
		Name:      "add",
		Type:      Decimal,
		Unique:    true,
		Required:  true,
		Unsigned:  true,
		Limit:     1000,
		Precision: 5,
		Scale:     2,
		Default:   0,
		Options:   "options",
	}, column)
}

func TestRenameColumn(t *testing.T) {
	var (
		options = []ColumnOption{
			Required(true),
			Unsigned(true),
			Limit(1000),
			Precision(5),
			Scale(2),
			Default(0),
			Options("options"),
		}
		column = renameColumn("add", "rename", options)
	)

	assert.Equal(t, Column{
		Op:        SchemaRename,
		Name:      "add",
		NewName:   "rename",
		Required:  true,
		Unsigned:  true,
		Limit:     1000,
		Precision: 5,
		Scale:     2,
		Default:   0,
		Options:   "options",
	}, column)
}

func TestDropColumn(t *testing.T) {
	var (
		options = []ColumnOption{
			Required(true),
			Unsigned(true),
			Limit(1000),
			Precision(5),
			Scale(2),
			Default(0),
			Options("options"),
		}
		column = dropColumn("drop", options)
	)

	assert.Equal(t, Column{
		Op:        SchemaDrop,
		Name:      "drop",
		Required:  true,
		Unsigned:  true,
		Limit:     1000,
		Precision: 5,
		Scale:     2,
		Default:   0,
		Options:   "options",
	}, column)
}

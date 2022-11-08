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
			Primary(true),
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
		Primary:   true,
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
		Rename:    "rename",
		Required:  true,
		Unsigned:  true,
		Limit:     1000,
		Precision: 5,
		Scale:     2,
		Default:   0,
		Options:   "options",
	}, column)
}

func TestAlterColumn(t *testing.T) {
	var (
		options = []ColumnOption{
			Required(true),
			Limit(1000),
		}
		columns = alterColumnType("alter", String, options)
	)

	assert.Equal(t, []Column{
		{
			Op:     SchemaAlter,
			Constr: AlterColumnType,
			Type:   String,
			Name:   "alter",
			Limit:  1000,
		},
		{
			Op:       SchemaAlter,
			Constr:   AlterColumnRequired,
			Name:     "alter",
			Required: true,
		},
	}, columns)
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

func TestColumn_InternalTableDefinition(t *testing.T) {
	assert.NotPanics(t, func() { Column{}.internalTableDefinition() })
}

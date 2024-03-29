package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Options("options"),
			Optional(true),
		}
		index = createIndex("table", "add_idx", []string{"add"}, options)
	)

	assert.Equal(t, Index{
		Table:    "table",
		Name:     "add_idx",
		Columns:  []string{"add"},
		Optional: true,
		Options:  "options",
	}, index)
}

func TestCreateUniqueIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Options("options"),
		}
		index = createUniqueIndex("table", "add_idx", []string{"add"}, options)
	)

	assert.Equal(t, Index{
		Table:   "table",
		Name:    "add_idx",
		Unique:  true,
		Columns: []string{"add"},
		Options: "options",
	}, index)
}

func TestCreateFilteredUniqueIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Options("options"),
			Eq("deleted", false),
		}
		index = createUniqueIndex("table", "add_idx", []string{"add"}, options)
	)

	assert.Equal(t, Index{
		Table:   "table",
		Name:    "add_idx",
		Unique:  true,
		Columns: []string{"add"},
		Filter: FilterQuery{
			Type:  FilterEqOp,
			Field: "deleted",
			Value: false,
		},
		Options: "options",
	}, index)
}

func TestDropIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Options("options"),
		}
		index = dropIndex("table", "drop", options)
	)

	assert.Equal(t, Index{
		Op:      SchemaDrop,
		Table:   "table",
		Name:    "drop",
		Options: "options",
	}, index)
}

func TestIndex_Description(t *testing.T) {
	assert.Equal(t, "create index idx_test on tests", Index{Name: "idx_test", Table: "tests"}.description())
}

func TestIndex_InternalMigration(t *testing.T) {
	assert.NotPanics(t, func() { Index{}.internalMigration() })
}

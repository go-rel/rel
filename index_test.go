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

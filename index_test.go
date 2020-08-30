package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
	// TODO: unique option
	var (
		options = []IndexOption{
			Options("options"),
		}
		index = createIndex("table", "add_idx", []string{"add"}, options)
	)

	assert.Equal(t, Index{
		Type:    SimpleIndex,
		Table:   "table",
		Name:    "add_idx",
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

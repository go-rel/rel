package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Name("simple"),
			Options("options"),
		}
		index = createIndex("table", []string{"add"}, SimpleIndex, options)
	)

	assert.Equal(t, Index{
		Type:    SimpleIndex,
		Table:   "table",
		Name:    "simple",
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

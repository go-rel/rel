package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Name("simple"),
			Options("options"),
		}
		index = createIndex([]string{"add"}, SimpleIndex, options)
	)

	assert.Equal(t, Index{
		Type:    SimpleIndex,
		Name:    "simple",
		Columns: []string{"add"},
		Options: "options",
	}, index)
}

func TestRenameIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Options("options"),
		}
		index = renameIndex("add", "rename", options)
	)

	assert.Equal(t, Index{
		Op:      SchemaRename,
		Name:    "add",
		NewName: "rename",
		Options: "options",
	}, index)
}

func TestDropIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Options("options"),
		}
		index = dropIndex("drop", options)
	)

	assert.Equal(t, Index{
		Op:      SchemaDrop,
		Name:    "drop",
		Options: "options",
	}, index)
}

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Name("simple"),
			Comment("comment"),
			Options("options"),
		}
		index = addIndex([]string{"add"}, SimpleIndex, options)
	)

	assert.Equal(t, Index{
		Type:    SimpleIndex,
		Name:    "simple",
		Columns: []string{"add"},
		Comment: "comment",
		Options: "options",
	}, index)
}

func TestAddForeignKey(t *testing.T) {
	var (
		options = []IndexOption{
			OnDelete("cascade"),
			OnUpdate("cascade"),
			Name("fk"),
			Comment("comment"),
			Options("options"),
		}
		index = addForeignKey("table_id", "table", "id", options)
	)

	assert.Equal(t, Index{
		Type:    ForeignKey,
		Name:    "fk",
		Columns: []string{"table_id"},
		Reference: ForeignKeyReference{
			Table:    "table",
			Columns:  []string{"id"},
			OnDelete: "cascade",
			OnUpdate: "cascade",
		},
		Comment: "comment",
		Options: "options",
	}, index)
}

func TestRenameIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Comment("comment"),
			Options("options"),
		}
		index = renameIndex("add", "rename", options)
	)

	assert.Equal(t, Index{
		Op:      Rename,
		Name:    "add",
		NewName: "rename",
		Comment: "comment",
		Options: "options",
	}, index)
}

func TestDropIndex(t *testing.T) {
	var (
		options = []IndexOption{
			Comment("comment"),
			Options("options"),
		}
		index = dropIndex("drop", options)
	)

	assert.Equal(t, Index{
		Op:      Drop,
		Name:    "drop",
		Comment: "comment",
		Options: "options",
	}, index)
}

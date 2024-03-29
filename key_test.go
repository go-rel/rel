package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateForeignKey(t *testing.T) {
	var (
		options = []KeyOption{
			OnDelete("cascade"),
			OnUpdate("cascade"),
			Name("fk"),
			Options("options"),
		}
		index = createForeignKey("table_id", "table", "id", options)
	)

	assert.Equal(t, Key{
		Type:    ForeignKey,
		Name:    "fk",
		Columns: []string{"table_id"},
		Reference: ForeignKeyReference{
			Table:    "table",
			Columns:  []string{"id"},
			OnDelete: "cascade",
			OnUpdate: "cascade",
		},
		Options: "options",
	}, index)
}

func TestCreateUniqueKey(t *testing.T) {
	var (
		options = []KeyOption{
			Name("uq"),
			Options("options"),
		}
		index = createKeys([]string{"code"}, UniqueKey, options)
	)

	assert.Equal(t, Key{
		Type:    UniqueKey,
		Name:    "uq",
		Columns: []string{"code"},
		Options: "options",
	}, index)
}

func TestKey_InternalTableDefinition(t *testing.T) {
	assert.NotPanics(t, func() { Key{}.internalTableDefinition() })
}

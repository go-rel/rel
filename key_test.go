package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateForeignKey(t *testing.T) {
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

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanMigration(t *testing.T) {
	var (
		migrations = scanMigration("testdata/migrations")
	)

	assert.Equal(t, []migration{
		{
			Version: "1",
			Name:    "CreateSamples",
		},
	}, migrations)
}

func TestScanMigration_invalidDir(t *testing.T) {
	assert.Panics(t, func() {
		scanMigration("invalid")
	})
}

func TestScanMigration_invalidFile(t *testing.T) {
	assert.Panics(t, func() {
		scanMigration("./")
	})
}
func TestGetMigrateCommand(t *testing.T) {
	assert.Equal(t, "m.Rollback(ctx)", getMigrateCommand("rollback"))
	assert.Equal(t, "m.Rollback(ctx)", getMigrateCommand("down"))
	assert.Equal(t, "m.Migrate(ctx)", getMigrateCommand("migrate"))
	assert.Equal(t, "m.Migrate(ctx)", getMigrateCommand("up"))
}

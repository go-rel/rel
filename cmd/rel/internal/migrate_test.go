package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanMigration(t *testing.T) {
	tests := []struct {
		dir        string
		migrations []migration
		err        error
	}{
		{
			dir: "testdata/migrations",
			migrations: []migration{
				{
					Version: "1",
					Name:    "CreateSamples",
				},
			},
		},
		{
			dir: "db",
			err: errors.New("rel: open db: no such file or directory"),
		},
		{
			dir: "../",
			err: errors.New("rel: invalid migration file: main.go"),
		},
	}

	for _, test := range tests {
		t.Run(test.dir, func(t *testing.T) {
			var (
				migrations, err = scanMigration(test.dir)
			)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.migrations, migrations)
		})
	}

}

func TestGetMigrateCommand(t *testing.T) {
	assert.Equal(t, "m.Rollback(ctx)", getMigrateCommand("rollback"))
	assert.Equal(t, "m.Rollback(ctx)", getMigrateCommand("down"))
	assert.Equal(t, "m.Migrate(ctx)", getMigrateCommand("migrate"))
	assert.Equal(t, "m.Migrate(ctx)", getMigrateCommand("up"))
}

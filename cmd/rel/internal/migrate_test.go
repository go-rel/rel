package internal

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecMigrate(t *testing.T) {
	t.Run("invalid migration dir", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			args = []string{
				"rel",
				"migrate",
				"-dir=db",
			}
		)

		assert.Equal(t, errors.New("rel: open db: no such file or directory"), ExecMigrate(ctx, args))
	})

	t.Run("success", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			args = []string{
				"rel",
				"migrate",
				"-dir=testdata/migrations",
				"-module=github.com/Fs02/rel/cmd/rel/internal",
				"-adapter=github.com/Fs02/rel/adapter/sqlite3",
				"-driver=github.com/mattn/go-sqlite3",
				"-dsn=:memory:",
				"-verbose=false",
			}
			dir  = "testdata"
			buff = &bytes.Buffer{}
		)

		stderr = buff
		defer func() { stderr = os.Stderr }()

		tempDir = func() string { return dir }
		defer func() { tempDir = os.TempDir }()

		err := ExecMigrate(ctx, args)
		assert.Empty(t, buff.String())
		assert.Nil(t, err)
	})
}

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

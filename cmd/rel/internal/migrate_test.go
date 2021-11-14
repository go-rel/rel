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
	t.Run("missing required parameters", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			args = []string{
				"rel",
				"migrate",
			}
		)

		assert.Equal(t, errors.New("rel: missing required parameters:\n\tadapter: \n\tdriver: \n\tdsn: "), ExecMigrate(ctx, args))
	})

	t.Run("invalid migration dir", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			args = []string{
				"rel",
				"migrate",
				"-adapter=github.com/go-rel/sqlite3",
				"-driver=github.com/mattn/go-sqlite3",
				"-dsn=:memory:",
				"-dir=db",
			}
		)

		assert.Equal(t, errors.New("rel: error accessing read migration directory: db"), ExecMigrate(ctx, args))
	})

	t.Run("success", func(t *testing.T) {
		var (
			ctx  = context.TODO()
			args = []string{
				"rel",
				"migrate",
				"-dir=testdata/migrations",
				"-module=github.com/go-rel/rel/cmd/rel/internal",
				"-adapter=github.com/go-rel/sqlite3",
				"-driver=github.com/mattn/go-sqlite3",
				"-dsn=:memory:",
				"-verbose=false",
			}
			dir  = "testdata"
			buff = &bytes.Buffer{}
		)

		tempdir = dir
		stderr = buff
		defer func() { stderr = os.Stderr }()

		err := ExecMigrate(ctx, args)
		assert.Contains(t, buff.String(), "Running: migrate 1 create table todos")
		assert.Contains(t, buff.String(), "Done: migrate 1 create table todos")
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
			err: errors.New("rel: error accessing read migration directory: db"),
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

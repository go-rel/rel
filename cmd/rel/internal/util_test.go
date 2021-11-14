package internal

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseInfo(t *testing.T) {
	t.Run("database", func(t *testing.T) {
		os.Setenv("DATABASE_ADAPTER", "github.com/go-rel/mysql")
		os.Setenv("DATABASE_DRIVER", "github.com/go-sql-driver/mysql")
		os.Setenv("DATABASE_URL", "user:password@(localhost:3306)/db?charset=utf8&parseTime=True&loc=Local")

		defer os.Setenv("DATABASE_ADAPTER", "")
		defer os.Setenv("DATABASE_DRIVER", "")
		defer os.Setenv("DATABASE_URL", "")

		adapter, driver, url := getDatabaseInfo()
		assert.Equal(t, "github.com/go-rel/mysql", adapter)
		assert.Equal(t, "github.com/go-sql-driver/mysql", driver)
		assert.Equal(t, "user:password@(localhost:3306)/db?charset=utf8&parseTime=True&loc=Local", url)
	})

	t.Run("mysql", func(t *testing.T) {
		os.Setenv("MYSQL_HOST", "localhost")
		os.Setenv("MYSQL_PORT", "3306")
		os.Setenv("MYSQL_DATABASE", "db")
		os.Setenv("MYSQL_USERNAME", "user")
		os.Setenv("MYSQL_PASSWORD", "password")

		defer os.Setenv("MYSQL_HOST", "")
		defer os.Setenv("MYSQL_PORT", "")
		defer os.Setenv("MYSQL_DATABASE", "")
		defer os.Setenv("MYSQL_USERNAME", "")
		defer os.Setenv("MYSQL_PASSWORD", "")

		adapter, driver, url := getDatabaseInfo()
		assert.Equal(t, "github.com/go-rel/mysql", adapter)
		assert.Equal(t, "github.com/go-sql-driver/mysql", driver)
		assert.Equal(t, "user:password@(localhost:3306)/db?charset=utf8&parseTime=True&loc=Local", url)
	})

	t.Run("postgresql", func(t *testing.T) {
		os.Setenv("POSTGRES_HOST", "localhost")
		os.Setenv("POSTGRES_PORT", "5432")
		os.Setenv("POSTGRES_DATABASE", "db")
		os.Setenv("POSTGRES_USERNAME", "user")
		os.Setenv("POSTGRES_PASSWORD", "password")

		defer os.Setenv("POSTGRES_HOST", "")
		defer os.Setenv("POSTGRES_PORT", "")
		defer os.Setenv("POSTGRES_DATABASE", "")
		defer os.Setenv("POSTGRES_USERNAME", "")
		defer os.Setenv("POSTGRES_PASSWORD", "")

		adapter, driver, url := getDatabaseInfo()
		assert.Equal(t, "github.com/go-rel/postgres", adapter)
		assert.Equal(t, "github.com/lib/pq", driver)
		assert.Equal(t, "postgres://user:password@localhost:5432/db?sslmode=disable", url)
	})

	t.Run("postgresql alternative", func(t *testing.T) {
		os.Setenv("POSTGRESQL_HOST", "localhost")
		os.Setenv("POSTGRESQL_PORT", "5432")
		os.Setenv("POSTGRESQL_DATABASE", "db")
		os.Setenv("POSTGRESQL_USERNAME", "user")
		os.Setenv("POSTGRESQL_PASSWORD", "password")

		defer os.Setenv("POSTGRESQL_HOST", "")
		defer os.Setenv("POSTGRESQL_PORT", "")
		defer os.Setenv("POSTGRESQL_DATABASE", "")
		defer os.Setenv("POSTGRESQL_USERNAME", "")
		defer os.Setenv("POSTGRESQL_PASSWORD", "")

		adapter, driver, url := getDatabaseInfo()
		assert.Equal(t, "github.com/go-rel/postgres", adapter)
		assert.Equal(t, "github.com/lib/pq", driver)
		assert.Equal(t, "postgres://user:password@localhost:5432/db?sslmode=disable", url)
	})

	t.Run("sqlite3", func(t *testing.T) {
		os.Setenv("SQLITE3_DATABASE", "test.db")
		defer os.Setenv("SQLITE3_DATABASE", "")

		adapter, driver, url := getDatabaseInfo()
		assert.Equal(t, "github.com/go-rel/sqlite3", adapter)
		assert.Equal(t, "github.com/mattn/go-sqlite3", driver)
		assert.Equal(t, "test.db", url)
	})

}

func TestGetModule(t *testing.T) {
	t.Run("gomod", func(t *testing.T) {
		var (
			file, _ = ioutil.TempFile(os.TempDir(), "go.mod")
		)

		defer os.Remove(file.Name())
		file.WriteString("module github.com/Fs02/go-todo-backend")
		file.Close()

		gomod = file.Name()
		assert.Equal(t, "github.com/Fs02/go-todo-backend", getModule())
	})

	t.Run("gomod invalid", func(t *testing.T) {
		var (
			file, _ = ioutil.TempFile(os.TempDir(), "go.mod")
		)

		defer os.Remove(file.Name())
		file.WriteString("pkg github.com/Fs02/go-todo-backend")
		file.Close()

		gomod = file.Name()
		assert.NotEqual(t, "github.com/Fs02/go-todo-backend", getModule())
	})

	t.Run("gopath", func(t *testing.T) {
		assert.NotEmpty(t, getModule())
	})
}

func TestInternal(t *testing.T) {
	assert.Panics(t, func() { check(errors.New("err")) })
	assert.NotPanics(t, func() { check(nil) })
}

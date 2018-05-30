package sql

import (
	db "database/sql"
	"testing"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func open() (*Adapter, error) {
	var err error
	adapter := New(
		func(err error) error { return err },
		func(Adapter) int { return 1 },
		Placeholder("?"),
		Ordinal(false),
		InsertDefaultValues(true),
	)

	// simplified tests using sqlite backend.
	adapter.DB, err = db.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	paranoid.Panic(err, "failed to open database connection")

	_, _, execerr := adapter.Exec(`CREATE TABLE test (
		id INTEGER PRIMARY KEY,
		name STRING
	);`, nil)
	paranoid.Panic(execerr, "failed creating test table")

	return adapter, err
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New(nil, nil))
}

func TestAdapter_Count(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	_, err = grimoire.New(adapter).From("test").Count()
	assert.Nil(t, err)
}

func TestAdapter_All(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	result := []struct{}{}
	assert.Nil(t, grimoire.New(adapter).From("test").All(&result))
}

func TestAdapter_Insert(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	result := struct {
		Name string
	}{}
	ch := changeset.Convert(result)
	assert.Nil(t, grimoire.New(adapter).From("test").Insert(&result, ch))
	assert.Nil(t, grimoire.New(adapter).From("test").Insert(nil, ch, ch))
}

func TestAdapter_Update(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	result := struct {
		Name string
	}{}
	ch := changeset.Convert(result)

	assert.Nil(t, grimoire.New(adapter).From("test").Update(nil, ch))
}

func TestAdapter_Delete(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	assert.Nil(t, grimoire.New(adapter).From("test").Delete())
}

func TestAdapter_Transaction_commit(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	result := struct {
		Name string
	}{}
	ch := changeset.Convert(result)

	err = grimoire.New(adapter).Transaction(func(repo grimoire.Repo) error {
		repo.From("test").MustInsert(&result, ch)
		return nil
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_rollback(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	err = grimoire.New(adapter).Transaction(func(repo grimoire.Repo) error {
		return errors.New("", "", errors.UniqueConstraint)
	})

	assert.NotNil(t, err)
}

func TestAdapter_InsertAll_error(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	fields := []string{"notexist"}
	allchanges := []map[string]interface{}{
		{"notexist": "12"},
		{"notexist": "13"},
	}

	_, err = adapter.InsertAll(grimoire.Repo{}.From("test"), fields, allchanges)

	assert.NotNil(t, err)
}

func TestAdapter_Transaction_commitError(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit())
}

func TestAdapter_Transaction_rollbackError(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback())
}

func TestAdapter_Query_error(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	out := struct{}{}

	_, err = adapter.Query(&out, "error", nil)
	assert.NotNil(t, err)
}

func TestAdapter_Exec_error(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	_, _, err = adapter.Exec("error", nil)
	assert.NotNil(t, err)
}

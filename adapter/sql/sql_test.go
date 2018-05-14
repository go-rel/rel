package sql

import (
	db "database/sql"
	"testing"

	paranoid "github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func open() (*Adapter, error) {
	var err error
	adapter := &Adapter{
		Placeholder:   "?",
		Ordinal:       false,
		IncrementFunc: func(Adapter) int { return 1 },
		ErrorFunc:     func(err error) error { return err },
	}

	// simplified tests using sqlite backend.
	adapter.DB, err = db.Open("sqlite3", "file::memory:?mode=memory&cache=shared")

	_, _, execerr := adapter.Exec(`CREATE TABLE test (
		id INTEGER PRIMARY KEY,
		name STRING
	);`, nil)
	paranoid.Panic(execerr)

	return adapter, err
}

func TestAdapterNew(t *testing.T) {
	assert.NotNil(t, New("?", false, nil, nil))
}

func TestAdapterCount(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	_, err = grimoire.New(adapter).From("test").Count()
	assert.Nil(t, err)
}

func TestAdapterAll(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	result := []struct{}{}
	assert.Nil(t, grimoire.New(adapter).From("test").All(&result))
}

func TestAdapterInsert(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	result := struct {
		Name string
	}{}
	ch := changeset.Change(result)
	assert.Nil(t, grimoire.New(adapter).From("test").Insert(&result, ch))
	assert.Nil(t, grimoire.New(adapter).From("test").Insert(nil, ch, ch))
}

func TestAdapterUpdate(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	result := struct {
		Name string
	}{}
	ch := changeset.Change(result)

	assert.Nil(t, grimoire.New(adapter).From("test").Update(nil, ch))
}

func TestAdapterDelete(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	assert.Nil(t, grimoire.New(adapter).From("test").Delete())
}

func TestAdapterTransactionCommit(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	result := struct {
		Name string
	}{}
	ch := changeset.Change(result)

	err = grimoire.New(adapter).Transaction(func(repo grimoire.Repo) error {
		repo.From("test").MustInsert(&result, ch)
		return nil
	})

	assert.Nil(t, err)
}

func TestAdapterTransactionRollback(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	err = grimoire.New(adapter).Transaction(func(repo grimoire.Repo) error {
		return errors.UniqueConstraintError("", "")
	})

	assert.NotNil(t, err)
}

func TestAdapterInsertAllError(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	fields := []string{"notexist"}
	allchanges := []map[string]interface{}{
		{"notexist": "12"},
		{"notexist": "13"},
	}

	_, err = adapter.InsertAll(grimoire.Repo{}.From("test"), fields, allchanges)

	assert.NotNil(t, err)
}

func TestAdapterTransactionCommitError(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit())
}

func TestAdapterTransactionRollbackError(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback())
}

func TestAdapterQueryError(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	out := struct{}{}

	_, err = adapter.Query(&out, "error", nil)
	assert.NotNil(t, err)
}

func TestAdapterExecError(t *testing.T) {
	adapter, err := open()
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	_, _, err = adapter.Exec("error", nil)
	assert.NotNil(t, err)
}

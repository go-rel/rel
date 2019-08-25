package sql

import (
	db "database/sql"
	"errors"
	"testing"

	"github.com/Fs02/grimoire"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func open(t *testing.T) *Adapter {
	var (
		err    error
		config = &Config{
			Placeholder:         "?",
			EscapeChar:          "`",
			InsertDefaultValues: true,
			ErrorFunc:           func(err error) error { return err },
			IncrementFunc:       func(Adapter) int { return 1 },
		}
		adapter = New(config)
	)

	// simplified tests using sqlite backend.
	adapter.DB, err = db.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	assert.Nil(t, err)

	_, _, err = adapter.Exec(`CREATE TABLE IF NOT EXISTS names (
		id INTEGER PRIMARY KEY,
		name STRING
	);`, nil)
	assert.Nil(t, err)

	return adapter
}

type Name struct {
	ID   int
	Name string
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New(nil))
}

// func TestAdapter_Count(t *testing.T) {
// 	adapter, err := open()
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	_, err = repo.From("test").Count()
// 	assert.Nil(t, err)
// }

// func TestAdapter_All(t *testing.T) {
// 	adapter, err := open()
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	result := []Name{}
// 	assert.Nil(t, repo.All(&result)) //From("test").All(&result))
// }

func TestAdapter_Insert(t *testing.T) {
	var (
		adapter = open(t)
		repo    = grimoire.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)
	defer adapter.Close()

	assert.Nil(t, repo.Insert(&name))
	assert.NotEqual(t, 0, name.ID)
}

func TestAdapter_Update(t *testing.T) {
	var (
		adapter = open(t)
		repo    = grimoire.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)

	defer adapter.Close()

	assert.Nil(t, repo.Insert(&name))
	assert.NotEqual(t, 0, name.ID)

	name.Name = "Zoro"

	assert.Nil(t, repo.Update(&name))
	assert.NotEqual(t, 0, name.ID)
	assert.Equal(t, "Zoro", name.Name)
}

func TestAdapter_Delete(t *testing.T) {
	var (
		adapter = open(t)
		repo    = grimoire.New(adapter)
		name    = Name{}
	)

	defer adapter.Close()

	assert.Nil(t, repo.Delete(&name))
}

func TestAdapter_Transaction_commit(t *testing.T) {
	var (
		adapter = open(t)
		repo    = grimoire.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)

	err := repo.Transaction(func(repo grimoire.Repo) error {
		repo.MustInsert(&name)
		return nil
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_rollback(t *testing.T) {
	var (
		adapter = open(t)
		repo    = grimoire.New(adapter)
	)

	err := repo.Transaction(func(repo grimoire.Repo) error {
		return errors.New("error")
	})

	assert.NotNil(t, err)
}

func TestAdapter_Transaction_nestedCommit(t *testing.T) {
	var (
		adapter = open(t)
		repo    = grimoire.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)

	defer adapter.Close()

	err := repo.Transaction(func(repo grimoire.Repo) error {
		return repo.Transaction(func(repo grimoire.Repo) error {
			repo.MustInsert(&name)
			return nil
		})
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_nestedRollback(t *testing.T) {
	var (
		adapter = open(t)
		repo    = grimoire.New(adapter)
	)

	defer adapter.Close()

	err := repo.Transaction(func(repo grimoire.Repo) error {
		return repo.Transaction(func(repo grimoire.Repo) error {
			return errors.New("error")
		})
	})

	assert.NotNil(t, err)
}

// func TestAdapter_InsertAll_error(t *testing.T) {
// 	adapter, err := open()
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	fields := []string{"notexist"}
// 	allchanges := []map[string]interface{}{
// 		{"notexist": "12"},
// 		{"notexist": "13"},
// 	}

// 	_, err = adapter.InsertAll(query.Query{}, fields, allchanges)

// 	assert.NotNil(t, err)
// }

func TestAdapter_Transaction_commitError(t *testing.T) {
	var (
		adapter = open(t)
	)

	defer adapter.Close()

	assert.NotNil(t, adapter.Commit())
}

func TestAdapter_Transaction_rollbackError(t *testing.T) {
	var (
		adapter = open(t)
	)

	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback())
}

// func TestAdapter_Query_error(t *testing.T) {
// 	adapter, err := open()
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	out := struct{}{}

// 	_, err = adapter.Query(&out, "error", nil)
// 	assert.NotNil(t, err)
// }

func TestAdapter_Exec_error(t *testing.T) {
	var (
		adapter = open(t)
	)

	defer adapter.Close()

	_, _, err := adapter.Exec("error", nil)
	assert.NotNil(t, err)
}

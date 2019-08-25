package sql

import (
	db "database/sql"
	"errors"
	"testing"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func open() (*Adapter, error) {
	var err error
	adapter := New(&Config{
		Placeholder:         "?",
		EscapeChar:          "`",
		InsertDefaultValues: true,
		ErrorFunc:           func(err error) error { return err },
		IncrementFunc:       func(Adapter) int { return 1 },
	})

	// simplified tests using sqlite backend.
	adapter.DB, err = db.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	paranoid.Panic(err, "failed to open database connection")

	_, _, execerr := adapter.Exec(`CREATE TABLE names (
		id INTEGER PRIMARY KEY,
		name STRING
	);`, nil)
	paranoid.Panic(execerr, "failed creating names table")

	return adapter, err
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

// 	_, err = grimoire.New(adapter).From("test").Count()
// 	assert.Nil(t, err)
// }

// func TestAdapter_All(t *testing.T) {
// 	adapter, err := open()
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	result := []Name{}
// 	assert.Nil(t, grimoire.New(adapter).All(&result)) //From("test").All(&result))
// }

func TestAdapter_Insert(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	var (
		name = Name{
			Name: "Luffy",
		}
	)

	assert.Nil(t, grimoire.New(adapter).Insert(&name))
	assert.NotEqual(t, 0, name.ID)
}

func TestAdapter_Update(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	var (
		name = Name{
			Name: "Luffy",
		}
	)

	assert.Nil(t, grimoire.New(adapter).Insert(&name))
	assert.NotEqual(t, 0, name.ID)

	name.Name = "Zoro"

	assert.Nil(t, grimoire.New(adapter).Update(&name))
	assert.NotEqual(t, 0, name.ID)
	assert.Equal(t, "Zoro", name.Name)
}

func TestAdapter_Delete(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	assert.Nil(t, grimoire.New(adapter).Delete(&Name{}))
}

func TestAdapter_Transaction_commit(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	var (
		name = Name{
			Name: "Luffy",
		}
	)

	err = grimoire.New(adapter).Transaction(func(repo grimoire.Repo) error {
		repo.MustInsert(&name)
		return nil
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_rollback(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	err = grimoire.New(adapter).Transaction(func(repo grimoire.Repo) error {
		return errors.New("error")
	})

	assert.NotNil(t, err)
}

func TestAdapter_Transaction_nestedCommit(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	var (
		name = Name{
			Name: "Luffy",
		}
	)

	err = grimoire.New(adapter).Transaction(func(repo grimoire.Repo) error {
		return repo.Transaction(func(repo grimoire.Repo) error {
			repo.MustInsert(&name)
			return nil
		})
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_nestedRollback(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	err = grimoire.New(adapter).Transaction(func(repo grimoire.Repo) error {
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

// func TestAdapter_Query_error(t *testing.T) {
// 	adapter, err := open()
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	out := struct{}{}

// 	_, err = adapter.Query(&out, "error", nil)
// 	assert.NotNil(t, err)
// }

func TestAdapter_Exec_error(t *testing.T) {
	adapter, err := open()
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	_, _, err = adapter.Exec("error", nil)
	assert.NotNil(t, err)
}

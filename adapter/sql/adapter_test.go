package sql

import (
	db "database/sql"
	"errors"
	"testing"

	"github.com/Fs02/rel"
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
			IncrementFunc:       func(Adapter) int { return -1 },
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

func TestAdapter_Aggregate(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	count, err := repo.Aggregate(rel.From("names"), "count", "id")
	assert.Equal(t, 0, count)
	assert.Nil(t, err)
}

func TestAdapter_Aggregate_transaction(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	repo.Transaction(func(repo rel.Repository) error {
		count, err := repo.Aggregate(rel.From("names"), "count", "id")
		assert.Equal(t, 0, count)
		assert.Nil(t, err)

		return nil
	})
}

func TestAdapter_Insert(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)
	defer adapter.Close()

	assert.Nil(t, repo.Insert(&name))
	assert.NotEqual(t, 0, name.ID)
}

func TestAdapter_InsertAll(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
		names   = []Name{
			{Name: "Luffy"},
			{Name: "Zoro"},
		}
	)
	defer adapter.Close()

	assert.Nil(t, repo.InsertAll(&names))
	assert.Len(t, names, 2)
	assert.NotEqual(t, 0, names[0].ID)
	assert.NotEqual(t, 0, names[1].ID)
	assert.Equal(t, "Luffy", names[0].Name)
	assert.Equal(t, "Zoro", names[1].Name)
}

func TestAdapter_Update(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
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
		repo    = rel.New(adapter)
		name    = Name{}
	)

	defer adapter.Close()

	assert.Nil(t, repo.Delete(&name))
}

func TestAdapter_Transaction_commit(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)

	err := repo.Transaction(func(repo rel.Repository) error {
		repo.MustInsert(&name)
		return nil
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_rollback(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	err := repo.Transaction(func(repo rel.Repository) error {
		return errors.New("error")
	})

	assert.NotNil(t, err)
}

func TestAdapter_Transaction_nestedCommit(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)

	defer adapter.Close()

	err := repo.Transaction(func(repo rel.Repository) error {
		return repo.Transaction(func(repo rel.Repository) error {
			repo.MustInsert(&name)
			return nil
		})
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_nestedRollback(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	err := repo.Transaction(func(repo rel.Repository) error {
		return repo.Transaction(func(repo rel.Repository) error {
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

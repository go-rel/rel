package sql

import (
	"context"
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

	_, _, err = adapter.Exec(context.TODO(), `CREATE TABLE IF NOT EXISTS names (
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

func TestAdapter_Ping(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	assert.Nil(t, repo.Ping(context.TODO()))
}

func TestAdapter_Aggregate(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	count, err := repo.Aggregate(context.TODO(), rel.From("names"), "count", "id")
	assert.Equal(t, 0, count)
	assert.Nil(t, err)
}

func TestAdapter_Aggregate_transaction(t *testing.T) {
	var (
		ctx     = context.TODO()
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	repo.Transaction(ctx, func(ctx context.Context) error {
		count, err := repo.Aggregate(ctx, rel.From("names"), "count", "id")
		assert.Equal(t, 0, count)
		assert.Nil(t, err)

		return nil
	})
}

func TestAdapter_FindAll(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	assert.Nil(t, repo.FindAll(context.TODO(), &[]struct{}{}, rel.From("names")))
}

func TestAdapter_FindAll_transaction(t *testing.T) {
	var (
		ctx     = context.TODO()
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	assert.Nil(t, repo.Transaction(ctx, func(ctx context.Context) error {
		return repo.FindAll(ctx, &[]struct{}{}, rel.From("names"))
	}))
}

func TestAdapter_Query_error(t *testing.T) {
	var (
		adapter = open(t)
	)
	defer adapter.Close()

	_, err := adapter.Query(context.TODO(), rel.Query{})
	assert.NotNil(t, err)
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

	assert.Nil(t, repo.Insert(context.TODO(), &name))
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

	assert.Nil(t, repo.InsertAll(context.TODO(), &names))
	assert.Len(t, names, 2)
	assert.NotEqual(t, 0, names[0].ID)
	assert.NotEqual(t, 0, names[1].ID)
	assert.Equal(t, "Luffy", names[0].Name)
	assert.Equal(t, "Zoro", names[1].Name)
}

func TestAdapter_InsertAll_customPrimary(t *testing.T) {
	var (
		adapter = open(t)
		repo    = rel.New(adapter)
		names   = []Name{
			{ID: 20, Name: "Luffy"},
			{ID: 21, Name: "Zoro"},
		}
	)
	defer adapter.Close()

	assert.Nil(t, repo.InsertAll(context.TODO(), &names))
	assert.Len(t, names, 2)
	assert.Equal(t, 20, names[0].ID)
	assert.Equal(t, 21, names[1].ID)
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

	assert.Nil(t, repo.Insert(context.TODO(), &name))
	assert.NotEqual(t, 0, name.ID)

	name.Name = "Zoro"

	assert.Nil(t, repo.Update(context.TODO(), &name))
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

	repo.MustInsert(context.TODO(), &name)

	assert.Nil(t, repo.Delete(context.TODO(), &name))
}

func TestAdapter_Transaction_commit(t *testing.T) {
	var (
		ctx     = context.TODO()
		adapter = open(t)
		repo    = rel.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)

	err := repo.Transaction(ctx, func(ctx context.Context) error {
		repo.MustInsert(ctx, &name)
		return nil
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_rollback(t *testing.T) {
	var (
		ctx     = context.TODO()
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	err := repo.Transaction(ctx, func(ctx context.Context) error {
		return errors.New("error")
	})

	assert.NotNil(t, err)
}

func TestAdapter_Transaction_nestedCommit(t *testing.T) {
	var (
		ctx     = context.TODO()
		adapter = open(t)
		repo    = rel.New(adapter)
		name    = Name{
			Name: "Luffy",
		}
	)

	defer adapter.Close()

	err := repo.Transaction(ctx, func(ctx context.Context) error {
		return repo.Transaction(ctx, func(ctx context.Context) error {
			repo.MustInsert(ctx, &name)
			return nil
		})
	})

	assert.Nil(t, err)
}

func TestAdapter_Transaction_nestedRollback(t *testing.T) {
	var (
		ctx     = context.TODO()
		adapter = open(t)
		repo    = rel.New(adapter)
	)

	defer adapter.Close()

	err := repo.Transaction(ctx, func(ctx context.Context) error {
		return repo.Transaction(ctx, func(ctx context.Context) error {
			return errors.New("error")
		})
	})

	assert.NotNil(t, err)
}

func TestAdapter_InsertAll_error(t *testing.T) {
	var (
		adapter = open(t)
	)
	defer adapter.Close()

	fields := []string{"notexist"}
	mutations := []map[string]rel.Mutate{
		{"notexist": rel.Set("notexist", "13")},
		{"notexist": rel.Set("notexist", "12")},
	}

	ids, err := adapter.InsertAll(context.TODO(), rel.Query{}, "id", fields, mutations)
	assert.NotNil(t, err)
	assert.Nil(t, ids)
}

func TestAdapter_Transaction_commitError(t *testing.T) {
	var (
		adapter = open(t)
	)

	defer adapter.Close()

	assert.NotNil(t, adapter.Commit(context.TODO()))
}

func TestAdapter_Transaction_rollbackError(t *testing.T) {
	var (
		adapter = open(t)
	)

	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback(context.TODO()))
}

func TestAdapter_Exec_error(t *testing.T) {
	var (
		adapter = open(t)
	)

	defer adapter.Close()

	_, _, err := adapter.Exec(context.TODO(), "error", nil)
	assert.NotNil(t, err)
}

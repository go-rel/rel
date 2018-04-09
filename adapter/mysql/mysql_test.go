package mysql

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/specs"
	"github.com/Fs02/grimoire/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestSpecs(t *testing.T) {
	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()
	repo := grimoire.New(adapter)

	// Query Specs
	specs.Query(t, repo)
	specs.QueryJoin(t, repo)
	specs.QueryNotFound(t, repo)

	// Insert Specs
	specs.Insert(t, repo)
	specs.InsertAll(t, repo)
	specs.InsertSet(t, repo)

	// Update Specs
	specs.Update(t, repo)
	specs.UpdateWhere(t, repo)
	specs.UpdateSet(t, repo)

	// Put Specs
	specs.SaveInsert(t, repo)
	specs.SaveInsertAll(t, repo)
	specs.PutUpdate(t, repo)

	//Delete specs
	specs.Delete(t, repo)
}

func TestAdapterInsertAllError(t *testing.T) {
	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	fields := []string{"notexist"}
	allchanges := []map[string]interface{}{
		{"notexist": "12"},
		{"notexist": "13"},
	}

	_, err = adapter.InsertAll(grimoire.Repo{}.From("users"), fields, allchanges)

	assert.NotNil(t, err)
}

func TestAdapterTransactionCommit(t *testing.T) {
	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	txAdapter, err := adapter.Begin()
	assert.Nil(t, err)
	assert.Nil(t, txAdapter.Commit())
}

func TestAdapterTransactionRollback(t *testing.T) {
	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	txAdapter, err := adapter.Begin()
	assert.Nil(t, err)
	assert.Nil(t, txAdapter.Rollback())
}

func TestAdapterTransactionCommitError(t *testing.T) {
	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit())
}

func TestAdapterTransactionRollbackError(t *testing.T) {
	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback())
}

// func TestAdapterQuery(t *testing.T) {
// 	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer adapter.Close()

// 	out := struct{}{}

// 	// normal query
// 	count, err := adapter.Query(&out, "SELECT 10;", []interface{}{})
// 	assert.Nil(t, err)
// 	assert.Equal(t, int64(1), count)

// 	// within transaction
// 	txAdapter, err := adapter.Begin()
// 	assert.Nil(t, err)

// 	count, err = txAdapter.Query(&out, "SELECT 10;", []interface{}{})
// 	assert.Nil(t, err)

// 	assert.Equal(t, int64(1), count)
// 	assert.Nil(t, txAdapter.Commit())
// }

// func TestAdapterQueryError(t *testing.T) {
// 	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer adapter.Close()

// 	out := struct{}{}

// 	_, err = adapter.Query(&out, ";;", []interface{}{})
// 	assert.NotNil(t, err)
// }

// func TestAdapterExec(t *testing.T) {
// 	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer adapter.Close()

// 	// normal exec
// 	stmt := "INSERT INTO users (name, created_at, updated_at) VALUES (?,?,?)"
// 	args := []interface{}{"find", time.Now().Round(time.Second), time.Now().Round(time.Second)}
// 	id, count, err := adapter.Exec(stmt, args)

// 	assert.Nil(t, err)
// 	assert.True(t, id > 0)
// 	assert.Equal(t, int64(1), count)

// 	// within transaction
// 	// within transaction
// 	txAdapter, err := adapter.Begin()
// 	assert.Nil(t, err)

// 	id, count, err = txAdapter.Exec(stmt, args)
// 	assert.Nil(t, err)
// 	assert.True(t, id > 0)
// 	assert.Equal(t, int64(1), count)

// 	assert.Nil(t, txAdapter.Commit())
// }

// func TestAdapterExecError(t *testing.T) {
// 	adapter, err := Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer adapter.Close()

// 	_, _, err = adapter.Exec(";;", []interface{}{})
// 	assert.NotNil(t, err)
// }

func TestAdapterError(t *testing.T) {
	adapter := new(Adapter)

	// error nil
	assert.Nil(t, adapter.Error(nil))

	// 1062 error
	rawerr := &mysql.MySQLError{Message: "duplicate", Number: 1062}
	duperr := errors.DuplicateError(rawerr.Message, "")
	assert.Equal(t, duperr, adapter.Error(rawerr))

	// other errors
	err := errors.UnexpectedError("error")
	assert.Equal(t, err, adapter.Error(err))
}

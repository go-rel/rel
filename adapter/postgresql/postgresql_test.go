package postgresql

import (
	"os"
	"testing"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/specs"
	"github.com/Fs02/grimoire/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	adapter, err := Open(dsn())
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS addresses;`, []interface{}{})
	paranoid.Panic(err)
	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS users;`, []interface{}{})
	paranoid.Panic(err)

	_, _, err = adapter.Exec(`CREATE TABLE users (
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(30) NOT NULL,
		gender VARCHAR(10) NOT NULL,
		age INT NOT NULL,
		note varchar(50),
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);`, []interface{}{})
	paranoid.Panic(err)

	_, _, err = adapter.Exec(`CREATE TABLE addresses (
		id SERIAL NOT NULL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		address VARCHAR(60) NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);`, []interface{}{})
	paranoid.Panic(err)
}

func dsn() string {
	if os.Getenv("POSTGRESQL_DATABASE") != "" {
		return os.Getenv("POSTGRESQL_DATABASE")
	}

	return "postgres://postgres@localhost/grimoire_test?sslmode=disable"
}

func TestSpecs(t *testing.T) {
	adapter, err := Open(dsn())
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
	// specs.Insert(t, repo)
	// specs.InsertAll(t, repo)
	// specs.InsertSet(t, repo)

	// Update Specs
	// specs.Update(t, repo)
	// specs.UpdateWhere(t, repo)
	// specs.UpdateSet(t, repo)

	// Put Specs
	// specs.SaveInsert(t, repo)
	// specs.SaveInsertAll(t, repo)
	// specs.SaveUpdate(t, repo)

	// Delete specs
	// specs.Delete(t, repo)

	// Transaction specs
	// specs.Transaction(t, repo)
}

func TestAdapterInsertAllError(t *testing.T) {
	adapter, err := Open(dsn())
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

func TestAdapterTransactionCommitError(t *testing.T) {
	adapter, err := Open(dsn())
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit())
}

func TestAdapterTransactionRollbackError(t *testing.T) {
	adapter, err := Open(dsn())
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback())
}

func TestAdapterQueryError(t *testing.T) {
	adapter, err := Open(dsn())
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	out := struct{}{}

	_, err = adapter.Query(&out, ";;", []interface{}{})
	assert.NotNil(t, err)
}

func TestAdapterExecError(t *testing.T) {
	adapter, err := Open(dsn())
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	_, _, err = adapter.Exec(";;", []interface{}{})
	assert.NotNil(t, err)
}

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

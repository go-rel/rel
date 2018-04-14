package mysql

import (
	"os"
	"testing"
	"time"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/specs"
	"github.com/Fs02/grimoire/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func logger(string, time.Duration, error) {}

func init() {
	adapter, err := Open(dsn())
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS addresses;`, []interface{}{}, logger)
	paranoid.Panic(err)
	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS users;`, []interface{}{}, logger)
	paranoid.Panic(err)

	_, _, err = adapter.Exec(`CREATE TABLE users (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(30) NOT NULL,
		gender VARCHAR(10) NOT NULL,
		age INT NOT NULL,
		note varchar(50),
		created_at DATETIME,
		updated_at DATETIME
	);`, []interface{}{}, logger)
	paranoid.Panic(err)

	_, _, err = adapter.Exec(`CREATE TABLE addresses (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		user_id INT UNSIGNED,
		address VARCHAR(60) NOT NULL,
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`, []interface{}{}, logger)
	paranoid.Panic(err)
}

func dsn() string {
	if os.Getenv("MYSQL_DATABASE") != "" {
		return os.Getenv("MYSQL_DATABASE") + "?charset=utf8&parseTime=True&loc=Local"
	}

	return "root@(127.0.0.1:3306)/grimoire_test?charset=utf8&parseTime=True&loc=Local"
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
	specs.SaveUpdate(t, repo)

	// Delete specs
	specs.Delete(t, repo)

	// Transaction specs
	specs.Transaction(t, repo)
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

	_, err = adapter.InsertAll(grimoire.Repo{}.From("users"), fields, allchanges, logger)

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

	_, err = adapter.Query(&out, "error", []interface{}{}, logger)
	assert.NotNil(t, err)
}

func TestAdapterExecError(t *testing.T) {
	adapter, err := Open(dsn())
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	_, _, err = adapter.Exec("error", []interface{}{}, logger)
	assert.NotNil(t, err)
}

func TestErrorFunc(t *testing.T) {
	// error nil
	assert.Nil(t, errorFunc(nil))

	// 1062 error
	rawerr := &mysql.MySQLError{Message: "duplicate", Number: 1062}
	duperr := errors.DuplicateError(rawerr.Message, "")
	assert.Equal(t, duperr, errorFunc(rawerr))

	// other errors
	err := errors.UnexpectedError("error")
	assert.Equal(t, err, errorFunc(err))
}

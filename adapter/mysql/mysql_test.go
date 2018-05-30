package mysql

import (
	"os"
	"testing"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/specs"
	"github.com/stretchr/testify/assert"
)

func init() {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")

	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS extras;`, nil)
	paranoid.Panic(err, "failed dropping extras table")
	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS addresses;`, nil)
	paranoid.Panic(err, "failed dropping addresses table")
	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS users;`, nil)
	paranoid.Panic(err, "failed dropping users table")

	_, _, err = adapter.Exec(`CREATE TABLE users (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(30) NOT NULL,
		gender VARCHAR(10) NOT NULL,
		age INT NOT NULL,
		note varchar(50),
		created_at DATETIME,
		updated_at DATETIME
	);`, nil)
	paranoid.Panic(err, "failed creating users table")

	_, _, err = adapter.Exec(`CREATE TABLE addresses (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		user_id INT UNSIGNED,
		address VARCHAR(60) NOT NULL,
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`, nil)
	paranoid.Panic(err, "failed creating addresses table")

	_, _, err = adapter.Exec(`CREATE TABLE extras (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		slug VARCHAR(30) DEFAULT NULL UNIQUE,
		user_id INT UNSIGNED,
		SCORE INT,
		CONSTRAINT extras_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
	);`, nil)

	paranoid.Panic(err, "failed creating extras table")
}

func dsn() string {
	if os.Getenv("MYSQL_DATABASE") != "" {
		return os.Getenv("MYSQL_DATABASE") + "?charset=utf8&parseTime=True&loc=Local"
	}

	return "root@(127.0.0.1:3306)/grimoire_test?charset=utf8&parseTime=True&loc=Local"
}

func TestAdapter__specs(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	repo := grimoire.New(adapter)

	// Query Specs
	specs.Query(t, repo)
	specs.QueryJoin(t, repo)
	specs.QueryNotFound(t, repo)

	// Preload specs
	specs.Preload(t, repo)

	// Count Specs
	specs.Count(t, repo)

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

	// Constraint specs
	// - Check constraint is not supported by mysql
	specs.UniqueConstraint(t, repo)
	specs.ForeignKeyConstraint(t, repo)
}

func TestAdapter_InsertAll_error(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	fields := []string{"notexist"}
	allchanges := []map[string]interface{}{
		{"notexist": "12"},
		{"notexist": "13"},
	}

	_, err = adapter.InsertAll(grimoire.Repo{}.From("users"), fields, allchanges)

	assert.NotNil(t, err)
}

func TestAdapter_Transaction_commitError(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit())
}

func TestAdapter_Transaction_rollbackError(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback())
}

func TestAdapter_Query_error(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	out := struct{}{}

	_, err = adapter.Query(&out, "error", nil)
	assert.NotNil(t, err)
}

func TestAdapter_Exec_error(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	_, _, err = adapter.Exec("error", nil)
	assert.NotNil(t, err)
}

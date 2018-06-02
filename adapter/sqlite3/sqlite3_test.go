package sqlite3

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
	defer adapter.Close()

	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS extras;`, nil)
	paranoid.Panic(err, "failed when dropping extras table")
	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS addresses;`, nil)
	paranoid.Panic(err, "failed when dropping addresses table")
	_, _, err = adapter.Exec(`DROP TABLE IF EXISTS users;`, nil)
	paranoid.Panic(err, "failed when dropping users table")

	_, _, err = adapter.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		slug VARCHAR(30) DEFAULT NULL,
		name VARCHAR(30) NOT NULL DEFAULT '',
		gender VARCHAR(10) NOT NULL DEFAULT 'male',
		age INTEGER NOT NULL DEFAULT 0,
		note varchar(50),
		created_at DATETIME,
		updated_at DATETIME,
		UNIQUE (slug)
	);`, nil)
	paranoid.Panic(err, "failed when creating users table")

	_, _, err = adapter.Exec(`CREATE TABLE addresses (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		address VARCHAR(60) NOT NULL DEFAULT '',
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`, nil)
	paranoid.Panic(err, "failed when creating addresses table")

	_, _, err = adapter.Exec(`CREATE TABLE extras (
		id INTEGER PRIMARY KEY,
		slug VARCHAR(30) DEFAULT NULL UNIQUE,
		user_id INTEGER,
		score INTEGER DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES users(id),
		CONSTRAINT extras_score_check CHECK (score>=0 AND score<=100)
	);`, nil)
	paranoid.Panic(err, "failed when creating extras table")
}

func dsn() string {
	if os.Getenv("SQLITE3_DATABASE") != "" {
		return os.Getenv("SQLITE3_DATABASE") + "?_foreign_keys=1"
	}

	return "./grimoire_test.db?_foreign_keys=1"
}

func TestAdapter_specs(t *testing.T) {
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

	// Aggregate Specs
	specs.Aggregate(t, repo)

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
	// - foreign key constraint is not supported because of lack of information in the error message.
	specs.UniqueConstraint(t, repo)
	specs.CheckConstraint(t, repo)
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

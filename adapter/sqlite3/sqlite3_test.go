package sqlite3

import (
	"context"
	"os"
	"testing"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/specs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

func init() {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS extras;`, nil)
	paranoid.Panic(err, "failed when dropping extras table")
	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS addresses;`, nil)
	paranoid.Panic(err, "failed when dropping addresses table")
	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS users;`, nil)
	paranoid.Panic(err, "failed when dropping users table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		slug VARCHAR(30) DEFAULT NULL,
		name VARCHAR(30) NOT NULL DEFAULT '',
		gender VARCHAR(10) NOT NULL DEFAULT '',
		age INTEGER NOT NULL DEFAULT 0,
		note varchar(50),
		created_at DATETIME,
		updated_at DATETIME,
		UNIQUE (slug)
	);`, nil)
	paranoid.Panic(err, "failed when creating users table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE addresses (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		name VARCHAR(60) NOT NULL DEFAULT '',
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`, nil)
	paranoid.Panic(err, "failed when creating addresses table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE extras (
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
		return os.Getenv("SQLITE3_DATABASE") + "?_foreign_keys=1&_loc=Local"
	}

	return "./rel_test.db?_foreign_keys=1&_loc=Local"
}

func TestAdapter_specs(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	repo := rel.New(adapter)

	// Query Specs
	specs.Query(t, repo)
	specs.QueryJoin(t, repo)
	specs.QueryNotFound(t, repo)

	// Preload specs
	specs.PreloadHasMany(t, repo)
	specs.PreloadHasManyWithQuery(t, repo)
	specs.PreloadHasManySlice(t, repo)
	specs.PreloadHasOne(t, repo)
	specs.PreloadHasOneWithQuery(t, repo)
	specs.PreloadHasOneSlice(t, repo)
	specs.PreloadBelongsTo(t, repo)
	specs.PreloadBelongsToWithQuery(t, repo)
	specs.PreloadBelongsToSlice(t, repo)

	// Aggregate Specs
	specs.Aggregate(t, repo)

	// Insert Specs
	specs.Insert(t, repo)
	specs.InsertHasMany(t, repo)
	specs.InsertHasOne(t, repo)
	specs.InsertBelongsTo(t, repo)
	specs.Inserts(t, repo)
	specs.InsertAll(t, repo)

	// Update Specs
	specs.Update(t, repo)
	specs.UpdateNotFound(t, repo)
	specs.UpdateHasManyInsert(t, repo)
	specs.UpdateHasManyUpdate(t, repo)
	specs.UpdateHasManyReplace(t, repo)
	specs.UpdateHasOneInsert(t, repo)
	specs.UpdateHasOneUpdate(t, repo)
	specs.UpdateBelongsToInsert(t, repo)
	specs.UpdateBelongsToUpdate(t, repo)
	specs.UpdateAtomic(t, repo)
	specs.Updates(t, repo)

	// Delete specs
	specs.Delete(t, repo)
	specs.DeleteAll(t, repo)

	// Constraint specs
	// - foreign key constraint is not supported because of lack of information in the error message.
	specs.UniqueConstraint(t, repo)
	specs.CheckConstraint(t, repo)
}

// func TestAdapter_InsertAll_error(t *testing.T) {
// 	adapter, err := Open(dsn())
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	fields := []string{"notexist"}
// 	mutations := []map[string]interface{}{
// 		{"notexist": "12"},
// 		{"notexist": "13"},
// 	}

// 	_, err = adapter.InsertAll(rel.Repo{}.From("users"), fields, mutations)

// 	assert.NotNil(t, err)
// }

func TestAdapter_Transaction_commitError(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit(ctx))
}

func TestAdapter_Transaction_rollbackError(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback(ctx))
}

// func TestAdapter_Query_error(t *testing.T) {
// 	adapter, err := Open(dsn())
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	out := struct{}{}

// 	_, err = adapter.Query(&out, "error", nil)
// 	assert.NotNil(t, err)
// }

func TestAdapter_Exec_error(t *testing.T) {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	_, _, err = adapter.Exec(ctx, "error", nil)
	assert.NotNil(t, err)
}

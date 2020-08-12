package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/specs"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

func init() {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")
	defer adapter.Close()

	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS extras;`, nil)
	paranoid.Panic(err, "failed dropping extras table")
	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS addresses;`, nil)
	paranoid.Panic(err, "failed dropping addresses table")
	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS users;`, nil)
	paranoid.Panic(err, "failed dropping users table")
	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS composites;`, nil)
	paranoid.Panic(err, "failed dropping composites table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE users (
		id SERIAL NOT NULL PRIMARY KEY,
		slug VARCHAR(30) DEFAULT NULL,
		name VARCHAR(30) NOT NULL DEFAULT '',
		gender VARCHAR(10) NOT NULL DEFAULT '',	
		age INT NOT NULL DEFAULT 0,
		note varchar(50),
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ,
		UNIQUE(slug)
	);`, nil)
	paranoid.Panic(err, "failed creating users table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE addresses (
		id SERIAL NOT NULL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		name VARCHAR(60) NOT NULL DEFAULT '',
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ
	);`, nil)
	paranoid.Panic(err, "failed creating addresses table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE extras (
		id SERIAL NOT NULL PRIMARY KEY,
		slug VARCHAR(30) DEFAULT NULL UNIQUE,
		user_id INTEGER REFERENCES users(id),
		score INTEGER DEFAULT 0 CHECK (score>=0 AND score<=100)
	);`, nil)
	paranoid.Panic(err, "failed creating extras table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE composites (
		primary1 SERIAL NOT NULL,
		primary2 SERIAL NOT NULL,
		data VARCHAR(255) DEFAULT NULL,
		PRIMARY KEY (primary1, primary2)
	);`, nil)
	paranoid.Panic(err, "failed creating composites table")

	// hack to make sure location it has the same location object as returned by pq driver.
	time.Local, err = time.LoadLocation("Asia/Jakarta")
	paranoid.Panic(err, "failed loading time location")
}

func dsn() string {
	if os.Getenv("POSTGRESQL_DATABASE") != "" {
		return os.Getenv("POSTGRESQL_DATABASE") + "?sslmode=disable&timezone=Asia/Jakarta"
	}

	return "postgres://rel@localhost:9920/rel_test?sslmode=disable&timezone=Asia/Jakarta"
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
	specs.InsertAllPartialCustomPrimary(t, repo)

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
	specs.UpdateAll(t, repo)

	// Delete specs
	specs.Delete(t, repo)
	specs.DeleteBelongsTo(t, repo)
	specs.DeleteHasOne(t, repo)
	specs.DeleteHasMany(t, repo)
	specs.DeleteAll(t, repo)

	// Constraint specs
	specs.UniqueConstraint(t, repo)
	specs.ForeignKeyConstraint(t, repo)
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

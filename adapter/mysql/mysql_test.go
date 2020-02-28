package mysql

import (
	"context"
	"os"
	"testing"

	paranoid "github.com/Fs02/go-paranoid"
	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/specs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

func init() {
	adapter, err := Open(dsn())
	paranoid.Panic(err, "failed to open database connection")

	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS extras;`, nil)
	paranoid.Panic(err, "failed dropping extras table")
	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS addresses;`, nil)
	paranoid.Panic(err, "failed dropping addresses table")
	_, _, err = adapter.Exec(ctx, `DROP TABLE IF EXISTS users;`, nil)
	paranoid.Panic(err, "failed dropping users table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE users (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(30) NOT NULL DEFAULT '',
		gender VARCHAR(10) NOT NULL DEFAULT '',
		age INT NOT NULL DEFAULT 0,
		note varchar(50),
		created_at DATETIME,
		updated_at DATETIME
	);`, nil)
	paranoid.Panic(err, "failed creating users table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE addresses (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		user_id INT UNSIGNED,
		name VARCHAR(60) NOT NULL DEFAULT '',
		created_at DATETIME,
		updated_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`, nil)
	paranoid.Panic(err, "failed creating addresses table")

	_, _, err = adapter.Exec(ctx, `CREATE TABLE extras (
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

	return "root@tcp(localhost:3306)/rel_test?charset=utf8&parseTime=True&loc=Local"
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
	// - Check constraint is not supported by mysql
	specs.UniqueConstraint(t, repo)
	specs.ForeignKeyConstraint(t, repo)
}

func TestAdapter_Open(t *testing.T) {
	// with parameter
	assert.NotPanics(t, func() {
		adapter, _ := Open("root@tcp(localhost:3306)/rel_test?charset=utf8")
		defer adapter.Close()
	})

	// without paremeter
	assert.NotPanics(t, func() {
		adapter, _ := Open("root@tcp(localhost:3306)/rel_test")
		defer adapter.Close()
	})
}

// func TestAdapter_InsertAll_error(t *testing.T) {
// 	adapter, err := Open(dsn())
// 	paranoid.Panic(err, "failed to open database connection")
// 	defer adapter.Close()

// 	fields := []string{"notexist"}
// 	modifications := []map[string]interface{}{
// 		{"notexist": "12"},
// 		{"notexist": "13"},
// 	}

// 	_, err = adapter.InsertAll(rel.Repo{}.From("users"), fields, modifications)

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

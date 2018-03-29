package mysql

import (
	"os"
	"testing"
	"time"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

func init() {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	adapter.Exec(`DROP TABLE IF EXISTS users;`, []interface{}{})
	adapter.Exec(`CREATE TABLE users (
		id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(30) NOT NULL,
		created_at DATETIME,
		updated_at DATETIME
	);`, []interface{}{})
}

type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func dsn() string {
	if os.Getenv("MYSQL_DATABASE") != "" {
		return os.Getenv("MYSQL_DATABASE")
	}

	return "root@(127.0.0.1:3306)/grimoire_test"
}

func TestRepoFind(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	stmt := "INSERT INTO users (name, created_at, updated_at) VALUES (?,?,?)"
	name := "find"
	createdAt := time.Now().Round(time.Second)
	updatedAt := time.Now().Round(time.Second)

	id, count, err := adapter.Exec(stmt, []interface{}{name, createdAt, updatedAt})
	assert.Nil(t, err)
	assert.True(t, id > 0)
	assert.Equal(t, int64(1), count)

	repo := grimoire.New(adapter)
	user := User{}
	users := []User{}

	// find inserted user
	err = repo.From("users").Find(id).One(&user)
	assert.Nil(t, err)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, createdAt, user.CreatedAt)
	assert.Equal(t, updatedAt, user.UpdatedAt)

	// find all user
	err = repo.From("users").Find(id).All(&users)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, id, users[0].ID)
	assert.Equal(t, name, users[0].Name)
	assert.Equal(t, createdAt, users[0].CreatedAt)
	assert.Equal(t, updatedAt, users[0].UpdatedAt)

	// find user error not found
	err = repo.From("users").Find(0).One(&user)
	assert.True(t, err.(errors.Error).NotFoundError())
}

func TestRepoInsert(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	name := "insert"
	createdAt := time.Now().Round(time.Second)
	updatedAt := time.Now().Round(time.Second)
	ch := changeset.Cast(map[string]interface{}{
		"name":       name,
		"created_at": createdAt,
		"updated_at": updatedAt,
	}, []string{"name", "created_at", "updated_at"})

	user := User{}
	err := grimoire.New(adapter).From("users").Insert(&user, ch)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, createdAt, user.CreatedAt)
	assert.Equal(t, updatedAt, user.UpdatedAt)
}

func TestRepoUpdate(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	stmt := "INSERT INTO users (name, created_at, updated_at) VALUES (?,?,?)"
	name := "update"
	createdAt := time.Now().Round(time.Second)
	updatedAt := time.Now().Round(time.Second)

	id, count, err := adapter.Exec(stmt, []interface{}{name, createdAt, updatedAt})
	assert.Nil(t, err)
	assert.True(t, id > 0)
	assert.Equal(t, int64(1), count)

	newName := "new update"
	ch := changeset.Cast(map[string]interface{}{
		"name": newName,
	}, []string{"name"})
	user := User{}
	err = grimoire.New(adapter).From("users").Find(id).Update(&user, ch)
	assert.Nil(t, err)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, newName, user.Name)
	assert.Equal(t, createdAt, user.CreatedAt)
	assert.Equal(t, updatedAt, user.UpdatedAt)
}

func TestRepoDelete(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	stmt := "INSERT INTO users (name, created_at, updated_at) VALUES (?,?,?)"
	name := "delete"
	createdAt := time.Now().Round(time.Second)
	updatedAt := time.Now().Round(time.Second)

	id, count, err := adapter.Exec(stmt, []interface{}{name, createdAt, updatedAt})
	assert.Nil(t, err)
	assert.True(t, id > 0)
	assert.Equal(t, int64(1), count)

	err = grimoire.New(adapter).From("users").Find(id).Delete()
	assert.Nil(t, err)

	user := User{}
	err = grimoire.New(adapter).From("users").Find(id).One(&user)
	assert.NotNil(t, err)
}

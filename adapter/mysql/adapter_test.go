package mysql

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Fs02/grimoire"
	. "github.com/Fs02/grimoire/c"
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
		gender VARCHAR(10) NOT NULL,
		age INT NOT NULL,
		note varchar(50),
		created_at DATETIME,
		updated_at DATETIME
	);`, []interface{}{})

	// prepare data
	repo := grimoire.New(adapter)
	data := repo.From("users").
		Set("gender", "male").
		Set("created_at", time.Now()).
		Set("updated_at", time.Now())

	for i := 1; i <= 5; i++ {
		data.Set("id", i).
			Set("name", "name"+strconv.Itoa(i)).
			Set("age", i*10).
			MustInsert(nil)
	}

	data = data.Set("gender", "female")
	for i := 6; i <= 8; i++ {
		data.Set("id", i).
			Set("name", "name"+strconv.Itoa(i)).
			Set("age", i*10).
			MustInsert(nil)
	}
}

type User struct {
	ID        int64
	Name      string
	Gender    string
	Age       int
	Note      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User table identifiers
const (
	users     = "users"
	id        = I("id")
	name      = I("name")
	gender    = I("gender")
	age       = I("age")
	note      = I("note")
	createdAt = I("created_at")
)

func dsn() string {
	if os.Getenv("MYSQL_DATABASE") != "" {
		return os.Getenv("MYSQL_DATABASE")
	}

	return "root@(127.0.0.1:3306)/grimoire_test"
}

func TestRepoQuery(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	repo := grimoire.New(adapter)

	queries := []grimoire.Query{
		repo.From(users).Where(Eq(id, 1)),
		repo.From(users).Where(Eq(name, "name1")),
		repo.From(users).Where(Eq(age, 10)),
		repo.From(users).Where(Eq(id, 1), Eq(name, "name1")),
		repo.From(users).Where(Eq(id, 1), Eq(name, "name1"), Eq(age, 10)),
		repo.From(users).Where(Eq(id, 1)).OrWhere(Eq(name, "name1")),
		repo.From(users).Where(Eq(id, 1)).OrWhere(Eq(name, "name1"), Eq(age, 10)),
		repo.From(users).Where(Eq(id, 1)).OrWhere(Eq(name, "name1")).OrWhere(Eq(age, 10)),
		repo.From(users).Where(Ne(gender, "male")),
		repo.From(users).Where(Gt(age, 79)),
		repo.From(users).Where(Gte(age, 80)),
		repo.From(users).Where(Lt(age, 11)),
		repo.From(users).Where(Lte(age, 10)),
		repo.From(users).Where(Nil(note)),
		repo.From(users).Where(NotNil(name)),
		repo.From(users).Order(Asc(name)),
		repo.From(users).Order(Asc(name), Desc(age)),
		repo.From(users).Where(In(id, 1, 2, 3)),
		repo.From(users).Where(Nin(id, 1, 2, 3)),
		repo.From(users).Where(Like(name, "name%")),
		repo.From(users).Where(NotLike(name, "noname%")),
		repo.From(users).Where(Fragment("id = ?", 1)),
		repo.From(users).Where(Not(Eq(id, 1), Eq(name, "name1"), Eq(age, 10))),
		repo.From(users).Where(Xor(Eq(id, 1), Eq(name, "name1"), Eq(age, 10))),
		repo.From(users).Limit(5),
		repo.From(users).Limit(5).Offset(5),
		repo.From(users).Find(1),
		repo.From(users).Select("name").Find(1),
		repo.From(users).Select("name", "age").Find(1),
		repo.From(users).Distinct().Find(1),
	}

	for _, q := range queries {
		str, _ := adapter.Find(q)
		t.Run("ALL|"+str, func(t *testing.T) {
			var result []User
			err := q.All(&result)
			assert.Nil(t, err)
			assert.NotEqual(t, 0, len(result))
		})
	}
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

	user := User{}
	name := "insert"

	ch := changeset.Cast(user, map[string]interface{}{
		"name": name,
	}, []string{"name", "created_at", "updated_at"})

	err := grimoire.New(adapter).From("users").Insert(&user, ch)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, name, user.Name)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
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

	user := User{}
	newName := "new update"

	ch := changeset.Cast(user, map[string]interface{}{
		"name": newName,
	}, []string{"name"})

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

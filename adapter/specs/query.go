package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

// Query tests query specifications without join.
func Query(t *testing.T, repo grimoire.Repo) {
	// preparte tests data
	user := User{Name: "name1", Gender: "male", Age: 10}
	repo.From(users).MustSave(&user)
	repo.From(users).MustSave(&User{Name: "name2", Gender: "male", Age: 20})
	repo.From(users).MustSave(&User{Name: "name3", Gender: "male", Age: 30})
	repo.From(users).MustSave(&User{Name: "name4", Gender: "female", Age: 40})
	repo.From(users).MustSave(&User{Name: "name5", Gender: "female", Age: 50})
	repo.From(users).MustSave(&User{Name: "name6", Gender: "female", Age: 60})

	repo.From(addresses).MustSave(&Address{Address: "address1", UserID: &user.ID})
	repo.From(addresses).MustSave(&Address{Address: "address2", UserID: &user.ID})
	repo.From(addresses).MustSave(&Address{Address: "address3", UserID: &user.ID})

	tests := []grimoire.Query{
		repo.From(users).Where(c.Eq(id, user.ID)),
		repo.From(users).Where(c.Eq(name, "name1")),
		repo.From(users).Where(c.Eq(age, 10)),
		repo.From(users).Where(c.Eq(id, user.ID), c.Eq(name, "name1")),
		repo.From(users).Where(c.Eq(id, user.ID), c.Eq(name, "name1"), c.Eq(age, 10)),
		repo.From(users).Where(c.Eq(id, user.ID)).OrWhere(c.Eq(name, "name1")),
		repo.From(users).Where(c.Eq(id, user.ID)).OrWhere(c.Eq(name, "name1"), c.Eq(age, 10)),
		repo.From(users).Where(c.Eq(id, user.ID)).OrWhere(c.Eq(name, "name1")).OrWhere(c.Eq(age, 10)),
		repo.From(users).Where(c.Ne(gender, "male")),
		repo.From(users).Where(c.Gt(age, 59)),
		repo.From(users).Where(c.Gte(age, 60)),
		repo.From(users).Where(c.Lt(age, 11)),
		repo.From(users).Where(c.Lte(age, 10)),
		repo.From(users).Where(c.Nil(note)),
		repo.From(users).Where(c.NotNil(name)),
		repo.From(users).Where(c.In(id, 1, 2, 3)),
		repo.From(users).Where(c.Nin(id, 1, 2, 3)),
		repo.From(users).Where(c.Like(name, "name%")),
		repo.From(users).Where(c.NotLike(name, "noname%")),
		repo.From(users).Where(c.Fragment("id > 0")),
		repo.From(users).Where(c.Not(c.Eq(id, 1), c.Eq(name, "name1"), c.Eq(age, 10))),
		repo.From(users).Order(c.Asc(name)),
		repo.From(users).Order(c.Desc(name)),
		repo.From(users).Order(c.Asc(name), c.Desc(age)),
		repo.From(users).Group("gender").Select("COUNT(id)"),
		repo.From(users).Group("age").Having(c.Gt(age, 10)).Select("COUNT(id)"),
		repo.From(users).Limit(5),
		repo.From(users).Limit(5).Offset(5),
		repo.From(users).Find(1),
		repo.From(users).Select("name").Find(1),
		repo.From(users).Select("name", "age").Find(1),
		repo.From(users).Distinct().Find(1),
	}

	run(t, tests)
}

// QueryJoin tests query specifications with join.
func QueryJoin(t *testing.T, repo grimoire.Repo) {
	tests := []grimoire.Query{
		repo.From(addresses).Join(users),
		repo.From(addresses).Join(users, c.Eq(c.I("addresses.user_id"), c.I("users.id"))),
		repo.From(addresses).Join(users).Find(1),
		repo.From(addresses).Join(users).Where(c.Eq(address, "address1")),
		repo.From(addresses).Join(users).Where(c.Eq(address, "address1")).Order(c.Asc(name)),
		repo.From(addresses).JoinWith("LEFT JOIN", users),
		repo.From(addresses).JoinWith("LEFT OUTER JOIN", users),
	}

	run(t, tests)
}

// QueryNotFound tests query specifications when no result found.
func QueryNotFound(t *testing.T, repo grimoire.Repo) {
	t.Run("NotFound", func(t *testing.T) {
		user := User{}

		// find user error not found
		err := repo.From("users").Find(0).One(&user)
		assert.NotNil(t, err)
		assert.Equal(t, errors.NotFound, err.(errors.Error).Kind())
	})
}

func run(t *testing.T, queries []grimoire.Query) {
	for _, query := range queries {
		statement, _ := builder.Find(query)
		t.Run("All|"+statement, func(t *testing.T) {
			var result []User
			assert.Nil(t, query.All(&result))
			assert.NotEqual(t, 0, len(result))
		})
	}

	for _, query := range queries {
		statement, _ := builder.Find(query)
		t.Run("One|"+statement, func(t *testing.T) {
			var result User
			assert.Nil(t, query.One(&result))
		})
	}
}

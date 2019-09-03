package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/sort"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

// Query tests query specifications without join.
func Query(t *testing.T, repo grimoire.Repo) {
	// preparte tests data
	var (
		user = User{Name: "name1", Gender: "male", Age: 10}
	)

	repo.MustInsert(&user)
	repo.MustInsert(&User{Name: "name2", Gender: "male", Age: 20})
	repo.MustInsert(&User{Name: "name3", Gender: "male", Age: 30})
	repo.MustInsert(&User{Name: "name4", Gender: "female", Age: 40})
	repo.MustInsert(&User{Name: "name5", Gender: "female", Age: 50})
	repo.MustInsert(&User{Name: "name6", Gender: "female", Age: 60})

	repo.MustInsert(&Address{Name: "address1", UserID: &user.ID})
	repo.MustInsert(&Address{Name: "address2", UserID: &user.ID})
	repo.MustInsert(&Address{Name: "address3", UserID: &user.ID})

	tests := []grimoire.Querier{
		where.Eq("id", user.ID),
		grimoire.Where(where.Eq("id", user.ID)),
		grimoire.Where(where.Eq("name", "name1")),
		grimoire.Where(where.Eq("age", 10)),
		grimoire.Where(where.Eq("id", user.ID), where.Eq("name", "name1")),
		grimoire.Where(where.Eq("id", user.ID), where.Eq("name", "name1"), where.Eq("age", 10)),
		grimoire.Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1")),
		grimoire.Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1"), where.Eq("age", 10)),
		grimoire.Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1")).OrWhere(where.Eq("age", 10)),
		grimoire.Where(where.Ne("gender", "male")),
		grimoire.Where(where.Gt("age", 59)),
		grimoire.Where(where.Gte("age", 60)),
		grimoire.Where(where.Lt("age", 11)),
		grimoire.Where(where.Lte("age", 10)),
		grimoire.Where(where.Nil("note")),
		grimoire.Where(where.NotNil("name")),
		grimoire.Where(where.In("id", 1, 2, 3)),
		grimoire.Where(where.Nin("id", 1, 2, 3)),
		grimoire.Where(where.Like("name", "name%")),
		grimoire.Where(where.NotLike("name", "noname%")),
		grimoire.Where(where.Fragment("id > 0")),
		grimoire.Where(where.Not(where.Eq("id", 1), where.Eq("name", "name1"), where.Eq("age", 10))),
		sort.Asc("name"),
		sort.Desc("name"),
		grimoire.Select().SortAsc("name").SortDesc("age"),
		grimoire.Group("gender").Select("COUNT(id)"),
		grimoire.Group("age").Having(where.Gt("age", 10)).Select("COUNT(id)"),
		grimoire.Limit(5),
		grimoire.Select().Limit(5),
		grimoire.Select().Limit(5).Offset(5),
		grimoire.Select("name").Where(where.Eq("id", 1)),
		grimoire.Select("name", "age").Where(where.Eq("id", 1)),
		grimoire.Select().Distinct().Where(where.Eq("id", 1)),
	}

	run(t, repo, tests)
}

// QueryJoin tests query specifications with join.
func QueryJoin(t *testing.T, repo grimoire.Repo) {
	tests := []grimoire.Querier{
		// grimoire.Join("users"),
		// grimoire.Join(users, where.Eq(where.I("addresses.user_id"), where.I("users.id"))),
		// grimoire.Join(users).Where(where.Eq("id", 1)),
		// grimoire.Join(users).Where(where.Eq(address, "address1")),
		// grimoire.Join(users).Where(where.Eq(address, "address1")).Sort("name"),
		// grimoire.JoinWith("LEFT JOIN", users),
		// grimoire.JoinWith("LEFT OUTER JOIN", users),
	}

	run(t, repo, tests)
}

// QueryNotFound tests query specifications when no result found.
func QueryNotFound(t *testing.T, repo grimoire.Repo) {
	t.Run("NotFound", func(t *testing.T) {
		var (
			user User
			err  = repo.One(&user, where.Eq("id", 0))
		)

		// find user error not found
		assert.Equal(t, grimoire.NoResultError{}, err)
	})
}

func run(t *testing.T, repo grimoire.Repo, queriers []grimoire.Querier) {
	for _, query := range queriers {
		t.Run("All", func(t *testing.T) {
			var (
				users []User
				err   = repo.All(&users, query)
			)

			assert.Nil(t, err)
			assert.NotEqual(t, 0, len(users))
		})
	}

	for _, query := range queriers {
		t.Run("One", func(t *testing.T) {
			var (
				user User
				err  = repo.One(&user, query)
			)

			assert.Nil(t, err)
		})
	}
}

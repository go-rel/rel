package specs

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/sort"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

// Query tests query specifications without join.
func Query(t *testing.T, repo rel.Repository) {
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

	tests := []rel.Querier{
		where.Eq("id", user.ID),
		rel.Where(where.Eq("id", user.ID)),
		rel.Where(where.Eq("name", "name1")),
		rel.Where(where.Eq("age", 10)),
		rel.Where(where.Eq("id", user.ID), where.Eq("name", "name1")),
		rel.Where(where.Eq("id", user.ID), where.Eq("name", "name1"), where.Eq("age", 10)),
		rel.Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1")),
		rel.Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1"), where.Eq("age", 10)),
		rel.Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1")).OrWhere(where.Eq("age", 10)),
		rel.Where(where.Ne("gender", "male")),
		rel.Where(where.Gt("age", 59)),
		rel.Where(where.Gte("age", 60)),
		rel.Where(where.Lt("age", 11)),
		rel.Where(where.Lte("age", 10)),
		rel.Where(where.Nil("note")),
		rel.Where(where.NotNil("name")),
		rel.Where(where.In("id", 1, 2, 3)),
		rel.Where(where.Nin("id", 1, 2, 3)),
		rel.Where(where.Like("name", "name%")),
		rel.Where(where.NotLike("name", "noname%")),
		rel.Where(where.Fragment("id > 0")),
		rel.Where(where.Not(where.Eq("id", 1), where.Eq("name", "name1"), where.Eq("age", 10))),
		sort.Asc("name"),
		sort.Desc("name"),
		rel.Select().SortAsc("name").SortDesc("age"),
		rel.Select("gender", "COUNT(id) AS count").Group("gender"),
		rel.Select("age", "COUNT(id) AS count").Group("age").Having(where.Gt("age", 10)),
		rel.Limit(5),
		rel.Select().Limit(5),
		rel.Select().Limit(5).Offset(5),
		rel.Select("name").Where(where.Eq("id", 1)),
		rel.Select("name", "age").Where(where.Eq("id", 1)),
		rel.Select().Distinct().Where(where.Eq("id", 1)),
	}

	run(t, repo, tests)
}

// QueryJoin tests query specifications with join.
func QueryJoin(t *testing.T, repo rel.Repository) {
	tests := []rel.Querier{
		rel.From("addresses").Join("users"),
		rel.From("addresses").JoinOn("users", "addresses.user_id", "users.id"),
		rel.From("addresses").Join("users").Where(where.Eq("addresses.id", 1)),
		rel.From("addresses").Join("users").Where(where.Eq("addresses.name", "address1")),
		rel.From("addresses").Join("users").Where(where.Eq("addresses.name", "address1")).SortAsc("addresses.name"),
		rel.From("addresses").JoinWith("LEFT JOIN", "users", "addresses.user_id", "users.id"),
	}

	run(t, repo, tests)
}

// QueryNotFound tests query specifications when no result found.
func QueryNotFound(t *testing.T, repo rel.Repository) {
	t.Run("NotFound", func(t *testing.T) {
		var (
			user User
			err  = repo.Find(&user, where.Eq("id", 0))
		)

		// find user error not found
		assert.Equal(t, rel.NoResultError{}, err)
	})
}

func run(t *testing.T, repo rel.Repository, queriers []rel.Querier) {
	for _, query := range queriers {
		t.Run("FindAll", func(t *testing.T) {
			var (
				users []User
				err   = repo.FindAll(&users, query)
			)

			assert.Nil(t, err)
			assert.NotEqual(t, 0, len(users))
		})
	}

	for _, query := range queriers {
		t.Run("Find", func(t *testing.T) {
			var (
				user User
				err  = repo.Find(&user, query)
			)

			assert.Nil(t, err)
		})
	}
}

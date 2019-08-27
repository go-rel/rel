package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

// Aggregate tests count specifications.
func Aggregate(t *testing.T, repo grimoire.Repo) {
	// preparte tests data
	var (
		user = User{Name: "name1", Gender: "male", Age: 10}
	)

	repo.MustInsert(&user)

	tests := []grimoire.Query{
		grimoire.From("users").Where(where.Eq("id", user.ID)),
		grimoire.From("users").Where(where.Eq("name", "name1")),
		grimoire.From("users").Where(where.Eq("age", 10)),
		grimoire.From("users").Where(where.Eq("id", user.ID), where.Eq("name", "name1")),
		grimoire.From("users").Where(where.Eq("id", user.ID), where.Eq("name", "name1"), where.Eq("age", 10)),
		grimoire.From("users").Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1")),
		grimoire.From("users").Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1"), where.Eq("age", 10)),
		grimoire.From("users").Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1")).OrWhere(where.Eq("age", 10)),
		grimoire.From("users").Where(where.Ne("gender", "male")),
		grimoire.From("users").Where(where.Gt("age", 59)),
		grimoire.From("users").Where(where.Gte("age", 60)),
		grimoire.From("users").Where(where.Lt("age", 11)),
		grimoire.From("users").Where(where.Lte("age", 10)),
		grimoire.From("users").Where(where.Nil("note")),
		grimoire.From("users").Where(where.NotNil("name")),
		grimoire.From("users").Where(where.In("id", 1, 2, 3)),
		grimoire.From("users").Where(where.Nin("id", 1, 2, 3)),
		grimoire.From("users").Where(where.Like("name", "name%")),
		grimoire.From("users").Where(where.NotLike("name", "noname%")),
		grimoire.From("users").Where(where.Fragment("id > 0")),
		grimoire.From("users").Where(where.Not(where.Eq("id", 1), where.Eq("name", "name1"), where.Eq("age", 10))),
		// this query is not supported.
		// group query is automatically removed.
		// use all instead for complex aggregation.
		grimoire.From("users").Limit(10),
		grimoire.From("users").Group("gender"),
		grimoire.From("users").Group("age").Having(where.Gt("age", 10)),
	}

	for _, query := range tests {
		statement, _ := builder.Find(query.Select("count(id) AS count"))

		t.Run("Aggregate|"+statement, func(t *testing.T) {
			count, err := repo.Aggregate(query, "count", "id")
			assert.Nil(t, err)
			assert.True(t, count >= 0)

			sum, err := repo.Aggregate(query, "sum", "id")
			assert.Nil(t, err)
			assert.True(t, sum >= 0)
		})
	}
}

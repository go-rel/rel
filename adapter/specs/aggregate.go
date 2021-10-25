package specs

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

// Aggregate tests count specifications.
func Aggregate(t *testing.T, repo rel.Repository) {
	// preparte tests data
	var (
		user = User{Name: "name1", Gender: "male", Age: 10}
	)

	repo.MustInsert(ctx, &user)

	waitForReplication()

	tests := []rel.Query{
		rel.From("users").Where(where.Eq("id", user.ID)),
		rel.From("users").Where(where.Eq("name", "name1")),
		rel.From("users").Where(where.Eq("age", 10)),
		rel.From("users").Where(where.Eq("id", user.ID), where.Eq("name", "name1")),
		rel.From("users").Where(where.Eq("id", user.ID), where.Eq("name", "name1"), where.Eq("age", 10)),
		rel.From("users").Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1")),
		rel.From("users").Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1"), where.Eq("age", 10)),
		rel.From("users").Where(where.Eq("id", user.ID)).OrWhere(where.Eq("name", "name1")).OrWhere(where.Eq("age", 10)),
		rel.From("users").Where(where.Ne("gender", "male")),
		rel.From("users").Where(where.Gt("age", 59)),
		rel.From("users").Where(where.Gte("age", 60)),
		rel.From("users").Where(where.Lt("age", 11)),
		rel.From("users").Where(where.Lte("age", 10)),
		rel.From("users").Where(where.Nil("note")),
		rel.From("users").Where(where.NotNil("name")),
		rel.From("users").Where(where.In("id", 1, 2, 3)),
		rel.From("users").Where(where.Nin("id", 1, 2, 3)),
		rel.From("users").Where(where.Like("name", "name%")),
		rel.From("users").Where(where.NotLike("name", "noname%")),
		rel.From("users").Where(where.Fragment("id > 0")),
		rel.From("users").Where(where.Not(where.Eq("id", 1), where.Eq("name", "name1"), where.Eq("age", 10))),
		// this query is not supported.
		// group query is automatically removed.
		// use all instead for complex aggregation.
		rel.From("users").Limit(10),
		rel.From("users").Group("gender"),
		rel.From("users").Group("age").Having(where.Gt("age", 10)),
	}

	for _, query := range tests {
		t.Run("Aggregate", func(t *testing.T) {
			count, err := repo.Aggregate(ctx, query, "count", "id")
			assert.Nil(t, err)
			assert.NotZero(t, count)

			sum, err := repo.Aggregate(ctx, query, "sum", "id")
			assert.Nil(t, err)
			assert.NotZero(t, sum)
		})
	}
}

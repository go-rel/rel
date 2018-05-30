package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
	"github.com/stretchr/testify/assert"
)

// Count tests count specifications.
func Count(t *testing.T, repo grimoire.Repo) {
	// preparte tests data
	user := User{Name: "name1", Gender: "male", Age: 10}
	repo.From(users).MustSave(&user)

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
	}

	for _, query := range tests {
		statement, _ := builder.Find(query.Select("COUNT(*) AS count"))
		t.Run("Count|"+statement, func(t *testing.T) {
			_, err := query.Count()
			assert.Nil(t, err)
		})
	}
}

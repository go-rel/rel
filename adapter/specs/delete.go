package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

// Delete tests delete specifications.
func Delete(t *testing.T, repo grimoire.Repo) {
	var (
		user = User{Name: "delete", Age: 100}
	)

	repo.MustInsert(&user)
	assert.NotEqual(t, 0, user.ID)

	assert.Nil(t, repo.Delete(&user))
	assert.NotNil(t, repo.One(&user, where.Eq("id", user.ID)))

}

// DeleteAll tests delete specifications.
func DeleteAll(t *testing.T, repo grimoire.Repo) {
	var (
		user = User{Name: "delete", Age: 100}
	)

	repo.MustInsert(&user)
	assert.NotEqual(t, 0, user.ID)

	assert.Nil(t, repo.Delete(&user))
	assert.NotNil(t, repo.One(&user, where.Eq("id", user.ID)))

	repo.MustInsert(&User{Name: "delete", Age: 100})
	repo.MustInsert(&User{Name: "delete", Age: 100})
	repo.MustInsert(&User{Name: "other delete", Age: 110})

	tests := []grimoire.Query{
		grimoire.From("users").Where(where.Eq("name", "delete")),
		grimoire.From("users").Where(where.Eq("name", "other delete"), where.Gt("age", 100)),
	}

	for _, query := range tests {
		statement, _ := builder.Delete(query.Collection, query.WhereQuery)
		t.Run("Delete|"+statement, func(t *testing.T) {
			var result []User
			assert.Nil(t, repo.All(&result, query))
			assert.NotEqual(t, 0, len(result))

			assert.Nil(t, repo.DeleteAll(query))

			assert.Nil(t, repo.All(&result, query))
			assert.Equal(t, 0, len(result))
		})
	}
}

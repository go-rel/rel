package specs

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

// Delete tests delete specifications.
func Delete(t *testing.T, repo rel.Repository) {
	var (
		user = User{Name: "delete", Age: 100}
	)

	repo.MustInsert(&user)
	assert.NotEqual(t, 0, user.ID)

	assert.Nil(t, repo.Delete(&user))
	assert.Equal(t, rel.NoResultError{}, repo.Find(&user, where.Eq("id", user.ID)))
}

// DeleteAll tests delete specifications.
func DeleteAll(t *testing.T, repo rel.Repository) {
	var (
		user = User{Name: "delete", Age: 100}
	)

	repo.MustInsert(&user)
	assert.NotEqual(t, 0, user.ID)

	assert.Nil(t, repo.Delete(&user))
	assert.NotNil(t, repo.Find(&user, where.Eq("id", user.ID)))

	repo.MustInsert(&User{Name: "delete", Age: 100})
	repo.MustInsert(&User{Name: "delete", Age: 100})
	repo.MustInsert(&User{Name: "other delete", Age: 110})

	tests := []rel.Query{
		rel.From("users").Where(where.Eq("name", "delete")),
		rel.From("users").Where(where.Eq("name", "other delete"), where.Gt("age", 100)),
	}

	for _, query := range tests {
		t.Run("Delete", func(t *testing.T) {
			var result []User
			assert.Nil(t, repo.FindAll(&result, query))
			assert.NotEqual(t, 0, len(result))

			assert.Nil(t, repo.DeleteAll(query))

			assert.Nil(t, repo.FindAll(&result, query))
			assert.Equal(t, 0, len(result))
		})
	}
}

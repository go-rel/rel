package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func Update(t *testing.T, repo grimoire.Repo) {
	var (
		name   = "zoro"
		gender = "male"
		age    = 23
		note   = "swordsman"
		user   = User{
			Name: "update",
		}
	)

	repo.MustInsert(&user)

	user.Name = name
	user.Gender = gender
	user.Age = age
	user.Note = &note

	err := repo.Update(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, gender, user.Gender)
	assert.Equal(t, age, user.Age)
	assert.Equal(t, &note, user.Note)

	var (
		queried User
	)

	user.Addresses = nil
	err = repo.One(&queried, where.Eq("id", user.ID))
	assert.Nil(t, err)
	assert.Equal(t, user, queried)
}

// Update tests update specifications.
// TODO: atomic update
// TODO: update all
// TODO: update with assoc
func Updates(t *testing.T, repo grimoire.Repo) {
	var (
		note    = "note"
		user    = User{Name: "update"}
		address = Address{Name: "update"}
	)

	repo.MustInsert(&user)
	repo.MustInsert(&address)

	tests := []interface{}{
		&User{ID: user.ID, Name: "changed", Age: 100},
		&User{ID: user.ID, Name: "changed", Age: 100, Note: &note},
		&User{ID: user.ID, Note: &note},
		&Address{ID: address.ID, Name: "address"},
		&Address{ID: address.ID, UserID: &user.ID},
		&Address{ID: address.ID, Name: "address", UserID: &user.ID},
	}

	for _, record := range tests {
		var (
			changes      = grimoire.BuildChanges(grimoire.Struct(record))
			statement, _ = builder.Update("collection", changes, where.Eq("id", 1))
		)

		t.Run("Update|"+statement, func(t *testing.T) {
			assert.Nil(t, repo.Update(record))
		})
	}
}

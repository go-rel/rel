package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func Update(t *testing.T, repo grimoire.Repo) {
	var (
		note = "swordsman"
		user = User{
			Name: "update",
		}
	)

	repo.MustInsert(&user)

	user.Name = "update"
	user.Gender = "male"
	user.Age = 23
	user.Note = &note

	err := repo.Update(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "update", user.Name)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, 23, user.Age)
	assert.Equal(t, &note, user.Note)

	var (
		queried User
	)

	user.Addresses = nil
	err = repo.One(&queried, where.Eq("id", user.ID))
	assert.Nil(t, err)
	assert.Equal(t, user, queried)
}

func UpdateHasManyInsert(t *testing.T, repo grimoire.Repo) {
	var (
		user = User{
			Name: "update init",
		}
	)

	repo.MustInsert(&user)

	user.Name = "update insert has many"
	user.Addresses = []Address{
		{Name: "primary"},
		{Name: "work"},
	}

	err := repo.Update(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "update insert has many", user.Name)

	assert.Len(t, user.Addresses, 2)
	assert.NotEqual(t, 0, user.Addresses[0].ID)
	assert.NotEqual(t, 0, user.Addresses[1].ID)
	assert.Equal(t, user.ID, *user.Addresses[0].UserID)
	assert.Equal(t, user.ID, *user.Addresses[1].UserID)
	assert.Equal(t, "primary", user.Addresses[0].Name)
	assert.Equal(t, "work", user.Addresses[1].Name)

	assert.Equal(t, 2, repo.MustCount("addresses", where.Eq("user_id", user.ID)))
}

func UpdateHasManyReplace(t *testing.T, repo grimoire.Repo) {
	var (
		user = User{
			Name: "update init",
			Addresses: []Address{
				{Name: "old address"},
			},
		}
	)

	repo.MustInsert(&user)

	user.Name = "update insert has many"
	user.Addresses = []Address{
		{Name: "primary"},
		{Name: "work"},
	}

	err := repo.Update(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "update insert has many", user.Name)

	assert.Len(t, user.Addresses, 2)
	assert.NotEqual(t, 0, user.Addresses[0].ID)
	assert.NotEqual(t, 0, user.Addresses[1].ID)
	assert.Equal(t, user.ID, *user.Addresses[0].UserID)
	assert.Equal(t, user.ID, *user.Addresses[1].UserID)
	assert.Equal(t, "primary", user.Addresses[0].Name)
	assert.Equal(t, "work", user.Addresses[1].Name)

	assert.Equal(t, 2, repo.MustCount("addresses", where.Eq("user_id", user.ID)))
}

func UpdateHasOneInsert(t *testing.T, repo grimoire.Repo) {
	var (
		user = User{
			Name: "update init",
		}
	)

	repo.MustInsert(&user)

	user.Name = "update insert has one"
	user.PrimaryAddress = &Address{Name: "primary"}

	err := repo.Update(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "update insert has one", user.Name)

	assert.NotEqual(t, 0, user.PrimaryAddress.ID)
	assert.Equal(t, user.ID, *user.PrimaryAddress.UserID)
	assert.Equal(t, "primary", user.PrimaryAddress.Name)

	assert.Equal(t, 1, repo.MustCount("addresses", where.Eq("user_id", user.ID)))
}

func UpdateHasOneUpdate(t *testing.T, repo grimoire.Repo) {
	var (
		user = User{
			Name:           "update init",
			PrimaryAddress: &Address{Name: "primary"},
		}
	)

	repo.MustInsert(&user)

	user.Name = "update update has one"
	user.PrimaryAddress.Name = "updated primary"

	err := repo.Update(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "update update has one", user.Name)

	assert.NotEqual(t, 0, user.PrimaryAddress.ID)
	assert.Equal(t, user.ID, *user.PrimaryAddress.UserID)
	assert.Equal(t, "updated primary", user.PrimaryAddress.Name)

	assert.Equal(t, 1, repo.MustCount("addresses", where.Eq("user_id", user.ID)))
}

func UpdateHasOneReplace(t *testing.T, repo grimoire.Repo) {
	var (
		user = User{
			Name:           "update init",
			PrimaryAddress: &Address{Name: "primary"},
		}
	)

	repo.MustInsert(&user)

	user.Name = "update replace has one"
	user.PrimaryAddress = &Address{Name: "replaced primary"}

	err := repo.Update(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "update replace has one", user.Name)

	assert.NotEqual(t, 0, user.PrimaryAddress.ID)
	assert.Equal(t, user.ID, *user.PrimaryAddress.UserID)
	assert.Equal(t, "replaced primary", user.PrimaryAddress.Name)

	assert.Equal(t, 1, repo.MustCount("addresses", where.Eq("user_id", user.ID)))
}

func UpdateBelongsToInsert(t *testing.T, repo grimoire.Repo) {
	var (
		address = Address{Name: "address init"}
	)

	repo.MustInsert(&address)

	address.Name = "update address belongs to"
	address.User = User{Name: "inserted user"}

	err := repo.Update(&address)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, address.ID)
	assert.Equal(t, "update address belongs to", address.Name)

	assert.NotEqual(t, 0, address.User.ID)
	assert.Equal(t, *address.UserID, address.User.ID)
	assert.Equal(t, "inserted user", address.User.Name)
}

func UpdateBelongsToUpdate(t *testing.T, repo grimoire.Repo) {
	var (
		address = Address{
			Name: "address init",
			User: User{Name: "user"},
		}
	)

	repo.MustInsert(&address)

	address.Name = "update address belongs to"
	address.User.Name = "updated user"

	err := repo.Update(&address)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, address.ID)
	assert.Equal(t, "update address belongs to", address.Name)

	assert.NotEqual(t, 0, address.User.ID)
	assert.Equal(t, *address.UserID, address.User.ID)
	assert.Equal(t, "updated user", address.User.Name)
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

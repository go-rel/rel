package specs

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

// Update tests specification for updating a record.
func Update(t *testing.T, repo rel.Repository) {
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
	err = repo.Find(&queried, where.Eq("id", user.ID))
	assert.Nil(t, err)
	assert.Equal(t, user, queried)
}

// UpdateHasManyInsert tests specification for updating a record and inserting has many association.
func UpdateHasManyInsert(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{
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

	repo.MustFind(&result, where.Eq("id", user.ID))
	repo.MustPreload(&result, "addresses")

	assert.Equal(t, result, user)
}

// UpdateHasManyUpdate tests specification for updating a record and updating has many association.
func UpdateHasManyUpdate(t *testing.T, repo rel.Repository) {
	var (
		user = User{
			Name: "update init",
			Addresses: []Address{
				{Name: "old address"},
			},
		}
		result User
	)

	repo.MustInsert(&user)
	assert.NotEqual(t, 0, user.Addresses[0].ID)

	user.Name = "update insert has many"
	user.Addresses[0].Name = "new address"

	assert.Nil(t, repo.Update(&user))
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "update insert has many", user.Name)

	assert.Len(t, user.Addresses, 1)
	assert.NotEqual(t, 0, user.Addresses[0].ID)
	assert.Equal(t, user.ID, *user.Addresses[0].UserID)
	assert.Equal(t, "new address", user.Addresses[0].Name)

	repo.MustFind(&result, where.Eq("id", user.ID))
	repo.MustPreload(&result, "addresses")

	assert.Equal(t, result, user)
}

// UpdateHasManyReplace tests specification for updating a record and replacing has many association.
func UpdateHasManyReplace(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{
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

	repo.MustFind(&result, where.Eq("id", user.ID))
	repo.MustPreload(&result, "addresses")

	assert.Equal(t, result, user)
}

// UpdateHasOneInsert tests specification for updating a record and inserting has many association.
func UpdateHasOneInsert(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{
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

	repo.MustFind(&result, where.Eq("id", user.ID))
	repo.MustPreload(&result, "primary_address")

	assert.Equal(t, result, user)
}

// UpdateHasOneUpdate tests specification for updating a record and updating has one association.
func UpdateHasOneUpdate(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{
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

	repo.MustFind(&result, where.Eq("id", user.ID))
	repo.MustPreload(&result, "primary_address")

	assert.Equal(t, result, user)
}

// UpdateHasOneReplace tests specification for updating a record and replacing has one association.
func UpdateHasOneReplace(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{
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

	repo.MustFind(&result, where.Eq("id", user.ID))
	repo.MustPreload(&result, "primary_address")

	assert.Equal(t, result, user)
}

// UpdateBelongsToInsert tests specification for updating a record and inserting belongs to association.
func UpdateBelongsToInsert(t *testing.T, repo rel.Repository) {
	var (
		result  Address
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

	repo.MustFind(&result, where.Eq("id", address.ID))
	repo.MustPreload(&result, "user")

	assert.Equal(t, result, address)
}

// UpdateBelongsToUpdate tests specification for updating a record and updating belongs to association.
func UpdateBelongsToUpdate(t *testing.T, repo rel.Repository) {
	var (
		result  Address
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

	repo.MustFind(&result, where.Eq("id", address.ID))
	repo.MustPreload(&result, "user")

	assert.Equal(t, result, address)
}

// UpdateAtomic tests increment and decerement operation when updating a record.
func UpdateAtomic(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{Name: "update", Age: 10}
	)

	repo.MustInsert(&user)

	assert.Nil(t, repo.Update(&user, rel.Inc("age")))
	assert.Equal(t, 11, user.Age)

	repo.MustFind(&result, where.Eq("id", user.ID))
	assert.Equal(t, result, user)

	assert.Nil(t, repo.Update(&user, rel.Dec("age")))
	assert.Equal(t, 10, user.Age)

	repo.MustFind(&result, where.Eq("id", user.ID))
	assert.Equal(t, result, user)
}

// Updates tests update specifications.
func Updates(t *testing.T, repo rel.Repository) {
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
		t.Run("Update", func(t *testing.T) {
			assert.Nil(t, repo.Update(record))

			switch v := record.(type) {
			case *User:
				var found User
				repo.MustFind(&found, where.Eq("id", v.ID))
				assert.Equal(t, found, *v)
			case *Address:
				var found Address
				repo.MustFind(&found, where.Eq("id", v.ID))
				assert.Equal(t, found, *v)
			}
		})
	}
}

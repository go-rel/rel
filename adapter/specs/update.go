package specs

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
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

	repo.MustInsert(ctx, &user)

	user.Name = "update"
	user.Gender = "male"
	user.Age = 23
	user.Note = &note

	err := repo.Update(ctx, &user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "update", user.Name)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, 23, user.Age)
	assert.Equal(t, &note, user.Note)

	// update unchanged
	assert.Nil(t, repo.Update(ctx, &user))

	var (
		queried User
	)

	waitForReplication()

	user.Addresses = nil
	err = repo.Find(ctx, &queried, where.Eq("id", user.ID))
	assert.Nil(t, err)
	assert.Equal(t, user, queried)
}

// UpdateNotFound tests specification for updating a not found record.
func UpdateNotFound(t *testing.T, repo rel.Repository) {
	var (
		user = User{
			ID:   0,
			Name: "update",
		}
	)

	// update unchanged
	assert.Equal(t, rel.NotFoundError{}, repo.Update(ctx, &user))
}

// UpdateHasManyInsert tests specification for updating a record and inserting has many association.
func UpdateHasManyInsert(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{
			Name: "update init",
		}
	)

	repo.MustInsert(ctx, &user)

	user.Name = "update insert has many"
	user.Addresses = []Address{
		{Name: "primary"},
		{Name: "work"},
	}

	err := repo.Update(ctx, &user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "update insert has many", user.Name)

	assert.Len(t, user.Addresses, 2)
	assert.NotZero(t, user.Addresses[0].ID)
	assert.NotZero(t, user.Addresses[1].ID)
	assert.Equal(t, user.ID, *user.Addresses[0].UserID)
	assert.Equal(t, user.ID, *user.Addresses[1].UserID)
	assert.Equal(t, "primary", user.Addresses[0].Name)
	assert.Equal(t, "work", user.Addresses[1].Name)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	repo.MustPreload(ctx, &result, "addresses")

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

	repo.MustInsert(ctx, &user)
	assert.NotZero(t, user.Addresses[0].ID)

	user.Name = "update insert has many"
	user.Addresses[0].Name = "new address"

	assert.Nil(t, repo.Update(ctx, &user))
	assert.NotZero(t, user.ID)
	assert.Equal(t, "update insert has many", user.Name)

	assert.Len(t, user.Addresses, 1)
	assert.NotZero(t, user.Addresses[0].ID)
	assert.Equal(t, user.ID, *user.Addresses[0].UserID)
	assert.Equal(t, "new address", user.Addresses[0].Name)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	repo.MustPreload(ctx, &result, "addresses")

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

	repo.MustInsert(ctx, &user)

	user.Name = "update insert has many"
	user.Addresses = []Address{
		{Name: "primary"},
		{Name: "work"},
	}

	err := repo.Update(ctx, &user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "update insert has many", user.Name)

	assert.Len(t, user.Addresses, 2)
	assert.NotZero(t, user.Addresses[0].ID)
	assert.NotZero(t, user.Addresses[1].ID)
	assert.Equal(t, user.ID, *user.Addresses[0].UserID)
	assert.Equal(t, user.ID, *user.Addresses[1].UserID)
	assert.Equal(t, "primary", user.Addresses[0].Name)
	assert.Equal(t, "work", user.Addresses[1].Name)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	repo.MustPreload(ctx, &result, "addresses")

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

	repo.MustInsert(ctx, &user)

	user.Name = "update insert has one"
	user.PrimaryAddress = &Address{Name: "primary"}

	err := repo.Update(ctx, &user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "update insert has one", user.Name)

	assert.NotZero(t, user.PrimaryAddress.ID)
	assert.Equal(t, user.ID, *user.PrimaryAddress.UserID)
	assert.Equal(t, "primary", user.PrimaryAddress.Name)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	repo.MustPreload(ctx, &result, "primary_address")

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

	repo.MustInsert(ctx, &user)

	user.Name = "update update has one"
	user.PrimaryAddress.Name = "updated primary"

	err := repo.Update(ctx, &user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "update update has one", user.Name)

	assert.NotZero(t, user.PrimaryAddress.ID)
	assert.Equal(t, user.ID, *user.PrimaryAddress.UserID)
	assert.Equal(t, "updated primary", user.PrimaryAddress.Name)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	repo.MustPreload(ctx, &result, "primary_address")

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

	repo.MustInsert(ctx, &user)

	user.Name = "update replace has one"
	user.PrimaryAddress = &Address{Name: "replaced primary"}

	err := repo.Update(ctx, &user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "update replace has one", user.Name)

	assert.NotZero(t, user.PrimaryAddress.ID)
	assert.Equal(t, user.ID, *user.PrimaryAddress.UserID)
	assert.Equal(t, "replaced primary", user.PrimaryAddress.Name)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	repo.MustPreload(ctx, &result, "primary_address")

	assert.Equal(t, result, user)
}

// UpdateBelongsToInsert tests specification for updating a record and inserting belongs to association.
func UpdateBelongsToInsert(t *testing.T, repo rel.Repository) {
	var (
		result  Address
		address = Address{Name: "address init"}
	)

	repo.MustInsert(ctx, &address)

	address.Name = "update address belongs to"
	address.User = User{Name: "inserted user"}

	err := repo.Update(ctx, &address)
	assert.Nil(t, err)
	assert.NotZero(t, address.ID)
	assert.Equal(t, "update address belongs to", address.Name)

	assert.NotZero(t, address.User.ID)
	assert.Equal(t, *address.UserID, address.User.ID)
	assert.Equal(t, "inserted user", address.User.Name)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", address.ID))
	repo.MustPreload(ctx, &result, "user")

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

	repo.MustInsert(ctx, &address)

	address.Name = "update address belongs to"
	address.User.Name = "updated user"

	err := repo.Update(ctx, &address)
	assert.Nil(t, err)
	assert.NotZero(t, address.ID)
	assert.Equal(t, "update address belongs to", address.Name)

	assert.NotZero(t, address.User.ID)
	assert.Equal(t, *address.UserID, address.User.ID)
	assert.Equal(t, "updated user", address.User.Name)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", address.ID))
	repo.MustPreload(ctx, &result, "user")

	assert.Equal(t, result, address)
}

// UpdateAtomic tests increment and decerement operation when updating a record.
func UpdateAtomic(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{Name: "update", Age: 10}
	)

	repo.MustInsert(ctx, &user)

	assert.Nil(t, repo.Update(ctx, &user, rel.Inc("age")))
	assert.Equal(t, 11, user.Age)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	assert.Equal(t, result, user)

	assert.Nil(t, repo.Update(ctx, &user, rel.Dec("age")))
	assert.Equal(t, 10, user.Age)

	waitForReplication()

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	assert.Equal(t, result, user)
}

// Updates tests update specifications.
func Updates(t *testing.T, repo rel.Repository) {
	var (
		note    = "note"
		user    = User{Name: "update"}
		address = Address{Name: "update"}
	)

	repo.MustInsert(ctx, &user)
	repo.MustInsert(ctx, &address)

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
			assert.Nil(t, repo.Update(ctx, record))
		})
	}
}

// UpdateAny tests update all specifications.
func UpdateAny(t *testing.T, repo rel.Repository) {
	repo.MustInsert(ctx, &User{Name: "update", Age: 100})
	repo.MustInsert(ctx, &User{Name: "update", Age: 100})
	repo.MustInsert(ctx, &User{Name: "other update", Age: 110})

	tests := []rel.Query{
		rel.From("users").Where(where.Eq("name", "update")),
		rel.From("users").Where(where.Eq("name", "other update"), where.Gt("age", 100)),
	}

	for _, query := range tests {
		t.Run("UpdateAny", func(t *testing.T) {
			var (
				result []User
				name   = "all updated"
			)

			updatedCount, err := repo.UpdateAny(ctx, query, rel.Set("name", name))
			assert.Nil(t, err)
			assert.NotZero(t, updatedCount)

			waitForReplication()

			assert.Nil(t, repo.FindAll(ctx, &result, query))
			assert.Zero(t, len(result))
			for i := range result {
				assert.Equal(t, name, result[i].Name)
			}
		})
	}
}

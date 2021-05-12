package specs

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

// Delete tests delete specifications.
func Delete(t *testing.T, repo rel.Repository) {
	var (
		address = Address{
			Name: "address",
			User: User{Name: "user", Age: 100},
		}
	)

	repo.MustInsert(ctx, &address)
	assert.NotEqual(t, 0, address.ID)
	assert.NotEqual(t, 0, address.User.ID)

	assert.Nil(t, repo.Delete(ctx, &address))
	assert.Equal(t, rel.NotFoundError{}, repo.Find(ctx, &Address{}, where.Eq("id", address.ID)))

	// not deleted because cascade disabled
	assert.Nil(t, repo.Find(ctx, &User{}, where.Eq("id", address.User.ID)))
}

// DeleteBelongsTo tests delete specifications.
func DeleteBelongsTo(t *testing.T, repo rel.Repository) {
	var (
		address = Address{
			Name: "address",
			User: User{Name: "user", Age: 100},
		}
	)

	repo.MustInsert(ctx, &address)
	assert.NotEqual(t, 0, address.ID)
	assert.NotEqual(t, 0, address.User.ID)

	assert.Nil(t, repo.Delete(ctx, &address, rel.Cascade(true)))
	assert.Equal(t, rel.NotFoundError{}, repo.Find(ctx, &Address{}, where.Eq("id", address.ID)))
	assert.Equal(t, rel.NotFoundError{}, repo.Find(ctx, &User{}, where.Eq("id", address.User.ID)))
}

// DeleteHasOne tests delete specifications.
func DeleteHasOne(t *testing.T, repo rel.Repository) {
	var (
		user = User{
			Name:           "user",
			Age:            100,
			PrimaryAddress: &Address{Name: "primary address"},
		}
	)

	repo.MustInsert(ctx, &user)
	assert.NotEqual(t, 0, user.ID)
	assert.NotEqual(t, 0, user.PrimaryAddress.ID)

	assert.Nil(t, repo.Delete(ctx, &user, rel.Cascade(true)))
	assert.Equal(t, rel.NotFoundError{}, repo.Find(ctx, &User{}, where.Eq("id", user.ID)))
	assert.Equal(t, rel.NotFoundError{}, repo.Find(ctx, &Address{}, where.Eq("id", user.PrimaryAddress.ID)))
}

// DeleteHasMany tests delete specifications.
func DeleteHasMany(t *testing.T, repo rel.Repository) {
	var (
		user = User{
			Name: "user",
			Age:  100,
			Addresses: []Address{
				{Name: "address 1"},
				{Name: "address 2"},
			},
		}
	)

	repo.MustInsert(ctx, &user)
	assert.NotEqual(t, 0, user.ID)
	assert.NotEqual(t, 0, user.Addresses[0].ID)
	assert.NotEqual(t, 0, user.Addresses[1].ID)

	assert.Nil(t, repo.Delete(ctx, &user, rel.Cascade(true)))
	assert.Equal(t, rel.NotFoundError{}, repo.Find(ctx, &User{}, where.Eq("id", user.ID)))
	assert.Equal(t, rel.NotFoundError{}, repo.Find(ctx, &Address{}, where.Eq("id", user.Addresses[0].ID)))
	assert.Equal(t, rel.NotFoundError{}, repo.Find(ctx, &Address{}, where.Eq("id", user.Addresses[1].ID)))
}

// DeleteAll tests delete all specifications.
func DeleteAll(t *testing.T, repo rel.Repository) {
	repo.MustInsert(ctx, &User{Name: "delete", Age: 100})
	repo.MustInsert(ctx, &User{Name: "delete", Age: 100})
	repo.MustInsert(ctx, &User{Name: "other delete", Age: 110})

	tests := []rel.Query{
		rel.From("users").Where(where.Eq("name", "delete")),
		rel.From("users").Where(where.Eq("name", "other delete"), where.Gt("age", 100)),
	}

	for _, query := range tests {
		var result []User
		t.Run("DeleteAll", func(t *testing.T) {
			assert.Nil(t, repo.FindAll(ctx, &result, query))
			assert.NotEqual(t, 0, len(result))

			deletedCount, err := repo.DeleteAll(ctx, query)
			assert.Nil(t, err)
			assert.NotZero(t, deletedCount)

			assert.Nil(t, repo.FindAll(ctx, &result, query))
			assert.Zero(t, len(result))
		})
	}
}

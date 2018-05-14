package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
	"github.com/stretchr/testify/assert"
)

// Preload tests query specifications for preloading.
func Preload(t *testing.T, repo grimoire.Repo) {
	// preparte tests data
	user := User{Name: "preload", Gender: "male", Age: 10}
	repo.From(users).MustSave(&user)

	userAddresses := []Address{
		{Address: "preload1", UserID: &user.ID},
		{Address: "preload2", UserID: &user.ID},
		{Address: "preload3", UserID: &user.ID},
	}

	repo.From(addresses).MustSave(&userAddresses[0])
	repo.From(addresses).MustSave(&userAddresses[1])
	repo.From(addresses).MustSave(&userAddresses[2])

	assert.Nil(t, user.Addresses)

	t.Run("Preload Addresses", func(t *testing.T) {
		repo.From(addresses).Preload(&user, "Addresses")
		assert.Equal(t, userAddresses, user.Addresses)
	})

	t.Run("Preload Addresses with query", func(t *testing.T) {
		repo.From(addresses).Where(c.Eq(address, "preload1")).Preload(&user, "Addresses")
		assert.Equal(t, 1, len(user.Addresses))
		assert.Equal(t, userAddresses[0], user.Addresses[0])
	})

	user.Addresses = nil

	t.Run("Preload User", func(t *testing.T) {
		repo.From(users).Preload(&userAddresses[0], "User")
		assert.Equal(t, user, userAddresses[0].User)
	})

	t.Run("Preload User slice", func(t *testing.T) {
		repo.From(users).Preload(&userAddresses, "User")
		assert.Equal(t, user, userAddresses[0].User)
		assert.Equal(t, user, userAddresses[1].User)
		assert.Equal(t, user, userAddresses[2].User)
	})
}

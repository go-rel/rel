package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

// Preload tests query specifications for preloading.
func Preload(t *testing.T, repo grimoire.Repo) {
	// preparte tests data
	var (
		user = User{
			Name:   "preload",
			Gender: "male",
			Age:    25,
			Addresses: []Address{
				{Address: "preload1"},
				{Address: "preload2"},
				{Address: "preload3"},
			},
		}
	)

	repo.MustInsert(&user)

	// t.Run("Preload Addresses", func(t *testing.T) {
	// 	var (
	// 		emptyUser = User{ID: user.ID}
	// 	)

	// 	err := repo.Preload(&emptyUser, "addresses")
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, len(emptyUser.Addresses), len(user.Addresses))
	// })

	t.Run("Preload Addresses with query", func(t *testing.T) {
		var (
			emptyUser = User{ID: user.ID}
		)

		repo.Preload(&emptyUser, "addresses", where.Eq("address", "preload1"))
		assert.Equal(t, 1, len(emptyUser.Addresses))
		assert.Equal(t, user.Addresses[0].Address, emptyUser.Addresses[0].Address)
	})

	// unload
	user.Addresses = nil

	t.Run("Preload User", func(t *testing.T) {
		var (
			emptyAddress = Address{UserID: &user.ID}
		)

		repo.Preload(&emptyAddress, "user")
		assert.Equal(t, user, emptyAddress.User)
	})

	t.Run("Preload User slice", func(t *testing.T) {
		var (
			emptyAddresses = []Address{
				{UserID: &user.ID},
				{UserID: &user.ID},
			}
		)

		repo.Preload(&emptyAddresses, "user")
		assert.Len(t, emptyAddresses, 2)
		assert.Equal(t, user, emptyAddresses[0].User)
		assert.Equal(t, user, emptyAddresses[0].User)
	})
}

package specs

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func createPreloadUser(repo rel.Repository) User {
	var (
		user = User{
			Name:   "preload",
			Gender: "male",
			Age:    25,
			Addresses: []Address{
				{Name: "primary"},
				{Name: "home"},
				{Name: "work"},
			},
		}
	)

	repo.MustInsert(&user)

	return user
}

func PreloadHasMany(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = createPreloadUser(repo)
	)

	err := repo.One(&result, where.Eq("id", user.ID))
	assert.Nil(t, err)

	err = repo.Preload(&result, "addresses")
	assert.Nil(t, err)
	assert.Equal(t, user, result)
}

func PreloadHasManyWithQuery(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = createPreloadUser(repo)
	)

	err := repo.One(&result, where.Eq("id", user.ID))
	assert.Nil(t, err)

	err = repo.Preload(&result, "addresses", where.Eq("name", "primary"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result.Addresses))
	assert.Equal(t, user.Addresses[0], result.Addresses[0])
}

func PreloadHasManySlice(t *testing.T, repo rel.Repository) {
	var (
		result []User
		users  = []User{
			createPreloadUser(repo),
			createPreloadUser(repo),
		}
	)

	err := repo.All(&result, where.In("id", users[0].ID, users[1].ID))
	assert.Nil(t, err)

	err = repo.Preload(&result, "addresses")
	assert.Nil(t, err)
	assert.Equal(t, users, result)
}

func PreloadHasOne(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = createPreloadUser(repo)
	)

	err := repo.One(&result, where.Eq("id", user.ID))
	assert.Nil(t, err)

	err = repo.Preload(&result, "primary_address")
	assert.Nil(t, err)
	assert.NotNil(t, result.PrimaryAddress)
}

func PreloadHasOneWithQuery(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = createPreloadUser(repo)
	)

	err := repo.One(&result, where.Eq("id", user.ID))
	assert.Nil(t, err)

	err = repo.Preload(&result, "primary_address", where.Eq("name", "primary"))
	assert.Nil(t, err)
	assert.Equal(t, user.Addresses[0], *result.PrimaryAddress)
}

func PreloadHasOneSlice(t *testing.T, repo rel.Repository) {
	var (
		result []User
		users  = []User{
			createPreloadUser(repo),
			createPreloadUser(repo),
		}
	)

	err := repo.All(&result, where.In("id", users[0].ID, users[1].ID))
	assert.Nil(t, err)

	err = repo.Preload(&result, "primary_address")
	assert.Nil(t, err)
	assert.NotNil(t, result[0].PrimaryAddress)
	assert.NotNil(t, result[1].PrimaryAddress)
}

func PreloadBelongsTo(t *testing.T, repo rel.Repository) {
	var (
		result Address
		user   = createPreloadUser(repo)
	)

	err := repo.One(&result, where.Eq("id", user.Addresses[0].ID))
	assert.Nil(t, err)

	user.Addresses = nil

	err = repo.Preload(&result, "user")
	assert.Nil(t, err)
	assert.Equal(t, user, result.User)
}

func PreloadBelongsToWithQuery(t *testing.T, repo rel.Repository) {
	var (
		result Address
		user   = createPreloadUser(repo)
	)

	err := repo.One(&result, where.Eq("id", user.Addresses[0].ID))
	assert.Nil(t, err)

	user.Addresses = nil

	err = repo.Preload(&result, "user", where.Eq("name", "not exists"))
	assert.Nil(t, err)
	assert.Zero(t, result.User)
}

func PreloadBelongsToSlice(t *testing.T, repo rel.Repository) {
	var (
		user      = createPreloadUser(repo)
		result    = user.Addresses
		resultLen = len(result)
	)

	user.Addresses = nil

	err := repo.Preload(&result, "user")
	assert.Nil(t, err)
	assert.Len(t, result, resultLen)

	for i := range result {
		assert.Equal(t, user, result[i].User)
	}
}

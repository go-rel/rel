package specs

import (
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
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

	repo.MustInsert(ctx, &user)

	return user
}

// PreloadHasMany tests specification for preloading has many association.
func PreloadHasMany(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = createPreloadUser(repo)
	)

	waitForReplication()

	err := repo.Find(ctx, &result, where.Eq("id", user.ID))
	assert.Nil(t, err)

	err = repo.Preload(ctx, &result, "addresses")
	assert.Nil(t, err)
	assert.Equal(t, user, result)
}

// PreloadHasManyWithQuery tests specification for preloading has many association.
func PreloadHasManyWithQuery(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = createPreloadUser(repo)
	)

	waitForReplication()

	err := repo.Find(ctx, &result, where.Eq("id", user.ID))
	assert.Nil(t, err)

	err = repo.Preload(ctx, &result, "addresses", where.Eq("name", "primary"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result.Addresses))
	assert.Equal(t, user.Addresses[0], result.Addresses[0])
}

// PreloadHasManySlice tests specification for preloading has many association from multiple records.
func PreloadHasManySlice(t *testing.T, repo rel.Repository) {
	var (
		result []User
		users  = []User{
			createPreloadUser(repo),
			createPreloadUser(repo),
		}
	)

	waitForReplication()

	err := repo.FindAll(ctx, &result, where.In("id", users[0].ID, users[1].ID))
	assert.Nil(t, err)

	err = repo.Preload(ctx, &result, "addresses")
	assert.Nil(t, err)
	assert.Equal(t, users, result)
}

// PreloadHasOne tests specification for preloading has one association.
func PreloadHasOne(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = createPreloadUser(repo)
	)

	waitForReplication()

	err := repo.Find(ctx, &result, where.Eq("id", user.ID))
	assert.Nil(t, err)

	err = repo.Preload(ctx, &result, "primary_address")
	assert.Nil(t, err)
	assert.NotNil(t, result.PrimaryAddress)
}

// PreloadHasOneWithQuery tests specification for preloading has one association.
func PreloadHasOneWithQuery(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = createPreloadUser(repo)
	)

	waitForReplication()

	err := repo.Find(ctx, &result, where.Eq("id", user.ID))
	assert.Nil(t, err)

	err = repo.Preload(ctx, &result, "primary_address", where.Eq("name", "primary"))
	assert.Nil(t, err)
	assert.Equal(t, user.Addresses[0], *result.PrimaryAddress)
}

// PreloadHasOneSlice tests specification for preloading has one association from multiple records.
func PreloadHasOneSlice(t *testing.T, repo rel.Repository) {
	var (
		result []User
		users  = []User{
			createPreloadUser(repo),
			createPreloadUser(repo),
		}
	)

	waitForReplication()

	err := repo.FindAll(ctx, &result, where.In("id", users[0].ID, users[1].ID))
	assert.Nil(t, err)

	err = repo.Preload(ctx, &result, "primary_address")
	assert.Nil(t, err)
	assert.NotNil(t, result[0].PrimaryAddress)
	assert.NotNil(t, result[1].PrimaryAddress)
}

// PreloadBelongsTo tests specification for preloading belongs to association.
func PreloadBelongsTo(t *testing.T, repo rel.Repository) {
	var (
		result Address
		user   = createPreloadUser(repo)
	)

	waitForReplication()

	err := repo.Find(ctx, &result, where.Eq("id", user.Addresses[0].ID))
	assert.Nil(t, err)

	user.Addresses = nil

	err = repo.Preload(ctx, &result, "user")
	assert.Nil(t, err)
	assert.Equal(t, user, result.User)
}

// PreloadBelongsToWithQuery tests specification for preloading belongs to association.
func PreloadBelongsToWithQuery(t *testing.T, repo rel.Repository) {
	var (
		result Address
		user   = createPreloadUser(repo)
	)

	waitForReplication()

	err := repo.Find(ctx, &result, where.Eq("id", user.Addresses[0].ID))
	assert.Nil(t, err)

	user.Addresses = nil

	err = repo.Preload(ctx, &result, "user", where.Eq("name", "not exists"))
	assert.Nil(t, err)
	assert.Zero(t, result.User)
}

// PreloadBelongsToSlice tests specification for preloading belongs to association from multiple records.
func PreloadBelongsToSlice(t *testing.T, repo rel.Repository) {
	var (
		user      = createPreloadUser(repo)
		result    = user.Addresses
		resultLen = len(result)
	)

	waitForReplication()

	user.Addresses = nil

	err := repo.Preload(ctx, &result, "user")
	assert.Nil(t, err)
	assert.Len(t, result, resultLen)

	for i := range result {
		assert.Equal(t, user, result[i].User)
	}
}

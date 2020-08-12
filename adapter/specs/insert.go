package specs

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

// Insert tests specification for database insertion.
func Insert(t *testing.T, repo rel.Repository) {
	var (
		note = "swordsman"
		user = User{
			Name:   "insert",
			Gender: "male",
			Age:    23,
			Note:   &note,
		}
	)

	err := repo.Insert(ctx, &user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "insert", user.Name)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, 23, user.Age)
	assert.Equal(t, &note, user.Note)

	var (
		queried User
	)

	user.Addresses = nil
	err = repo.Find(ctx, &queried, where.Eq("id", user.ID))
	assert.Nil(t, err)
	assert.Equal(t, user, queried)
}

// InsertHasMany tests specification insertion with has many association.
func InsertHasMany(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{
			Name:   "insert has many",
			Gender: "male",
			Age:    23,
			Addresses: []Address{
				{Name: "primary"},
				{Name: "work"},
			},
		}
	)

	err := repo.Insert(ctx, &user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "insert has many", user.Name)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, 23, user.Age)

	assert.Len(t, user.Addresses, 2)
	assert.NotEqual(t, 0, user.Addresses[0].ID)
	assert.NotEqual(t, 0, user.Addresses[1].ID)
	assert.Equal(t, user.ID, *user.Addresses[0].UserID)
	assert.Equal(t, user.ID, *user.Addresses[1].UserID)
	assert.Equal(t, "primary", user.Addresses[0].Name)
	assert.Equal(t, "work", user.Addresses[1].Name)

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	repo.MustPreload(ctx, &result, "addresses")

	assert.Equal(t, result, user)
}

// InsertHasOne tests specification for insertion with has one association.
func InsertHasOne(t *testing.T, repo rel.Repository) {
	var (
		result User
		user   = User{
			Name:           "insert has one",
			Gender:         "male",
			Age:            23,
			PrimaryAddress: &Address{Name: "primary"},
		}
	)

	err := repo.Insert(ctx, &user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "insert has one", user.Name)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, 23, user.Age)

	assert.NotEqual(t, 0, user.PrimaryAddress.ID)
	assert.Equal(t, user.ID, *user.PrimaryAddress.UserID)
	assert.Equal(t, "primary", user.PrimaryAddress.Name)

	repo.MustFind(ctx, &result, where.Eq("id", user.ID))
	repo.MustPreload(ctx, &result, "primary_address")

	assert.Equal(t, result, user)
}

// InsertBelongsTo tests specification for insertion with belongs to association.
func InsertBelongsTo(t *testing.T, repo rel.Repository) {
	var (
		result  Address
		address = Address{
			Name: "insert belongs to",
			User: User{
				Name:   "zoro",
				Gender: "male",
				Age:    23,
			},
		}
	)

	err := repo.Insert(ctx, &address)
	assert.Nil(t, err)

	assert.NotEqual(t, 0, address.ID)
	assert.Equal(t, address.User.ID, *address.UserID)
	assert.Equal(t, "insert belongs to", address.Name)

	assert.NotEqual(t, 0, address.User.ID)
	assert.Equal(t, "zoro", address.User.Name)
	assert.Equal(t, "male", address.User.Gender)
	assert.Equal(t, 23, address.User.Age)

	repo.MustFind(ctx, &result, where.Eq("id", address.ID))
	repo.MustPreload(ctx, &result, "user")

	assert.Equal(t, result, address)
}

// Inserts tests insert specifications.
func Inserts(t *testing.T, repo rel.Repository) {
	var (
		user User
		note = "note"
	)

	repo.MustInsert(ctx, &user)

	tests := []interface{}{
		&User{},
		&User{Name: "insert", Age: 100},
		&User{Name: "insert", Age: 100, Note: &note},
		&User{Note: &note},
		&User{ID: 123, Name: "insert", Age: 100, Note: &note},
		&Address{},
		&Address{Name: "work"},
		&Address{UserID: &user.ID},
		&Address{Name: "work", UserID: &user.ID},
		&Address{ID: 123, Name: "work", UserID: &user.ID},
		&Composite{Primary1: 1, Primary2: 2, Data: "data-1-2"},
	}

	for _, record := range tests {
		t.Run("Insert", func(t *testing.T) {
			assert.Nil(t, repo.Insert(ctx, record))
			assertRecord(t, repo, record)
		})
	}
}

func assertRecord(t *testing.T, repo rel.Repository, record interface{}) {
	switch v := record.(type) {
	case *User:
		var found User
		repo.MustFind(ctx, &found, where.Eq("id", v.ID))
		assert.Equal(t, found, *v)
	case *Address:
		var found Address
		repo.MustFind(ctx, &found, where.Eq("id", v.ID))
		assert.Equal(t, found, *v)
	case *Composite:
		var found Composite
		repo.MustFind(ctx, &found, where.Eq("primary1", v.Primary1).AndEq("primary2", v.Primary2))
		assert.Equal(t, found, *v)
	}
}

// InsertAll tests insert multiple specifications.
func InsertAll(t *testing.T, repo rel.Repository) {
	var (
		user User
		note = "note"
	)

	repo.MustInsert(ctx, &user)

	tests := []interface{}{
		&[]User{{}},
		&[]User{{Name: "insert", Age: 100}},
		&[]User{{Name: "insert", Age: 100, Note: &note}},
		&[]User{{Note: &note}},
		&[]User{{Name: "insert", Age: 100}, {Name: "insert too"}},
		&[]Address{{}},
		&[]Address{{Name: "work"}},
		&[]Address{{UserID: &user.ID}},
		&[]Address{{Name: "work", UserID: &user.ID}},
		&[]Address{{Name: "work"}, {Name: "home"}},
	}

	for _, record := range tests {
		t.Run("InsertAll", func(t *testing.T) {
			assert.Nil(t, repo.InsertAll(ctx, record))

			switch v := record.(type) {
			case *[]User:
				var (
					found []User
					ids   = make([]int, len(*v))
				)

				for i := range *v {
					ids[i] = int((*v)[i].ID)
				}

				repo.MustFindAll(ctx, &found, where.InInt("id", ids))
				assert.Equal(t, found, *v)
			case *[]Address:
				var (
					found []Address
					ids   = make([]int, len(*v))
				)

				for i := range *v {
					ids[i] = int((*v)[i].ID)
				}

				repo.MustFindAll(ctx, &found, where.InInt("id", ids))
				assert.Equal(t, found, *v)
			}
		})
	}
}

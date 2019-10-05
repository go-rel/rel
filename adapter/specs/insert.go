package specs

import (
	"reflect"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func Insert(t *testing.T, repo rel.Repo) {
	var (
		note = "swordsman"
		user = User{
			Name:   "insert",
			Gender: "male",
			Age:    23,
			Note:   &note,
		}
	)

	err := repo.Insert(&user)
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
	err = repo.One(&queried, where.Eq("id", user.ID))
	assert.Nil(t, err)
	assert.Equal(t, user, queried)
}

func InsertHasMany(t *testing.T, repo rel.Repo) {
	var (
		user = User{
			Name:   "insert has many",
			Gender: "male",
			Age:    23,
			Addresses: []Address{
				{Name: "primary"},
				{Name: "work"},
			},
		}
	)

	err := repo.Insert(&user)
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
}

func InsertHasOne(t *testing.T, repo rel.Repo) {
	var (
		user = User{
			Name:           "insert has one",
			Gender:         "male",
			Age:            23,
			PrimaryAddress: &Address{Name: "primary"},
		}
	)

	err := repo.Insert(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, "insert has one", user.Name)
	assert.Equal(t, "male", user.Gender)
	assert.Equal(t, 23, user.Age)

	assert.NotEqual(t, 0, user.PrimaryAddress.ID)
	assert.Equal(t, user.ID, *user.PrimaryAddress.UserID)
	assert.Equal(t, "primary", user.PrimaryAddress.Name)
}

func InsertBelongsTo(t *testing.T, repo rel.Repo) {
	var (
		address = Address{
			Name: "insert belongs to",
			User: User{
				Name:   "zoro",
				Gender: "male",
				Age:    23,
			},
		}
	)

	err := repo.Insert(&address)
	assert.Nil(t, err)

	assert.NotEqual(t, 0, address.ID)
	assert.Equal(t, address.User.ID, *address.UserID)
	assert.Equal(t, "insert belongs to", address.Name)

	assert.NotEqual(t, 0, address.User.ID)
	assert.Equal(t, "zoro", address.User.Name)
	assert.Equal(t, "male", address.User.Gender)
	assert.Equal(t, 23, address.User.Age)
}

// Inserts tests insert specifications.
func Inserts(t *testing.T, repo rel.Repo) {
	var (
		user User
		note = "note"
	)

	repo.MustInsert(&user)

	tests := []interface{}{
		&User{},
		&User{Name: "insert", Age: 100},
		&User{Name: "insert", Age: 100, Note: &note},
		&User{Note: &note},
		&Address{},
		&Address{Name: "work"},
		&Address{UserID: &user.ID},
		&Address{Name: "work", UserID: &user.ID},
	}

	for _, record := range tests {
		var (
			changes      = rel.BuildChanges(rel.NewStructset(record))
			statement, _ = builder.Insert("collection", changes)
		)

		t.Run("Insert|"+statement, func(t *testing.T) {
			assert.Nil(t, repo.Insert(record))
		})
	}
}

// InsertAll tests insert multiple specifications.
func InsertAll(t *testing.T, repo rel.Repo) {
	var (
		user User
		note = "note"
	)

	repo.MustInsert(&user)

	tests := []interface{}{
		// &[]User{{}},
		&[]User{{Name: "insert", Age: 100}},
		&[]User{{Name: "insert", Age: 100, Note: &note}},
		&[]User{{Note: &note}},
		// &[]Address{{}},
		&[]Address{{Name: "work"}},
		&[]Address{{UserID: &user.ID}},
		&[]Address{{Name: "work", UserID: &user.ID}},
	}

	for _, record := range tests {
		// var (
		// 	changes      = rel.BuildChanges(rel.Struct(record))
		// 	statement, _ = builder.Insert("collection", changes)
		// )

		t.Run("InsertAll", func(t *testing.T) {
			// multiple insert
			assert.Nil(t, repo.InsertAll(record))
			assert.Equal(t, 1, reflect.ValueOf(record).Elem().Len())
		})
	}
}

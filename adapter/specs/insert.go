package specs

import (
	"reflect"
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func Insert(t *testing.T, repo grimoire.Repo) {
	var (
		name   = "zoro"
		gender = "male"
		age    = 23
		note   = "swordsman"
		user   = User{
			Name:   name,
			Gender: gender,
			Age:    age,
			Note:   &note,
		}
	)

	err := repo.Insert(&user)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, user.ID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, gender, user.Gender)
	assert.Equal(t, age, user.Age)
	assert.Equal(t, &note, user.Note)

	var (
		queried User
	)

	user.Addresses = nil
	err = repo.One(&queried, where.Eq("id", user.ID))
	assert.Nil(t, err)
	assert.Equal(t, user, queried)
}

// Insert tests insert specifications.
// TODO: insert with assocs
func Inserts(t *testing.T, repo grimoire.Repo) {
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
			changes      = grimoire.BuildChanges(grimoire.Struct(record))
			statement, _ = builder.Insert("collection", changes)
		)

		t.Run("Insert|"+statement, func(t *testing.T) {
			assert.Nil(t, repo.Insert(record))
		})
	}
}

// InsertAll tests insert multiple specifications.
func InsertAll(t *testing.T, repo grimoire.Repo) {
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
		// 	changes      = grimoire.BuildChanges(grimoire.Struct(record))
		// 	statement, _ = builder.Insert("collection", changes)
		// )

		t.Run("InsertAll", func(t *testing.T) {
			// multiple insert
			assert.Nil(t, repo.InsertAll(record))
			assert.Equal(t, 1, reflect.ValueOf(record).Elem().Len())
		})
	}
}

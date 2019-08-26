package specs

import (
	"reflect"
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

func InsertBasic(t *testing.T, repo grimoire.Repo) {
	var (
		name       = "zoro"
		gender     = "male"
		age        = 23
		note       = "swordsman"
		insertUser = User{
			Name:   name,
			Gender: gender,
			Age:    age,
			Note:   &note,
		}
	)

	err := repo.Insert(&insertUser)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, insertUser.ID)
	assert.Equal(t, name, insertUser.Name)
	assert.Equal(t, gender, insertUser.Gender)
	assert.Equal(t, age, insertUser.Age)
	assert.Equal(t, &note, insertUser.Note)

	var (
		queryUser User
	)

	insertUser.Addresses = nil
	err = repo.One(&queryUser, where.Eq("id", insertUser.ID))
	assert.Nil(t, err)
	assert.Equal(t, insertUser, queryUser)
}

// Insert tests insert specifications.
// TODO: insert with assocs
func Insert(t *testing.T, repo grimoire.Repo) {
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
		&Address{Address: "address"},
		&Address{UserID: &user.ID},
		&Address{Address: "address", UserID: &user.ID},
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

// InsertExplicit tests insert specifications.
func InsertExplicit(t *testing.T, repo grimoire.Repo) {
	var (
		user User
	)

	repo.MustInsert(&user)

	tests := []struct {
		record  interface{}
		changer grimoire.Changer
	}{
		{&User{}, grimoire.Map{}},
		{&User{}, grimoire.Map{"name": "insert", "age": 100}},
		{&User{}, grimoire.Map{"name": "insert", "age": 100, "note": "note"}},
		{&User{}, grimoire.Map{"note": "note"}},
		{&Address{}, grimoire.Map{}},
		{&Address{}, grimoire.Map{"address": "address"}},
		{&Address{}, grimoire.Map{"user_id": user.ID}},
		{&Address{}, grimoire.Map{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		var (
			changes      = grimoire.BuildChanges(test.changer)
			statement, _ = builder.Insert("collection", changes)
		)

		t.Run("Insert|"+statement, func(t *testing.T) {
			assert.Nil(t, repo.Insert(test.record, test.changer))
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
		&[]Address{{Address: "address"}},
		&[]Address{{UserID: &user.ID}},
		&[]Address{{Address: "address", UserID: &user.ID}},
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

// InsertAllExplicit tests insert multiple specifications.
func InsertAllExplicit(t *testing.T, repo grimoire.Repo) {
	var (
		user User
	)

	repo.MustInsert(&user)

	tests := []struct {
		record  interface{}
		changer grimoire.Changer
	}{
		// {&[]User{}, grimoire.Map{}},
		{&[]User{}, grimoire.Map{"name": "insert", "age": 100}},
		{&[]User{}, grimoire.Map{"name": "insert", "age": 100, "note": "note"}},
		{&[]User{}, grimoire.Map{"note": "note"}},
		// {&[]Address{}, grimoire.Map{}},
		{&[]Address{}, grimoire.Map{"address": "address"}},
		{&[]Address{}, grimoire.Map{"user_id": user.ID}},
		{&[]Address{}, grimoire.Map{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		var (
			changes      = grimoire.BuildChanges(test.changer)
			statement, _ = builder.Insert("collection", changes)
		)
		println(statement)

		t.Run("InsertAllExplicit|"+statement, func(t *testing.T) {
			// multiple insert
			assert.Nil(t, repo.InsertAll(test.record, changes, changes, changes))
			assert.Equal(t, 3, reflect.ValueOf(test.record).Elem().Len())
		})
	}
}

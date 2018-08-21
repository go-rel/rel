package specs

import (
	"reflect"
	"testing"
	"time"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

// Insert tests insert specifications.
func Insert(t *testing.T, repo grimoire.Repo) {
	user := User{}
	repo.From(users).MustSave(&user)

	tests := []struct {
		query  grimoire.Query
		record interface{}
		input  params.Params
	}{
		{repo.From(users), &User{}, params.Map{}},
		{repo.From(users), &User{}, params.Map{"name": "insert", "age": 100}},
		{repo.From(users), &User{}, params.Map{"name": "insert", "age": 100, "note": "note"}},
		{repo.From(users), &User{}, params.Map{"note": "note"}},
		{repo.From(addresses), &Address{}, params.Map{}},
		{repo.From(addresses), &Address{}, params.Map{"address": "address"}},
		{repo.From(addresses), &Address{}, params.Map{"user_id": user.ID}},
		{repo.From(addresses), &Address{}, params.Map{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		ch := changeset.Cast(test.record, test.input, []string{"name", "age", "note", "address", "user_id"})
		statement, _ := builder.Insert(test.query.Collection, ch.Changes())

		t.Run("Insert|"+statement, func(t *testing.T) {
			assert.Nil(t, ch.Error())

			assert.Nil(t, test.query.Insert(nil, ch))
			assert.Nil(t, test.query.Insert(test.record, ch))

			// multiple insert
			assert.Nil(t, test.query.Insert(nil, ch, ch, ch))
		})
	}
}

// InsertAll tests insert multiple specifications.
func InsertAll(t *testing.T, repo grimoire.Repo) {
	user := User{}
	repo.From(users).MustSave(&user)

	tests := []struct {
		query  grimoire.Query
		schema interface{}
		record interface{}
		params params.Params
	}{
		{repo.From(users), User{}, &[]User{}, params.Map{}},
		{repo.From(users), User{}, &[]User{}, params.Map{"name": "insert", "age": 100}},
		{repo.From(users), User{}, &[]User{}, params.Map{"name": "insert", "age": 100, "note": "note"}},
		{repo.From(users), User{}, &[]User{}, params.Map{"note": "note"}},
		{repo.From(addresses), &Address{}, &[]Address{}, params.Map{}},
		{repo.From(addresses), &Address{}, &[]Address{}, params.Map{"address": "address"}},
		{repo.From(addresses), &Address{}, &[]Address{}, params.Map{"user_id": user.ID}},
		{repo.From(addresses), &Address{}, &[]Address{}, params.Map{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		ch := changeset.Cast(test.schema, test.params, []string{"name", "age", "note", "address", "user_id"})
		statement, _ := builder.Insert(test.query.Collection, ch.Changes())

		t.Run("InsertAll|"+statement, func(t *testing.T) {
			assert.Nil(t, ch.Error())

			// multiple insert
			assert.Nil(t, test.query.Insert(test.record, ch, ch, ch))
			assert.Equal(t, 3, reflect.ValueOf(test.record).Elem().Len())
		})
	}
}

// InsertSet tests insert specifications only using Set query.
func InsertSet(t *testing.T, repo grimoire.Repo) {
	user := User{}
	repo.From(users).MustSave(&user)
	now := time.Now()

	tests := []struct {
		query  grimoire.Query
		record interface{}
	}{
		{repo.From(users).Set("created_at", now).Set("updated_at", now).Set("name", "insert set"), &User{}},
		{repo.From(users).Set("created_at", now).Set("updated_at", now).Set("name", "insert set").Set("age", 100), &User{}},
		{repo.From(users).Set("created_at", now).Set("updated_at", now).Set("name", "insert set").Set("age", 100).Set("note", "note"), &User{}},
		{repo.From(users).Set("created_at", now).Set("updated_at", now).Set("note", "note"), &User{}},
		{repo.From(addresses).Set("created_at", now).Set("updated_at", now).Set("address", "address"), &Address{}},
		{repo.From(addresses).Set("created_at", now).Set("updated_at", now).Set("address", "address").Set("user_id", user.ID), &Address{}},
		{repo.From(addresses).Set("created_at", now).Set("updated_at", now).Set("user_id", user.ID), &Address{}},
	}

	for _, test := range tests {
		statement, _ := builder.Insert(test.query.Collection, test.query.Changes)

		t.Run("InsertSet|"+statement, func(t *testing.T) {
			assert.Nil(t, test.query.Insert(nil))
			assert.Nil(t, test.query.Insert(test.record))
		})
	}
}

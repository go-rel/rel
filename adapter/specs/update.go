package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

// Update tests update specifications.
func Update(t *testing.T, repo grimoire.Repo) {
	user := User{Name: "update"}
	repo.From(users).MustSave(&user)

	address := Address{Address: "update"}
	repo.From(addresses).MustSave(&address)

	tests := []struct {
		query  grimoire.Query
		record interface{}
		params params.Params
	}{
		{repo.From(users).Find(user.ID), &User{}, params.Map{"name": "insert", "age": 100}},
		{repo.From(users).Find(user.ID), &User{}, params.Map{"name": "insert", "age": 100, "note": "note"}},
		{repo.From(users).Find(user.ID), &User{}, params.Map{"note": "note"}},
		{repo.From(addresses).Find(address.ID), &Address{}, params.Map{"address": "address"}},
		{repo.From(addresses).Find(address.ID), &Address{}, params.Map{"user_id": user.ID}},
		{repo.From(addresses).Find(address.ID), &Address{}, params.Map{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		ch := changeset.Cast(test.record, test.params, []string{"name", "age", "note", "address", "user_id"})
		statement, _ := builder.Update(test.query.Collection, ch.Changes(), test.query.Condition)

		t.Run("Update|"+statement, func(t *testing.T) {
			assert.Nil(t, ch.Error())

			assert.Nil(t, test.query.Update(nil, ch))
			assert.Nil(t, test.query.Update(test.record, ch))
		})
	}
}

// UpdateWhere tests update specifications.
func UpdateWhere(t *testing.T, repo grimoire.Repo) {
	user := User{Name: "update all"}
	repo.From(users).MustSave(&user)

	address := Address{Address: "update all"}
	repo.From(addresses).MustSave(&address)

	tests := []struct {
		query  grimoire.Query
		schema interface{}
		record interface{}
		params params.Params
	}{
		{repo.From(users).Where(c.Eq(name, "update all")), User{}, &[]User{}, params.Map{"name": "insert", "age": 100}},
		{repo.From(addresses).Where(c.Eq(c.I("address"), "update_all")), Address{}, &[]Address{}, params.Map{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		ch := changeset.Cast(test.schema, test.params, []string{"name", "age", "note", "address", "user_id"})
		statement, _ := builder.Update(test.query.Collection, ch.Changes(), test.query.Condition)

		t.Run("UpdateWhere|"+statement, func(t *testing.T) {
			assert.Nil(t, ch.Error())

			assert.Nil(t, test.query.Update(nil, ch))
			assert.Nil(t, test.query.Update(test.record, ch))
		})
	}
}

// UpdateSet tests update specifications using Set query.
func UpdateSet(t *testing.T, repo grimoire.Repo) {
	user := User{Name: "update"}
	repo.From(users).MustSave(&user)

	address := Address{Address: "update"}
	repo.From(addresses).MustSave(&address)

	tests := []struct {
		query  grimoire.Query
		record interface{}
	}{
		{repo.From(users).Find(user.ID).Set("name", "update set"), &User{}},
		{repo.From(users).Find(user.ID).Set("name", "update set").Set("age", 18), &User{}},
		{repo.From(users).Find(user.ID).Set("note", "note set"), &User{}},
		{repo.From(addresses).Find(address.ID).Set("address", "address set"), &Address{}},
		{repo.From(addresses).Find(address.ID).Set("user_id", user.ID), &Address{}},
	}

	for _, test := range tests {
		statement, _ := builder.Update(test.query.Collection, test.query.Changes, test.query.Condition)

		t.Run("UpdateSet|"+statement, func(t *testing.T) {
			assert.Nil(t, test.query.Update(nil))
			assert.Nil(t, test.query.Update(test.record))
		})
	}
}

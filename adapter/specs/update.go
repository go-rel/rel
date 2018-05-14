package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sql"
	"github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
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
		params map[string]interface{}
	}{
		{repo.From(users).Find(user.ID), &User{}, map[string]interface{}{"name": "insert", "age": 100}},
		{repo.From(users).Find(user.ID), &User{}, map[string]interface{}{"name": "insert", "age": 100, "note": "note"}},
		{repo.From(users).Find(user.ID), &User{}, map[string]interface{}{"note": "note"}},
		{repo.From(addresses).Find(address.ID), &Address{}, map[string]interface{}{"address": "address"}},
		{repo.From(addresses).Find(address.ID), &Address{}, map[string]interface{}{"user_id": user.ID}},
		{repo.From(addresses).Find(address.ID), &Address{}, map[string]interface{}{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		ch := changeset.Cast(test.record, test.params, []string{"name", "age", "note", "address", "user_id"})
		statement, _ := sql.NewBuilder("?", false).Update(test.query.Collection, ch.Changes(), test.query.Condition)

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
		params map[string]interface{}
	}{
		{repo.From(users).Where(c.Eq(name, "update all")), User{}, &[]User{}, map[string]interface{}{"name": "insert", "age": 100}},
		{repo.From(addresses).Where(c.Eq(c.I("address"), "update_all")), Address{}, &[]Address{}, map[string]interface{}{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		ch := changeset.Cast(test.schema, test.params, []string{"name", "age", "note", "address", "user_id"})
		statement, _ := sql.NewBuilder("?", false).Update(test.query.Collection, ch.Changes(), test.query.Condition)

		t.Run("Update|"+statement, func(t *testing.T) {
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
		statement, _ := sql.NewBuilder("?", false).Update(test.query.Collection, test.query.Changes, test.query.Condition)

		t.Run("Update|"+statement, func(t *testing.T) {
			assert.Nil(t, test.query.Update(nil))
			assert.Nil(t, test.query.Update(test.record))
		})
	}
}

// UpdateConstraint tests update constraint specifications.
func UpdateConstraint(t *testing.T, repo grimoire.Repo) {
	user := User{}
	repo.From(users).MustSave(&user)

	repo.From(users).Set("slug", "update-taken").MustInsert(nil)

	tests := []struct {
		name  string
		query grimoire.Query
		field string
		code  int
	}{
		{"UniqueConstraintError", repo.From(users).Find(user.ID).Set("slug", "update-taken"), "slug", errors.UniqueConstraintErrorCode},
	}

	for _, test := range tests {
		t.Run("UpdateConstraint|"+test.name, func(t *testing.T) {
			checkConstraint(t, test.query.Insert(nil), test.code, test.field)
		})
	}
}

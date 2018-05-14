package specs

import (
	"reflect"
	"testing"
	"time"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sql"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

// Insert tests insert specifications.
func Insert(t *testing.T, repo grimoire.Repo) {
	user := User{}
	repo.From(users).MustSave(&user)

	tests := []struct {
		query  grimoire.Query
		record interface{}
		params map[string]interface{}
	}{
		{repo.From(users), &User{}, map[string]interface{}{}},
		{repo.From(users), &User{}, map[string]interface{}{"name": "insert", "age": 100}},
		{repo.From(users), &User{}, map[string]interface{}{"name": "insert", "age": 100, "note": "note"}},
		{repo.From(users), &User{}, map[string]interface{}{"note": "note"}},
		{repo.From(addresses), &Address{}, map[string]interface{}{}},
		{repo.From(addresses), &Address{}, map[string]interface{}{"address": "address"}},
		{repo.From(addresses), &Address{}, map[string]interface{}{"user_id": user.ID}},
		{repo.From(addresses), &Address{}, map[string]interface{}{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		ch := changeset.Cast(test.record, test.params, []string{"name", "age", "note", "address", "user_id"})
		statement, _ := sql.NewBuilder("?", false).Insert(test.query.Collection, ch.Changes())

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
		params map[string]interface{}
	}{
		{repo.From(users), User{}, &[]User{}, map[string]interface{}{}},
		{repo.From(users), User{}, &[]User{}, map[string]interface{}{"name": "insert", "age": 100}},
		{repo.From(users), User{}, &[]User{}, map[string]interface{}{"name": "insert", "age": 100, "note": "note"}},
		{repo.From(users), User{}, &[]User{}, map[string]interface{}{"note": "note"}},
		{repo.From(addresses), &Address{}, &[]Address{}, map[string]interface{}{}},
		{repo.From(addresses), &Address{}, &[]Address{}, map[string]interface{}{"address": "address"}},
		{repo.From(addresses), &Address{}, &[]Address{}, map[string]interface{}{"user_id": user.ID}},
		{repo.From(addresses), &Address{}, &[]Address{}, map[string]interface{}{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		ch := changeset.Cast(test.schema, test.params, []string{"name", "age", "note", "address", "user_id"})
		statement, _ := sql.NewBuilder("?", false).Insert(test.query.Collection, ch.Changes())

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
		statement, _ := sql.NewBuilder("?", false).Insert(test.query.Collection, test.query.Changes)

		t.Run("InsertSet|"+statement, func(t *testing.T) {
			assert.Nil(t, test.query.Insert(nil))
			assert.Nil(t, test.query.Insert(test.record))
		})
	}
}

// InsertConstraint tests insert constraint specifications.
func InsertConstraint(t *testing.T, repo grimoire.Repo) {
	repo.From(users).Set("slug", "insert-taken").MustInsert(nil)

	tests := []struct {
		name  string
		query grimoire.Query
		field string
		code  int
	}{
		{"UniqueConstraintError", repo.From(users).Set("slug", "insert-taken"), "slug", errors.UniqueConstraintErrorCode},
	}

	for _, test := range tests {
		t.Run("InsertConstraint|"+test.name, func(t *testing.T) {
			checkConstraint(t, test.query.Insert(nil), test.code, test.field)
		})
	}
}

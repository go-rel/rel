package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/changeset"
	"github.com/stretchr/testify/assert"
)

// Insert tests insert specifications.
func Insert(t *testing.T, repo grimoire.Repo) {
	user := User{}
	assert.Nil(t, repo.From(users).Put(&user))

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
		statement, _ := sqlutil.NewBuilder("?", false).Insert(test.query.Collection, ch.Changes())

		t.Run("Insert|"+statement, func(t *testing.T) {
			assert.Nil(t, ch.Error())

			assert.Nil(t, test.query.Insert(nil, ch))
			assert.Nil(t, test.query.Insert(test.record, ch))

			// multiple insert
			assert.Nil(t, test.query.Insert(nil, ch, ch, ch))
		})

		t.Run("InsertAll|"+statement, func(t *testing.T) {
			assert.Nil(t, ch.Error())

			// multiple insert
			assert.Nil(t, test.query.Insert(nil, ch, ch, ch))
		})
	}
}

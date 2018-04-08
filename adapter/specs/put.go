package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/changeset"
	"github.com/stretchr/testify/assert"
)

// PutInsert tests put insert specifications.
func PutInsert(t *testing.T, repo grimoire.Repo) {
	tests := []grimoire.Query{
		repo.From(users),
	}

	for _, query := range tests {
		record := User{Name: "put insert", Age: 100}
		ch := changeset.Change(record)
		statement, _ := sqlutil.NewBuilder("?", false).Insert(query.Collection, ch.Changes())

		t.Run("PutInsert|"+statement, func(t *testing.T) {
			assert.Nil(t, query.Put(&record))

			var result User
			assert.Nil(t, query.Find(record.ID).One(&result))
			assert.Equal(t, record, result)
		})
	}
}

// PutUpdate tests put update specifications.
func PutUpdate(t *testing.T, repo grimoire.Repo) {
	record := User{Name: "put update", Age: 100}
	assert.Nil(t, repo.From(users).Put(&record))
	assert.Nil(t, repo.From(users).Put(&User{Name: "put update", Age: 100}))
	assert.Nil(t, repo.From(users).Put(&User{Name: "put update", Age: 100}))

	tests := []grimoire.Query{
		repo.From(users).Find(record.ID),
		repo.From(users).Where(c.Eq(name, "put update")),
	}

	for _, query := range tests {
		statement, _ := sqlutil.NewBuilder("?", false).Update(query.Collection, map[string]interface{}{}, query.Condition)
		t.Run("PutUpdate|"+statement, func(t *testing.T) {
			var result []User
			assert.Nil(t, query.All(&result))
			count := len(result)
			assert.NotEqual(t, 0, count)

			record := []User{{Name: "put update", Age: 100}}
			assert.Nil(t, query.Put(&record))
			assert.Equal(t, count, len(record))
		})
	}
}

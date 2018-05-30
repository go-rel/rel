package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
	"github.com/stretchr/testify/assert"
)

// SaveInsert tests insert specifications.
func SaveInsert(t *testing.T, repo grimoire.Repo) {
	tests := []grimoire.Query{
		repo.From(users),
	}

	for _, query := range tests {
		record := User{Name: "save insert", Age: 100}

		t.Run("SaveInsert", func(t *testing.T) {
			assert.Nil(t, query.Save(&record))

			var result User
			assert.Nil(t, query.Find(record.ID).One(&result))
			assert.Equal(t, record, result)
		})
	}
}

// SaveInsertAll tests insert multiple recors specifications.
func SaveInsertAll(t *testing.T, repo grimoire.Repo) {
	tests := []grimoire.Query{
		repo.From(users),
	}

	for _, query := range tests {
		records := []User{
			{Name: "save insert all 1", Age: 100},
			{Name: "save insert all 2", Age: 100},
		}

		t.Run("SaveInsertAll", func(t *testing.T) {
			assert.Nil(t, query.Save(&records))
		})
	}
}

// SaveUpdate tests update specifications.
func SaveUpdate(t *testing.T, repo grimoire.Repo) {
	record := User{Name: "save update", Age: 100}
	repo.From(users).MustSave(&record)
	repo.From(users).MustSave(&User{Name: "save update", Age: 100})
	repo.From(users).MustSave(&User{Name: "save update", Age: 100})

	tests := []grimoire.Query{
		repo.From(users).Find(record.ID),
		repo.From(users).Where(c.Eq(name, "save update")),
	}

	for _, query := range tests {
		statement, _ := builder.Update(query.Collection, map[string]interface{}{}, query.Condition)
		t.Run("SaveUpdate|"+statement, func(t *testing.T) {
			var result []User
			assert.Nil(t, query.All(&result))
			count := len(result)
			assert.NotEqual(t, 0, count)

			record := []User{{Name: "save update", Age: 100}}
			assert.Nil(t, query.Save(&record))
			assert.Equal(t, count, len(record))
		})
	}
}

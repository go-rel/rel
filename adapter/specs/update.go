package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
)

// Update tests update specifications.
// TODO: atomic update
// TODO: update all
// TODO: update with assoc
func Update(t *testing.T, repo grimoire.Repo) {
	var (
		note    = "note"
		user    = User{Name: "update"}
		address = Address{Address: "update"}
	)

	repo.MustInsert(&user)
	repo.MustInsert(&address)

	tests := []interface{}{
		&User{ID: user.ID, Name: "changed", Age: 100},
		&User{ID: user.ID, Name: "changed", Age: 100, Note: &note},
		&User{ID: user.ID, Note: &note},
		&Address{ID: address.ID, Address: "address"},
		&Address{ID: address.ID, UserID: &user.ID},
		&Address{ID: address.ID, Address: "address", UserID: &user.ID},
	}

	for _, record := range tests {
		var (
			changes      = grimoire.BuildChanges(grimoire.Struct(record))
			statement, _ = builder.Update("collection", changes, where.Eq("id", 1))
		)

		t.Run("Update|"+statement, func(t *testing.T) {
			assert.Nil(t, repo.Update(record))
		})
	}
}

func UpdateExplicit(t *testing.T, repo grimoire.Repo) {
	var (
		user    = User{Name: "update"}
		address = Address{Address: "update"}
	)

	repo.MustInsert(&user)
	repo.MustInsert(&address)

	tests := []struct {
		record  interface{}
		changer grimoire.Changer
	}{
		{&user, grimoire.Map{"name": "changed", "age": 100}},
		{&user, grimoire.Map{"name": "changed", "age": 100, "note": "note"}},
		{&user, grimoire.Map{"note": "note"}},
		{&address, grimoire.Map{"address": "address"}},
		{&address, grimoire.Map{"user_id": user.ID}},
		{&address, grimoire.Map{"address": "address", "user_id": user.ID}},
	}

	for _, test := range tests {
		var (
			changes      = grimoire.BuildChanges(test.changer)
			statement, _ = builder.Update("collection", changes, where.Eq("id", 1))
		)

		t.Run("Update|"+statement, func(t *testing.T) {
			assert.Nil(t, repo.Update(test.record, test.changer))
		})
	}
}

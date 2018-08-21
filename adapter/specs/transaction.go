package specs

import (
	"testing"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

var input = params.Map{
	"name":   "whiteviolet",
	"gender": "male",
	"age":    18,
	"note":   "some note here",
	"addresses": []params.Map{
		{
			"address": "Aceh, Indonesia",
		},
		{
			"address": "Bandung, Indonesia",
		},
	},
}

// Transaction tests insert specifications.
func Transaction(t *testing.T, repo grimoire.Repo) {
	tests := []struct {
		name  string
		block func(*testing.T) func(grimoire.Repo) error
		err   error
	}{
		{"QueryAll", queryAll, nil},
		{"InsertWithAssoc", insertWithAssoc, nil},
		{"InsertWithAssocError", insertWithAssocError, errors.New("let's rollback", "", errors.NotFound)},
		{"InsertWithAssocPanic", insertWithAssocPanic, errors.New("let's rollback", "", errors.NotFound)},
		{"ReplaceAssoc", replaceAssoc, nil},
	}

	for _, tt := range tests {
		t.Run("Transaction|"+tt.name, func(t *testing.T) {
			assert.Equal(t, tt.err, repo.Transaction(tt.block(t)))
		})
	}
}

func queryAll(t *testing.T) func(repo grimoire.Repo) error {
	users := []User{}

	// transaction block
	return func(repo grimoire.Repo) error {
		repo.From("users").All(&users)

		return nil
	}
}

func insertWithAssoc(t *testing.T) func(repo grimoire.Repo) error {
	user := User{}

	ch := changeUser(user, input)
	assert.Nil(t, ch.Error())

	// transaction block
	return func(repo grimoire.Repo) error {
		repo.From("users").MustInsert(&user, ch)
		addresses := ch.Changes()["addresses"].([]*changeset.Changeset)
		repo.From("addresses").Set("user_id", user.ID).MustInsert(&user.Addresses, addresses...)

		return nil
	}
}

func insertWithAssocError(t *testing.T) func(repo grimoire.Repo) error {
	user := User{}

	ch := changeUser(user, input)
	assert.Nil(t, ch.Error())

	// transaction block
	return func(repo grimoire.Repo) error {
		repo.From("users").MustInsert(&user, ch)
		addresses := ch.Changes()["addresses"].([]*changeset.Changeset)
		repo.From("addresses").Set("user_id", user.ID).MustInsert(&user.Addresses, addresses...)

		// should rollback
		return errors.New("let's rollback", "", errors.NotFound)
	}
}

func insertWithAssocPanic(t *testing.T) func(repo grimoire.Repo) error {
	user := User{}

	ch := changeUser(user, input)
	assert.Nil(t, ch.Error())

	// transaction block
	return func(repo grimoire.Repo) error {
		repo.From("users").MustInsert(&user, ch)
		addresses := ch.Changes()["addresses"].([]*changeset.Changeset)
		repo.From("addresses").Set("user_id", user.ID).MustInsert(&user.Addresses, addresses...)

		// should rollback
		panic(errors.New("let's rollback", "", errors.NotFound))
	}
}

func replaceAssoc(t *testing.T) func(repo grimoire.Repo) error {
	user := User{}

	ch := changeUser(user, input)
	assert.Nil(t, ch.Error())

	// transaction block
	return func(repo grimoire.Repo) error {
		repo.From("users").MustOne(&user)
		addresses := ch.Changes()["addresses"].([]*changeset.Changeset)

		repo.From("addresses").Where(c.Eq(c.I("user_id"), user.ID)).MustDelete()
		repo.From("addresses").Set("user_id", user.ID).MustInsert(&user.Addresses, addresses...)

		return nil
	}
}

func changeUser(user interface{}, params params.Params) *changeset.Changeset {
	ch := changeset.Cast(user, params, []string{
		"name",
		"gender",
		"age",
		"note",
	})
	changeset.CastAssoc(ch, "addresses", changeAddress)
	return ch
}

func changeAddress(address interface{}, params params.Params) *changeset.Changeset {
	ch := changeset.Cast(address, params, []string{"address"})
	return ch
}

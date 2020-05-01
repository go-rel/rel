package rel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDirty(t *testing.T) {
	var (
		now  = time.Now()
		user = User{
			ID:        1,
			Name:      "User 1",
			Age:       20,
			UpdatedAt: now,
			CreatedAt: now,
		}
		snapshot = []interface{}{1, "User 1", 20, now, now}
		doc      = NewDocument(&user)
	)

	t.Run("init", func(t *testing.T) {
		user.Init(doc)

		assert.Equal(t, snapshot, user.Dirty.snapshot)
		assert.Empty(t, user.Dirty.Changes())
	})

	t.Run("update", func(t *testing.T) {
		user.Name = "User 2"
		user.Age = 21

		assert.Equal(t, snapshot, user.Dirty.snapshot)
		assert.Equal(t, map[string]Pair{
			"name": Pair{"User 1", "User 2"},
			"age":  Pair{20, 21},
		}, user.Dirty.Changes())
	})
}

func TestDirty_ptr(t *testing.T) {
	var (
		dirty   Dirty
		userID  = 2
		address = Address{
			ID:     1,
			UserID: &userID,
		}
		snapshot = []interface{}{1, 2, "", Notes(""), nil}
		doc      = NewDocument(&address)
	)

	t.Run("init", func(t *testing.T) {
		dirty.Init(doc)

		assert.Equal(t, snapshot, dirty.snapshot)
		assert.Empty(t, dirty.Changes())
	})

	t.Run("update", func(t *testing.T) {
		userID = 3
		address.UserID = &userID

		assert.Equal(t, snapshot, dirty.snapshot)
		assert.Equal(t, map[string]Pair{
			"user_id": Pair{2, 3},
		}, dirty.Changes())

	})
}

func TestDirty_belongsTo(t *testing.T) {
	var (
		dirty   Dirty
		address = Address{
			User: &User{
				ID:   1,
				Name: "User 1",
			},
		}
		snapshot = []interface{}{1, "User 1", 0, time.Time{}, time.Time{}}
		doc      = NewDocument(&address)
	)

	t.Run("init", func(t *testing.T) {
		dirty.Init(doc)
		dirty = *dirty.assoc["user"]

		assert.Equal(t, dirty, address.User.Dirty)
		assert.Equal(t, snapshot, dirty.snapshot)
		assert.Empty(t, dirty.Changes())
	})

	t.Run("update", func(t *testing.T) {
		address.User.Name = "User Satu"

		assert.Equal(t, snapshot, dirty.snapshot)
		assert.Equal(t, map[string]Pair{
			"name": Pair{"User 1", "User Satu"},
		}, dirty.Changes())
	})
}

func TestDirty_hasOne(t *testing.T) {
	var (
		dirty Dirty
		user  = User{
			ID: 1,
			Address: Address{
				ID:     1,
				Street: "Grove Street",
				Notes:  "HQ",
			},
		}
		snapshot = []interface{}{1, nil, "Grove Street", Notes("HQ"), nil}
		doc      = NewDocument(&user)
	)

	t.Run("init", func(t *testing.T) {
		dirty.Init(doc)
		dirty = *dirty.assoc["address"]

		assert.Equal(t, snapshot, dirty.snapshot)
		assert.Empty(t, dirty.Changes())
	})

	t.Run("update", func(t *testing.T) {
		user.Address.UserID = &user.ID
		user.Address.Street = "Grove Street Blvd"
		user.Address.Notes = Notes("Home")

		assert.Equal(t, snapshot, dirty.snapshot)
		assert.Equal(t, map[string]Pair{
			"user_id": Pair{nil, user.ID},
			"street":  Pair{"Grove Street", "Grove Street Blvd"},
			"notes":   Pair{Notes("HQ"), Notes("Home")},
		}, dirty.Changes())
	})
}

func TestDirty_hasMany(t *testing.T) {
	var (
		dirties map[interface{}]*Dirty
		user    = User{
			ID: 1,
			Transactions: []Transaction{
				{ID: 11, Item: "Book", Status: "pending"},
				{ID: 12, Item: "Eraser", Status: "pending"},
			},
		}
		snapshots = [][]interface{}{
			{11, "Book", Status("pending"), 0},
			{12, "Eraser", Status("pending"), 0},
		}
		doc = NewDocument(&user)
	)

	t.Run("init", func(t *testing.T) {
		user.Dirty.Init(doc)
		dirties = user.Dirty.assocMany["transactions"]

		assert.Equal(t, snapshots[0], dirties[11].snapshot)
		assert.Equal(t, snapshots[1], dirties[12].snapshot)

		assert.Empty(t, dirties[11].Changes())
		assert.Empty(t, dirties[12].Changes())
	})

	t.Run("update", func(t *testing.T) {
		user.Transactions[0].Status = "paid"
		// replaced struct is new, so there's no dirty states to check.
		user.Transactions[1] = Transaction{Item: "Paper", Status: "pending"}

		assert.Equal(t, map[string]Pair{
			"status": Pair{Status("pending"), Status("paid")},
		}, dirties[11].Changes())
	})
}

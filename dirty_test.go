package rel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDirty(t *testing.T) {
	var (
		ts   = time.Now()
		user = User{
			ID:        1,
			Name:      "User 1",
			Age:       20,
			UpdatedAt: ts,
			CreatedAt: ts,
		}
		snapshot = []interface{}{1, "User 1", 20, ts, ts}
		doc      = NewDocument(&user)
	)

	t.Run("init", func(t *testing.T) {
		user.Init(doc)

		assert.Equal(t, snapshot, user.Dirty.snapshot)
		assert.Empty(t, user.Dirty.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, user))
	})

	t.Run("update", func(t *testing.T) {
		user.Name = "User 2"
		user.Age = 21

		assert.Equal(t, snapshot, user.Dirty.snapshot)
		assert.Equal(t, map[string]pair{
			"name": pair{"User 1", "User 2"},
			"age":  pair{20, 21},
		}, user.Dirty.Changes())
	})

	t.Run("apply dirty", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{
				"name":       Set("name", "User 2"),
				"age":        Set("age", 21),
				"updated_at": Set("updated_at", now()),
			},
			Assoc: map[string]AssocModification{},
		}, Apply(doc, user))
	})

	t.Run("apply new assoc as structset", func(t *testing.T) {
		user.Address.Street = "Grove Street"

		assert.Equal(t, Modification{
			Modifies: map[string]Modify{
				"name":       Set("name", "User 2"),
				"age":        Set("age", 21),
				"updated_at": Set("updated_at", now()),
			},
			Assoc: map[string]AssocModification{
				"address": AssocModification{
					Modifications: []Modification{
						{
							Modifies: map[string]Modify{
								"user_id":    Set("user_id", nil),
								"street":     Set("street", "Grove Street"),
								"notes":      Set("notes", Notes("")),
								"deleted_at": Set("deleted_at", nil),
							},
							Assoc: map[string]AssocModification{},
						},
					},
				},
			},
		}, Apply(doc, user))
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

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, dirty))
	})

	t.Run("update", func(t *testing.T) {
		userID = 3
		address.UserID = &userID

		assert.Equal(t, snapshot, dirty.snapshot)
		assert.Equal(t, map[string]pair{
			"user_id": pair{2, 3},
		}, dirty.Changes())
	})

	t.Run("apply dirty", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{
				"user_id": Set("user_id", 3),
			},
			Assoc: map[string]AssocModification{},
		}, Apply(doc, dirty))
	})

	t.Run("apply new assoc as structset", func(t *testing.T) {
		address.User = &User{Name: "User 3"}

		assert.Equal(t, Modification{
			Modifies: map[string]Modify{
				"user_id": Set("user_id", 3),
			},
			Assoc: map[string]AssocModification{
				"user": AssocModification{
					Modifications: []Modification{
						{
							Modifies: map[string]Modify{
								"age":        Set("age", 0),
								"name":       Set("name", "User 3"),
								"created_at": Set("created_at", now()),
								"updated_at": Set("updated_at", now()),
							},
							Assoc: map[string]AssocModification{},
						},
					},
				},
			},
		}, Apply(doc, dirty))
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

		assert.Equal(t, *dirty.assoc["user"], address.User.Dirty)
		assert.Equal(t, snapshot, dirty.assoc["user"].snapshot)
		assert.Empty(t, dirty.assoc["user"].Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, dirty))
	})

	t.Run("update", func(t *testing.T) {
		address.User.Name = "User Satu"

		assert.Equal(t, snapshot, dirty.assoc["user"].snapshot)
		assert.Equal(t, map[string]pair{
			"name": pair{"User 1", "User Satu"},
		}, dirty.assoc["user"].Changes())
	})

	t.Run("apply dirty", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc: map[string]AssocModification{
				"user": AssocModification{
					Modifications: []Modification{
						{
							Modifies: map[string]Modify{
								"name":       Set("name", "User Satu"),
								"updated_at": Set("updated_at", now()),
							},
							Assoc: map[string]AssocModification{},
						},
					},
				},
			},
		}, Apply(doc, dirty))
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

		assert.Equal(t, snapshot, dirty.assoc["address"].snapshot)
		assert.Empty(t, dirty.assoc["address"].Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, dirty))
	})

	t.Run("update", func(t *testing.T) {
		user.Address.UserID = &user.ID
		user.Address.Street = "Grove Street Blvd"
		user.Address.Notes = Notes("Home")

		assert.Equal(t, snapshot, dirty.assoc["address"].snapshot)
		assert.Equal(t, map[string]pair{
			"user_id": pair{nil, user.ID},
			"street":  pair{"Grove Street", "Grove Street Blvd"},
			"notes":   pair{Notes("HQ"), Notes("Home")},
		}, dirty.assoc["address"].Changes())
	})

	t.Run("apply dirty", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc: map[string]AssocModification{
				"address": AssocModification{
					Modifications: []Modification{
						{
							Modifies: map[string]Modify{
								"user_id": Set("user_id", user.ID),
								"street":  Set("street", "Grove Street Blvd"),
								"notes":   Set("notes", Notes("Home")),
							},
							Assoc: map[string]AssocModification{},
						},
					},
				},
			},
		}, Apply(doc, dirty))
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

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, user))
	})

	t.Run("update", func(t *testing.T) {
		user.Transactions[0].Status = "paid"
		// replaced struct is new, so there's no dirty states to check.
		user.Transactions[1] = Transaction{Item: "Paper", Status: "pending"}

		assert.Equal(t, map[string]pair{
			"status": pair{Status("pending"), Status("paid")},
		}, dirties[11].Changes())
	})

	t.Run("apply dirty", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc: map[string]AssocModification{
				"transactions": AssocModification{
					Modifications: []Modification{
						{
							Modifies: map[string]Modify{
								"status": Set("status", Status("paid")),
							},
							Assoc: map[string]AssocModification{},
						},
						{
							Modifies: map[string]Modify{
								"item":    Set("item", "Paper"),
								"status":  Set("status", Status("pending")),
								"user_id": Set("user_id", 0),
							},
							Assoc: map[string]AssocModification{},
						},
					},
					DeletedIDs: []interface{}{12},
				},
			},
		}, Apply(doc, user))
	})

	t.Run("apply clear assoc", func(t *testing.T) {
		user.Transactions = []Transaction{}
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc: map[string]AssocModification{
				"transactions": AssocModification{
					Modifications: []Modification{},
					DeletedIDs:    []interface{}{11, 12},
				},
			},
		}, Apply(doc, user))
	})
}

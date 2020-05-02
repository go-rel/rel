package rel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkSmallSliceLookup(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var (
			index  = 0
			values = []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			fields = []string{"field1", "field2", "field3", "field4", "field6", "field7", "field8", "field9", "field10"}
		)

		for i := range fields {
			if fields[i] == "field10" {
				index = i
				break
			}
		}
		_ = values[index]
		_ = values
		_ = fields
	}
}

func BenchmarkSmallMapLookup(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var (
			values = map[string]interface{}{
				"field1":  1,
				"field2":  2,
				"field3":  3,
				"field4":  4,
				"field5":  5,
				"field6":  6,
				"field7":  7,
				"field8":  8,
				"field9":  9,
				"field10": 10,
			}
		)
		_ = values["fields10"]
		_ = values
	}
}

func BenchmarkChangeset(b *testing.B) {
	var (
		user = User{
			ID:   1,
			Name: "Luffy",
			Age:  20,
			Transactions: []Transaction{
				{ID: 1, Item: "Sword"},
				{ID: 2, Item: "Shield"},
			},
			Address: Address{
				ID:     1,
				Street: "Grove Street",
			},
			CreatedAt: time.Now(),
		}
		doc = NewDocument(&user)
	)

	for n := 0; n < b.N; n++ {
		changeset := NewChangeset(&user)
		user.Name = "Zoro"

		Apply(doc, changeset)
	}
}

func BenchmarkChangeset_assoc(b *testing.B) {
	var (
		user = User{
			ID:   1,
			Name: "Luffy",
			Age:  20,
			Transactions: []Transaction{
				{ID: 1, Item: "Sword"},
				{ID: 2, Item: "Shield"},
			},
			Address: Address{
				ID:     1,
				Street: "Grove Street",
			},
			CreatedAt: time.Now(),
		}
		doc = NewDocument(&user)
	)

	for n := 0; n < b.N; n++ {
		changeset := NewChangeset(&user)
		user.Name = "Zoro"
		user.Transactions[0].Item = "Sake"
		user.Address.Street = "Thousand Sunny"

		Apply(doc, changeset)
	}
}

func TestChangeset(t *testing.T) {
	var (
		ts   = time.Now()
		user = User{
			ID:        1,
			Name:      "User 1",
			Age:       20,
			UpdatedAt: ts,
			CreatedAt: ts,
		}
		snapshot  = []interface{}{1, "User 1", 20, ts, ts}
		doc       = NewDocument(&user)
		changeset = NewChangeset(&user)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		user.Name = "User 2"
		user.Age = 21

		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Equal(t, map[string]pair{
			"name": pair{"User 1", "User 2"},
			"age":  pair{20, 21},
		}, changeset.Changes())

		assert.True(t, changeset.FieldChanged("name"))
		assert.True(t, changeset.FieldChanged("age"))

		assert.False(t, changeset.FieldChanged("id"))
		assert.False(t, changeset.FieldChanged("created_at"))
		assert.False(t, changeset.FieldChanged("unknown"))
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{
				"name":       Set("name", "User 2"),
				"age":        Set("age", 21),
				"updated_at": Set("updated_at", now()),
			},
			Assoc: map[string]AssocModification{},
		}, Apply(doc, changeset))
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
		}, Apply(doc, changeset))
	})
}

func TestChangeset_ptr(t *testing.T) {
	var (
		userID  = 2
		address = Address{
			ID:     1,
			UserID: &userID,
		}
		snapshot  = []interface{}{1, 2, "", Notes(""), nil}
		doc       = NewDocument(&address)
		changeset = NewChangeset(&address)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		userID = 3
		address.UserID = &userID

		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Equal(t, map[string]pair{
			"user_id": pair{2, 3},
		}, changeset.Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{
				"user_id": Set("user_id", 3),
			},
			Assoc: map[string]AssocModification{},
		}, Apply(doc, changeset))
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
		}, Apply(doc, changeset))
	})
}

func TestChangeset_belongsTo(t *testing.T) {
	var (
		address = Address{
			User: &User{
				ID:   1,
				Name: "User 1",
			},
		}
		snapshot  = []interface{}{1, "User 1", 0, time.Time{}, time.Time{}}
		doc       = NewDocument(&address)
		changeset = NewChangeset(&address)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.assoc["user"].snapshot)
		assert.Empty(t, changeset.assoc["user"].Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		address.User.Name = "User Satu"

		assert.Equal(t, snapshot, changeset.assoc["user"].snapshot)
		assert.Equal(t, map[string]pair{
			"name": pair{"User 1", "User Satu"},
		}, changeset.assoc["user"].Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
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
		}, Apply(doc, changeset))
	})
}

func TestChangeset_hasOne(t *testing.T) {
	var (
		user = User{
			ID: 1,
			Address: Address{
				ID:     1,
				Street: "Grove Street",
				Notes:  "HQ",
			},
		}
		snapshot  = []interface{}{1, nil, "Grove Street", Notes("HQ"), nil}
		doc       = NewDocument(&user)
		changeset = NewChangeset(&user)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.assoc["address"].snapshot)
		assert.Empty(t, changeset.assoc["address"].Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		user.Address.UserID = &user.ID
		user.Address.Street = "Grove Street Blvd"
		user.Address.Notes = Notes("Home")

		assert.Equal(t, snapshot, changeset.assoc["address"].snapshot)
		assert.Equal(t, map[string]pair{
			"user_id": pair{nil, user.ID},
			"street":  pair{"Grove Street", "Grove Street Blvd"},
			"notes":   pair{Notes("HQ"), Notes("Home")},
		}, changeset.assoc["address"].Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
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
		}, Apply(doc, changeset))
	})
}

func TestChangeset_hasMany(t *testing.T) {
	var (
		user = User{
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
		doc       = NewDocument(&user)
		changeset = NewChangeset(&user)
	)

	t.Run("snapshot", func(t *testing.T) {
		trxch := changeset.assocMany["transactions"]

		assert.Equal(t, snapshots[0], trxch[11].snapshot)
		assert.Equal(t, snapshots[1], trxch[12].snapshot)

		assert.Empty(t, trxch[11].Changes())
		assert.Empty(t, trxch[12].Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Modification{
			Modifies: map[string]Modify{},
			Assoc:    map[string]AssocModification{},
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		trxch := changeset.assocMany["transactions"]

		user.Transactions[0].Status = "paid"
		// replaced struct is new, so there's no changeset states to check.
		user.Transactions[1] = Transaction{Item: "Paper", Status: "pending"}

		assert.Equal(t, map[string]pair{
			"status": pair{Status("pending"), Status("paid")},
		}, trxch[11].Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
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
		}, Apply(doc, changeset))
	})
}

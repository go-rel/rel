package rel

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkSmallSliceLookup(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var (
			index  = 0
			values = []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
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
			values = map[string]any{
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
		snapshot  = []any{1, "User 1", 20, ts, ts}
		doc       = NewDocument(&user)
		changeset = NewChangeset(&user)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		user.Name = "User 2"
		user.Age = 21

		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Equal(t, map[string]any{
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
		assert.Equal(t, Mutation{
			Cascade: true,
			Mutates: map[string]Mutate{
				"name":       Set("name", "User 2"),
				"age":        Set("age", 21),
				"updated_at": Set("updated_at", Now()),
			},
		}, Apply(doc, changeset))
	})
}

func TestChangeset_byte_slice(t *testing.T) {
	var (
		ts   = time.Now()
		user = extendedUser{
			ID:        1,
			Password:  []byte("foo"),
			Metadata:  []byte(`{"baz":"foo"}`),
			CreatedAt: ts,
			UpdatedAt: ts,
		}
		snapshot  = []any{1, []byte("foo"), json.RawMessage(`{"baz":"foo"}`), ts, ts}
		doc       = NewDocument(&user)
		changeset = NewChangeset(&user)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		user.Password = []byte("bar")
		user.Metadata = []byte("{}")

		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Equal(t, map[string]any{
			"password": pair{[]byte("foo"), []byte("bar")},
			"metadata": pair{json.RawMessage(`{"baz":"foo"}`), json.RawMessage("{}")},
		}, changeset.Changes())

		assert.True(t, changeset.FieldChanged("password"))
		assert.True(t, changeset.FieldChanged("metadata"))

		assert.False(t, changeset.FieldChanged("id"))
		assert.False(t, changeset.FieldChanged("created_at"))
		assert.False(t, changeset.FieldChanged("unknown"))
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
			Mutates: map[string]Mutate{
				"password":   Set("password", []byte("bar")),
				"metadata":   Set("metadata", json.RawMessage(`{}`)),
				"updated_at": Set("updated_at", Now()),
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
		snapshot  = []any{1, 2, "", Notes(""), nil}
		doc       = NewDocument(&address)
		changeset = NewChangeset(&address)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		userID = 3
		address.UserID = &userID

		assert.Equal(t, snapshot, changeset.snapshot)
		assert.Equal(t, map[string]any{
			"user_id": pair{2, 3},
		}, changeset.Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
			Mutates: map[string]Mutate{
				"user_id": Set("user_id", 3),
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
		snapshot  = []any{1, "User 1", 0, time.Time{}, time.Time{}}
		doc       = NewDocument(&address)
		changeset = NewChangeset(&address)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.assoc["user"].snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		address.User.Name = "User Satu"

		assert.Equal(t, snapshot, changeset.assoc["user"].snapshot)
		assert.Equal(t, map[string]any{
			"user": map[string]any{
				"name": pair{"User 1", "User Satu"},
			},
		}, changeset.Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
			Assoc: map[string]AssocMutation{
				"user": {
					Mutations: []Mutation{
						{
							Cascade: true,
							Mutates: map[string]Mutate{
								"name":       Set("name", "User Satu"),
								"updated_at": Set("updated_at", Now()),
							},
						},
					},
				},
			},
		}, Apply(doc, changeset))
	})
}

func TestChangeset_belongsTo_new(t *testing.T) {
	var (
		address = Address{
			ID: 1,
		}
		doc       = NewDocument(&address)
		changeset = NewChangeset(&address)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Nil(t, changeset.assoc["user"].snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
		}, Apply(doc, changeset))
	})

	t.Run("create", func(t *testing.T) {
		address.User = &User{
			Name: "User Satu",
			Age:  20,
		}

		assert.Nil(t, changeset.assoc["user"].snapshot)
		assert.Equal(t, map[string]any{
			"user": map[string]any{
				"id":         pair{nil, 0},
				"name":       pair{nil, "User Satu"},
				"age":        pair{nil, 20},
				"created_at": pair{nil, time.Time{}},
				"updated_at": pair{nil, time.Time{}},
			},
		}, changeset.Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
			Assoc: map[string]AssocMutation{
				"user": {
					Mutations: []Mutation{
						{
							Cascade: true,
							Mutates: map[string]Mutate{
								"name":       Set("name", "User Satu"),
								"age":        Set("age", 20),
								"created_at": Set("created_at", Now()),
								"updated_at": Set("updated_at", Now()),
							},
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
		snapshot  = []any{1, nil, "Grove Street", Notes("HQ"), nil}
		doc       = NewDocument(&user)
		changeset = NewChangeset(&user)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Equal(t, snapshot, changeset.assoc["address"].snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		user.Address.UserID = &user.ID
		user.Address.Street = "Grove Street Blvd"
		user.Address.Notes = Notes("Home")

		assert.Equal(t, snapshot, changeset.assoc["address"].snapshot)
		assert.Equal(t, map[string]any{
			"address": map[string]any{
				"user_id": pair{nil, user.ID},
				"street":  pair{"Grove Street", "Grove Street Blvd"},
				"notes":   pair{Notes("HQ"), Notes("Home")},
			},
		}, changeset.Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
			Assoc: map[string]AssocMutation{
				"address": {
					Mutations: []Mutation{
						{
							Cascade: true,
							Mutates: map[string]Mutate{
								"user_id": Set("user_id", user.ID),
								"street":  Set("street", "Grove Street Blvd"),
								"notes":   Set("notes", Notes("Home")),
							},
						},
					},
				},
			},
		}, Apply(doc, changeset))
	})
}

func TestChangeset_hasOne_new(t *testing.T) {
	var (
		user = User{
			ID: 1,
		}
		doc       = NewDocument(&user)
		changeset = NewChangeset(&user)
	)

	t.Run("snapshot", func(t *testing.T) {
		assert.Nil(t, changeset.assoc["address"].snapshot)
		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
		}, Apply(doc, changeset))
	})

	t.Run("create", func(t *testing.T) {
		user.Address = Address{
			UserID: &user.ID,
			Street: "Grove Street Blvd",
			Notes:  Notes("Home"),
		}

		assert.Nil(t, changeset.assoc["address"].snapshot)
		assert.Equal(t, map[string]any{
			"address": map[string]any{
				"id":      pair{nil, 0},
				"user_id": pair{nil, user.ID},
				"street":  pair{nil, "Grove Street Blvd"},
				"notes":   pair{nil, Notes("Home")},
			},
		}, changeset.Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
			Assoc: map[string]AssocMutation{
				"address": {
					Mutations: []Mutation{
						{
							Cascade: true,
							Mutates: map[string]Mutate{
								"user_id":    Set("user_id", user.ID),
								"street":     Set("street", "Grove Street Blvd"),
								"notes":      Set("notes", Notes("Home")),
								"deleted_at": Set("deleted_at", nil),
							},
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
		snapshots = [][]any{
			{11, "Book", Status("pending"), 0, 0},
			{12, "Eraser", Status("pending"), 0, 0},
		}
		doc       = NewDocument(&user)
		changeset = NewChangeset(&user)
	)

	t.Run("snapshot", func(t *testing.T) {
		trxch := changeset.assocMany["transactions"]

		assert.Equal(t, snapshots[0], trxch[11].snapshot)
		assert.Equal(t, snapshots[1], trxch[12].snapshot)

		assert.Empty(t, changeset.Changes())
	})

	t.Run("apply clean", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
		}, Apply(doc, changeset))
	})

	t.Run("update", func(t *testing.T) {
		user.Transactions[0].Status = "paid"
		user.Transactions[1] = Transaction{Item: "Paper", Status: "pending"}

		assert.Equal(t, map[string]any{
			"transactions": []map[string]any{
				{
					"status": pair{Status("pending"), Status("paid")},
				},
				{
					"id":         pair{nil, 0},
					"item":       pair{nil, "Paper"},
					"status":     pair{nil, Status("pending")},
					"user_id":    pair{nil, 0},
					"address_id": pair{nil, 0},
				},
				{
					"id":         pair{12, nil},
					"item":       pair{"Eraser", nil},
					"status":     pair{Status("pending"), nil},
					"user_id":    pair{0, nil},
					"address_id": pair{0, nil},
				},
			},
		}, changeset.Changes())
	})

	t.Run("apply changeset", func(t *testing.T) {
		assert.Equal(t, Mutation{
			Cascade: true,
			Assoc: map[string]AssocMutation{
				"transactions": {
					Mutations: []Mutation{
						{
							Cascade: true,
							Mutates: map[string]Mutate{
								"status": Set("status", Status("paid")),
							},
						},
						{
							Cascade: true,
							Mutates: map[string]Mutate{
								"item":       Set("item", "Paper"),
								"status":     Set("status", Status("pending")),
								"user_id":    Set("user_id", 0),
								"address_id": Set("address_id", 0),
							},
						},
					},
					DeletedIDs: []any{12},
				},
			},
		}, Apply(doc, changeset))
	})
}

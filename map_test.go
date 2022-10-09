package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	var (
		user User
		doc  = NewDocument(&user)
		data = Map{
			"name": "Luffy",
			"age":  20,
			"transactions": []Map{
				{"item": "Sword"},
				{"item": "Shield"},
			},
			"address": Map{
				"street": "Grove Street",
			},
		}
		userMutation = Apply(NewDocument(&User{}),
			Set("name", "Luffy"),
			Set("age", 20),
		)
		transaction1Mutation = Apply(NewDocument(&Transaction{}),
			Set("item", "Sword"),
		)
		transaction2Mutation = Apply(NewDocument(&Transaction{}),
			Set("item", "Shield"),
		)
		addressMutation = Apply(NewDocument(&Address{}),
			Set("street", "Grove Street"),
		)
	)

	userMutation.SetAssoc("transactions", transaction1Mutation, transaction2Mutation)
	userMutation.SetAssoc("address", addressMutation)
	userMutation.SetDeletedIDs("transactions", []any{})

	assert.Equal(t, userMutation, Apply(doc, data))
	assert.Equal(t, User{
		Name: "Luffy",
		Age:  20,
		Transactions: []Transaction{
			{Item: "Sword"},
			{Item: "Shield"},
		},
		Address: Address{
			Street: "Grove Street",
		},
	}, user)
}

func TestMap_CascadeDisabled(t *testing.T) {
	var (
		user User
		doc  = NewDocument(&user)
		data = Map{
			"name": "Luffy",
			"age":  20,
			"transactions": []Map{
				{"item": "Sword"},
				{"item": "Shield"},
			},
			"address": Map{
				"street": "Grove Street",
			},
		}
		userMutation = Apply(NewDocument(&User{}),
			Cascade(false),
			Set("name", "Luffy"),
			Set("age", 20),
		)
	)

	assert.Equal(t, userMutation, Apply(doc, Cascade(false), data))
	assert.Equal(t, User{
		Name: "Luffy",
		Age:  20,
	}, user)
}

func TestMap_update(t *testing.T) {
	var (
		user = User{
			ID: 1,
			Transactions: []Transaction{
				{ID: 2},
				{ID: 3},
			},
			Address: Address{
				ID: 4,
			},
		}
		doc  = NewDocument(&user)
		data = Map{
			"name": "Luffy",
			"age":  20,
			"transactions": []Map{
				{"id": 2, "item": "Sword"},
				{"id": 3, "item": "Shield"},
			},
			"address": Map{
				"street": "Grove Street",
			},
		}
		userMutation = Apply(NewDocument(&User{}),
			Set("name", "Luffy"),
			Set("age", 20),
		)
		transaction1Mutation = Apply(NewDocument(&Transaction{}),
			Set("item", "Sword"),
		)
		transaction2Mutation = Apply(NewDocument(&Transaction{}),
			Set("item", "Shield"),
		)
		addressMutation = Apply(NewDocument(&Address{}),
			Set("street", "Grove Street"),
		)
	)

	userMutation.SetAssoc("transactions", transaction1Mutation, transaction2Mutation)
	userMutation.SetAssoc("address", addressMutation)
	userMutation.SetDeletedIDs("transactions", []any{})

	assert.Equal(t, userMutation, Apply(doc, data))
	assert.Equal(t, User{
		ID:   1,
		Name: "Luffy",
		Age:  20,
		Transactions: []Transaction{
			{ID: 2, Item: "Sword"},
			{ID: 3, Item: "Shield"},
		},
		Address: Address{
			ID:     4,
			Street: "Grove Street",
		},
	}, user)
}

func TestMap_hasManyUpdateDeleteInsert(t *testing.T) {
	var (
		user = User{
			Transactions: []Transaction{
				{ID: 2},
				{ID: 3},
			},
		}
		doc  = NewDocument(&user)
		data = Map{
			"transactions": []Map{
				{"item": "Sword"},
				{"id": 3, "item": "Shield"},
			},
		}
		userMutation         = Mutation{Cascade: true}
		transaction1Mutation = Apply(NewDocument(&Transaction{}),
			Set("item", "Sword"),
		)
		transaction2Mutation = Apply(NewDocument(&Transaction{}),
			Set("item", "Shield"),
		)
	)

	userMutation.SetAssoc("transactions", transaction2Mutation, transaction1Mutation)
	userMutation.SetDeletedIDs("transactions", []any{2})

	assert.Equal(t, userMutation, Apply(doc, data))
	assert.Equal(t, User{
		Transactions: []Transaction{
			{ID: 3, Item: "Shield"},
			{Item: "Sword"},
		},
	}, user)
}

func TestMap_hasManyUpdateNotLoaded(t *testing.T) {
	var (
		user = User{
			Transactions: []Transaction{
				{ID: 2},
			},
		}
		doc  = NewDocument(&user)
		data = Map{
			"transactions": []Map{
				{"id": 3, "item": "Sword"},
			},
		}
	)

	assert.PanicsWithValue(t, "rel: cannot update has many assoc that is not loaded or doesn't belong to this record", func() {
		Apply(doc, data)
	})
}

func TestMap_hasManyWrongType(t *testing.T) {
	var (
		user = User{
			Transactions: []Transaction{},
		}
		doc  = NewDocument(&user)
		data = Map{
			"transactions": Map{"item": "Sword"},
		}
	)

	assert.Panics(t, func() {
		Apply(doc, data)
	})
}

func TestMap_wrongType(t *testing.T) {
	var (
		user User
		doc  = NewDocument(&user)
		data = Map{
			"name": false,
		}
	)

	assert.Panics(t, func() {
		Apply(doc, Cascade(false), data)
	})
}

func TestMap_replacingPrimaryKey(t *testing.T) {
	var (
		user = User{ID: 1}
		doc  = NewDocument(&user)
		data = Map{
			"id": 2,
		}
	)

	assert.Panics(t, func() {
		Apply(doc, Cascade(false), data)
	})
}

func TestMap_String(t *testing.T) {
	var (
		data = Map{
			"name": "Luffy",
			"age":  20,
			"transactions": []Map{
				{"item": "Sword"},
				{"item": "Shield"},
			},
			"address": Map{
				"street": "Grove Street",
			},
		}
	)

	assert.Contains(t, data.String(), "\"name\": \"Luffy\"")
	assert.Contains(t, data.String(), "\"age\": 20")
	assert.Contains(t, data.String(), "\"transactions\": []rel.Map{rel.Map{\"item\": \"Sword\"}, rel.Map{\"item\": \"Shield\"}}")
	assert.Contains(t, data.String(), "\"address\": rel.Map{\"street\": \"Grove Street\"}")
}

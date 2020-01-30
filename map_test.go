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
		userModification = Apply(NewDocument(&User{}),
			Set("name", "Luffy"),
			Set("age", 20),
		)
		transaction1Modification = Apply(NewDocument(&Transaction{}),
			Set("item", "Sword"),
		)
		transaction2Modification = Apply(NewDocument(&Transaction{}),
			Set("item", "Shield"),
		)
		addressModification = Apply(NewDocument(&Address{}),
			Set("street", "Grove Street"),
		)
	)

	userModification.SetAssoc("transactions", transaction1Modification, transaction2Modification)
	userModification.SetAssoc("address", addressModification)
	userModification.SetDeletedIDs("transactions", []interface{}{})

	assert.Equal(t, userModification, Apply(doc, data))
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
		userModification = Apply(NewDocument(&User{}),
			Set("name", "Luffy"),
			Set("age", 20),
		)
		transaction1Modification = Apply(NewDocument(&Transaction{}),
			Set("item", "Sword"),
		)
		transaction2Modification = Apply(NewDocument(&Transaction{}),
			Set("item", "Shield"),
		)
		addressModification = Apply(NewDocument(&Address{}),
			Set("street", "Grove Street"),
		)
	)

	userModification.SetAssoc("transactions", transaction1Modification, transaction2Modification)
	userModification.SetAssoc("address", addressModification)
	userModification.SetDeletedIDs("transactions", []interface{}{})

	assert.Equal(t, userModification, Apply(doc, data))
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
		userModification         = Apply(NewDocument(&User{}))
		transaction1Modification = Apply(NewDocument(&Transaction{}),
			Set("item", "Sword"),
		)
		transaction2Modification = Apply(NewDocument(&Transaction{}),
			Set("item", "Shield"),
		)
	)

	userModification.SetAssoc("transactions", transaction2Modification, transaction1Modification)
	userModification.SetDeletedIDs("transactions", []interface{}{2})

	assert.Equal(t, userModification, Apply(doc, data))
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

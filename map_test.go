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
			"id":   1,
			"name": "Luffy",
			"age":  20,
			"transactions": []Map{
				{"id": 1, "item": "Sword"},
				{"id": 2, "item": "Shield"},
			},
			"address": Map{
				"id":     1,
				"street": "Grove Street",
			},
		}
		userModification = Apply(NewDocument(&User{}),
			Set("id", 1),
			Set("name", "Luffy"),
			Set("age", 20),
		)
		transaction1Modification = Apply(NewDocument(&Transaction{}),
			Set("id", 1),
			Set("item", "Sword"),
		)
		transaction2Modification = Apply(NewDocument(&Transaction{}),
			Set("id", 2),
			Set("item", "Shield"),
		)
		addressModification = Apply(NewDocument(&Address{}),
			Set("id", 1),
			Set("street", "Grove Street"),
		)
	)

	userModification.SetAssoc("transactions", transaction1Modification, transaction2Modification)
	userModification.SetAssoc("address", addressModification)

	assertModification(t, userModification, Apply(doc, data))
	assert.Equal(t, User{
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
	}, user)
}

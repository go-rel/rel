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
		userChanges = ApplyChanges(NewDocument(&User{}),
			Set("id", 1),
			Set("name", "Luffy"),
			Set("age", 20),
		)
		transaction1Changes = ApplyChanges(NewDocument(&Transaction{}),
			Set("id", 1),
			Set("item", "Sword"),
		)
		transaction2Changes = ApplyChanges(NewDocument(&Transaction{}),
			Set("id", 2),
			Set("item", "Shield"),
		)
		addressChanges = ApplyChanges(NewDocument(&Address{}),
			Set("id", 1),
			Set("street", "Grove Street"),
		)
	)

	userChanges.SetAssoc("transactions", transaction1Changes, transaction2Changes)
	userChanges.SetAssoc("address", addressChanges)

	assertChanges(t, userChanges, ApplyChanges(doc, data))
	assert.Equal(t, User{
		ID:   1,
		Name: "Luffy",
		Age:  20,
		Transactions: []Transaction{
			// {ID: 1, Item: "Sword"}, FIXME
			{ID: 2, Item: "Shield"},
		},
		Address: Address{
			ID:     1,
			Street: "Grove Street",
		},
	}, user)
}

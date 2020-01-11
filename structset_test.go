package rel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertChanges(t *testing.T, ch1 Changes, ch2 Changes) {
	assert.Equal(t, len(ch1.fields), len(ch2.fields))
	assert.Equal(t, len(ch1.changes), len(ch2.changes))
	assert.Equal(t, len(ch1.assoc), len(ch2.assoc))
	assert.Equal(t, len(ch1.assocChanges), len(ch2.assocChanges))

	for field := range ch1.fields {
		assert.Equal(t, ch1.changes[ch1.fields[field]], ch2.changes[ch2.fields[field]])
	}

	for assoc := range ch1.assoc {
		var (
			ac1 = ch1.assocChanges[ch1.assoc[assoc]]
			ac2 = ch2.assocChanges[ch2.assoc[assoc]]
		)
		assert.Equal(t, ac1.StaleIDs, ac2.StaleIDs)
		assert.Equal(t, len(ac1.Changes), len(ac2.Changes))

		for i := range ac1.Changes {
			assertChanges(t, ac1.Changes[i], ac2.Changes[i])
		}
	}
}

func BenchmarkStructset(b *testing.B) {
	var (
		user = &User{
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
	)

	for n := 0; n < b.N; n++ {
		ApplyChanges(nil, NewStructset(user))
	}
}

func TestStructset(t *testing.T) {
	var (
		user = &User{
			ID:   1,
			Name: "Luffy",
			Age:  20,
		}
		changes, _ = ApplyChanges(NewDocument(user),
			Set("name", "Luffy"),
			Set("age", 20),
			Set("created_at", now()),
			Set("updated_at", now()),
		)
		result, err = ApplyChanges(nil, NewStructset(user))
	)

	assert.Nil(t, err)
	assertChanges(t, changes, result)
}

func TestStructset_withAssoc(t *testing.T) {
	var (
		user = &User{
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
		userChanges, _ = ApplyChanges(NewDocument(&User{}),
			Set("name", "Luffy"),
			Set("age", 20),
			Set("updated_at", now()),
		)
		transaction1Changes, _ = ApplyChanges(NewDocument(&Transaction{}),
			Set("item", "Sword"),
		)
		transaction2Changes, _ = ApplyChanges(NewDocument(&Transaction{}),
			Set("item", "Shield"),
		)
		addressChanges, _ = ApplyChanges(NewDocument(&Address{}),
			Set("street", "Grove Street"),
		)
	)

	userChanges.SetAssoc("transactions", transaction1Changes, transaction2Changes)
	userChanges.SetAssoc("address", addressChanges)

	result, err := ApplyChanges(nil, NewStructset(user))
	assert.Nil(t, err)
	assertChanges(t, userChanges, result)
}

func TestStructset_invalidCreatedAtType(t *testing.T) {
	type tmp struct {
		ID        int
		Name      string
		CreatedAt int
	}

	var (
		user = &tmp{
			Name:      "Luffy",
			CreatedAt: 1,
		}
		changes, _ = ApplyChanges(NewDocument(user),
			Set("name", "Luffy"),
			Set("created_at", 1),
		)
	)

	result, err := ApplyChanges(nil, NewStructset(user))
	assert.Nil(t, err)
	assertChanges(t, changes, result)
}

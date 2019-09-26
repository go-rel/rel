package grimoire

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertChanges(t *testing.T, ch1 Changes, ch2 Changes) {
	assert.Equal(t, len(ch1.Fields), len(ch2.Fields))
	assert.Equal(t, len(ch1.Changes), len(ch2.Changes))
	assert.Equal(t, len(ch1.Assoc), len(ch2.Assoc))
	assert.Equal(t, len(ch1.AssocChanges), len(ch2.AssocChanges))
	assert.Equal(t, ch1.constraints, ch2.constraints)

	for field := range ch1.Fields {
		assert.Equal(t, ch1.Changes[ch1.Fields[field]], ch2.Changes[ch2.Fields[field]])
	}

	for assoc := range ch1.Assoc {
		assert.Equal(t, ch1.AssocChanges[ch1.Assoc[assoc]], ch2.AssocChanges[ch2.Assoc[assoc]])
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
		BuildChanges(NewStructset(user))
	}
}

func TestStructset(t *testing.T) {
	var (
		user = &User{
			ID:        1,
			Name:      "Luffy",
			Age:       20,
			CreatedAt: time.Now(),
		}
		changes = BuildChanges(
			Set("name", "Luffy"),
			Set("age", 20),
			Set("created_at", user.CreatedAt),
		)
	)

	assertChanges(t, changes, BuildChanges(NewStructset(user)))
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
		userChanges = BuildChanges(
			Set("name", "Luffy"),
			Set("age", 20),
			Set("created_at", user.CreatedAt),
		)
		transaction1Changes = BuildChanges(
			Set("item", "Sword"),
		)
		transaction2Changes = BuildChanges(
			Set("item", "Shield"),
		)
		addressChanges = BuildChanges(
			Set("street", "Grove Street"),
		)
	)

	userChanges.SetAssoc("transactions", transaction1Changes, transaction2Changes)
	userChanges.SetAssoc("address", addressChanges)

	assertChanges(t, userChanges, BuildChanges(NewStructset(user)))
}

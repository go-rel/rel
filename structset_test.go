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
		assert.Equal(t, ch1.assocChanges[ch1.assoc[assoc]], ch2.assocChanges[ch2.assoc[assoc]])
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
			ID:   1,
			Name: "Luffy",
			Age:  20,
		}
		changes = BuildChanges(
			Set("name", "Luffy"),
			Set("age", 20),
			Set("created_at", now()),
			Set("updated_at", now()),
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
			Set("updated_at", now()),
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
		changes = BuildChanges(
			Set("name", "Luffy"),
			Set("created_at", 1),
		)
	)

	assertChanges(t, changes, BuildChanges(NewStructset(user)))
}

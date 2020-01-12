package rel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertModification(t *testing.T, ch1 Modification, ch2 Modification) {
	assert.Equal(t, len(ch1.fields), len(ch2.fields))
	assert.Equal(t, len(ch1.modification), len(ch2.modification))
	assert.Equal(t, len(ch1.assoc), len(ch2.assoc))
	assert.Equal(t, len(ch1.assocModification), len(ch2.assocModification))

	for field := range ch1.fields {
		assert.Equal(t, ch1.modification[ch1.fields[field]], ch2.modification[ch2.fields[field]])
	}

	for assoc := range ch1.assoc {
		var (
			ac1 = ch1.assocModification[ch1.assoc[assoc]]
			ac2 = ch2.assocModification[ch2.assoc[assoc]]
		)
		assert.Equal(t, len(ac1), len(ac2))

		for i := range ac1 {
			assertModification(t, ac1[i], ac2[i])
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
		Apply(nil, NewStructset(user))
	}
}

func TestStructset(t *testing.T) {
	var (
		user = &User{
			ID:   1,
			Name: "Luffy",
			Age:  20,
		}
		modification = Apply(NewDocument(user),
			Set("name", "Luffy"),
			Set("age", 20),
			Set("created_at", now()),
			Set("updated_at", now()),
		)
	)

	assertModification(t, modification, Apply(nil, NewStructset(user)))
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
		userModification = Apply(NewDocument(&User{}),
			Set("name", "Luffy"),
			Set("age", 20),
			Set("updated_at", now()),
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

	assertModification(t, userModification, Apply(nil, NewStructset(user)))
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
		modification = Apply(NewDocument(user),
			Set("name", "Luffy"),
			Set("created_at", 1),
		)
	)

	assertModification(t, modification, Apply(nil, NewStructset(user)))
}

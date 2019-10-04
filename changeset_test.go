package grimoire

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkChangeset(b *testing.B) {
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
		BuildChanges(NewChangeset(user))
	}
}

func TestChangeset(t *testing.T) {
	var (
		user = User{
			ID:        1,
			Name:      "Luffy",
			Age:       20,
			CreatedAt: time.Now(),
		}
		changeset = NewChangeset(&user)
	)

	// update
	user.Name = "Zoro"
	user.Age = 21

	assertChanges(t, BuildChanges(
		Set("name", "Zoro"),
		Set("age", 21),
	), BuildChanges(changeset))

	// update time
	user.CreatedAt = time.Date(2019, 9, 9, 16, 32, 0, 0, time.Local)

	assertChanges(t, BuildChanges(
		Set("name", "Zoro"),
		Set("age", 21),
		Set("created_at", user.CreatedAt),
	), BuildChanges(changeset))
}

func TestChangeset_withAssoc(t *testing.T) {
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
		changeset = NewChangeset(&user)
	)

	// update
	user.Age = 21
	user.Address.Street = "Jl. Lingkar"
	user.Transactions[1].Item = "Bow"

	fmt.Printf("%#v", BuildChanges(changeset))

	assertChanges(t, BuildChanges(Map{
		"age": 21,
		"address": Map{
			"street": "Jl. Lingkar",
		},
		"transactions": []Map{
			{
				"id":   2,
				"item": "Bow",
			},
		},
	}), BuildChanges(changeset))
}

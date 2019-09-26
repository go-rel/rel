package grimoire

import (
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

// func TestChangeset_withAssoc(t *testing.T) {
// 	var (
// 		user = User{
// 			ID:   1,
// 			Name: "Luffy",
// 			Age:  20,
// 			Transactions: []Transaction{
// 				{ID: 1, Item: "Sword"},
// 				{ID: 2, Item: "Shield"},
// 			},
// 			Address: Address{
// 				ID:     1,
// 				Street: "Grove Street",
// 			},
// 			CreatedAt: time.Now(),
// 		}
// 		changeset = NewChangeset(&user)
// 	)

// 	// update without assoc
// 	user.Age = 21

// 	assertChanges(t, BuildChanges(Map{
// 		"age": 21,
// 	}), BuildChanges(changeset))

// 	assertChanges(t, BuildChanges(Map{
// 		"name":       "Luffy",
// 		"age":        20,
// 		"created_at": user.CreatedAt,
// 		"transactions": []Map{
// 			{"item": "Sword"},
// 			{"item": "Shield"},
// 		},
// 		"address": Map{
// 			"street": "Grove Street",
// 		},
// 	}), BuildChanges(NewStructset(user)))
// }

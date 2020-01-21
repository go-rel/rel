package rel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
		doc = NewDocument(user)
	)

	for n := 0; n < b.N; n++ {
		Apply(doc, NewStructset(user))
	}
}

func TestStructset(t *testing.T) {
	var (
		user = &User{
			ID:   1,
			Name: "Luffy",
			Age:  20,
		}
		doc          = NewDocument(user)
		modification = Apply(NewDocument(&User{}),
			Set("name", "Luffy"),
			Set("age", 20),
			Set("created_at", now()),
			Set("updated_at", now()),
		)
	)

	assert.Equal(t, modification, Apply(doc, NewStructset(user)))
}

func TestStructset_withAssoc(t *testing.T) {
	var (
		createdAt = time.Now().Add(-time.Hour) // should retains
		user      = &User{
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
			CreatedAt: createdAt,
		}
		doc     = NewDocument(user)
		userMod = Apply(NewDocument(&User{}),
			Set("name", "Luffy"),
			Set("age", 20),
			Set("created_at", createdAt),
			Set("updated_at", now()),
		)
		trx1Mod = Apply(NewDocument(&Transaction{}),
			Set("item", "Sword"),
		)
		trx2Mod = Apply(NewDocument(&Transaction{}),
			Set("item", "Shield"),
		)
		addrMod = Apply(NewDocument(&Address{}),
			Set("street", "Grove Street"),
		)
	)

	userMod.SetAssoc("transactions", trx1Mod, trx2Mod)
	userMod.SetAssoc("address", addrMod)

	assert.Equal(t, userMod, Apply(doc, NewStructset(user)))
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
		doc          = NewDocument(user)
		modification = Apply(NewDocument(user),
			Set("name", "Luffy"),
			Set("created_at", 1),
		)
	)

	assert.Equal(t, modification, Apply(doc, NewStructset(user)))
}

func TestStructset_differentStruct(t *testing.T) {
	type UserTmp struct {
		ID   int
		Name string
		Age  int
	}

	var (
		usertmp UserTmp
		user    = &User{
			ID:   1,
			Name: "Luffy",
			Age:  20,
		}
		doc          = NewDocument(&usertmp)
		modification = Apply(NewDocument(user),
			Set("name", "Luffy"),
			Set("age", 20),
		)
	)

	assert.Equal(t, modification, Apply(doc, NewStructset(user)))
	assert.Equal(t, user.Name, usertmp.Name)
	assert.Equal(t, user.Age, usertmp.Age)
}

func TestStructset_differentStructMissingField(t *testing.T) {
	// missing age field.
	type UserTmp struct {
		ID   int
		Name string
	}

	var (
		user = &User{
			ID:   1,
			Name: "Luffy",
			Age:  20,
		}
		doc = NewDocument(&UserTmp{})
	)

	assert.Panics(t, func() {
		Apply(doc, NewStructset(user))
	})
}

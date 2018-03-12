package mysql

import (
	"fmt"
	"testing"

	. "github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID            uint
	State         string
	PaymentMethod string
	// CreatedAt time.Time
	// UpdatedAt time.Time
}

func TestQuery(t *testing.T) {
	adapter := Adapter{}
	adapter.Open("root@(127.0.0.1:3306)/papyrus_test")
	defer adapter.Close()
	qs, args := adapter.All(From("transactions AS t").Join("corporate_users AS c", Eq(I("t.corporate_id"), I("c.id"))).Limit(2))
	println(qs)

	users := []User{}
	err := adapter.Query(&users, qs, args)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))
	fmt.Printf("%v", users)
}

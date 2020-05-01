package rel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDirty(t *testing.T) {
	var (
		now  = time.Now()
		user = User{
			ID:        1,
			Name:      "User 1",
			Age:       20,
			UpdatedAt: now,
			CreatedAt: now,
		}
		snapshot = []interface{}{1, "User 1", 20, now, now}
		doc      = NewDocument(&user)
	)

	user.Init(doc)

	assert.Equal(t, snapshot, user.Dirty.snapshot)
	assert.Empty(t, user.Dirty.Changes())

	// update
	user.Name = "User 2"
	user.Age = 21

	assert.Equal(t, snapshot, user.Dirty.snapshot)
	assert.Equal(t, map[string][2]interface{}{
		"name": [2]interface{}{"User 1", "User 2"},
		"age":  [2]interface{}{20, 21},
	}, user.Dirty.Changes())
}

func TestDirty_ptr(t *testing.T) {
	var (
		dirty   Dirty
		userID  = 2
		address = Address{
			ID:     1,
			UserID: &userID,
		}
		snapshot = []interface{}{1, 2, "", Notes(""), nil}
		doc      = NewDocument(&address)
	)

	dirty.Init(doc)

	assert.Equal(t, snapshot, dirty.snapshot)
	assert.Empty(t, dirty.Changes())

	// update
	userID = 3
	address.UserID = &userID

	assert.Equal(t, snapshot, dirty.snapshot)
	assert.Equal(t, map[string][2]interface{}{
		"user_id": [2]interface{}{2, 3},
	}, dirty.Changes())
}

func TestDirty_belongsTo(t *testing.T) {
	var (
		dirty   Dirty
		address = Address{
			User: &User{
				ID:   1,
				Name: "User 1",
			},
		}
		snapshot = []interface{}{1, "User 1", 0, time.Time{}, time.Time{}}
		doc      = NewDocument(&address)
	)

	dirty.Init(doc)
	dirty = *dirty.assoc["user"]
	assert.Equal(t, dirty, address.User.Dirty)

	assert.Equal(t, snapshot, dirty.snapshot)
	assert.Empty(t, dirty.Changes())

	// update
	address.User.Name = "User Satu"

	assert.Equal(t, snapshot, dirty.snapshot)
	assert.Equal(t, map[string][2]interface{}{
		"name": [2]interface{}{"User 1", "User Satu"},
	}, dirty.Changes())
}

func TestDirty_hasOne(t *testing.T) {
	var (
		dirty Dirty
		user  = User{
			ID: 1,
			Address: Address{
				ID:     1,
				Street: "Grove Street",
				Notes:  "HQ",
			},
		}
		snapshot = []interface{}{1, nil, "Grove Street", Notes("HQ"), nil}
		doc      = NewDocument(&user)
	)

	dirty.Init(doc)
	dirty = *dirty.assoc["address"]

	assert.Equal(t, snapshot, dirty.snapshot)
	assert.Empty(t, dirty.Changes())

	// update
	user.Address.UserID = &user.ID
	user.Address.Street = "Grove Street Blvd"
	user.Address.Notes = Notes("Home")

	assert.Equal(t, snapshot, dirty.snapshot)
	assert.Equal(t, map[string][2]interface{}{
		"user_id": [2]interface{}{nil, user.ID},
		"street":  [2]interface{}{"Grove Street", "Grove Street Blvd"},
		"notes":   [2]interface{}{Notes("HQ"), Notes("Home")},
	}, dirty.Changes())
}

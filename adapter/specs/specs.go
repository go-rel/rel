package specs

import (
	"time"

	"github.com/Fs02/grimoire/c"
)

type User struct {
	ID        int64
	Name      string
	Gender    string
	Age       int
	Note      *string
	Addresses []Address
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Address struct {
	ID        int64
	UserID    int64
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User table identifiers
const (
	users     = "users"
	addresses = "addresses"
	id        = c.I("id")
	name      = c.I("name")
	gender    = c.I("gender")
	age       = c.I("age")
	note      = c.I("note")
	createdAt = c.I("created_at")
	address   = c.I("address")
)

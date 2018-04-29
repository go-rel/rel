// Package specs defines test specifications for grimoire's adapter.
package specs

import (
	"time"

	"github.com/Fs02/grimoire/c"
)

// User defines users schema.
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

// Address defines addresses schema.
type Address struct {
	ID        int64
	User      User
	UserID    *int64
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

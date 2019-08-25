// Package specs defines test specifications for grimoire's adapter.
package specs

import (
	"strings"
	"testing"
	"time"

	"github.com/Fs02/grimoire/adapter/sql"
	"github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
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

// Extra defines baz schema.
type Extra struct {
	ID    uint
	Slug  *string
	Score int
	UserID int
}

// User table identifiers
const (
	users     = "users"
	addresses = "addresses"
	extras    = "extras"
	id        = c.I("id")
	name      = c.I("name")
	gender    = c.I("gender")
	age       = c.I("age")
	note      = c.I("note")
	createdAt = c.I("created_at")
	address   = c.I("address")
)

var builder = sql.NewBuilder(&sql.Config{
	Placeholder: "?",
	EscapeChar:  "`",
})

func assertConstraint(t *testing.T, err error, ctype grimoire.ConstraintType, key string) {
	assert.NotNil(t, err)
	cerr, _ := err.(grimoire.ConstraintError)
	assert.True(t, strings.Contains(cerr.Key, key))
	assert.Equal(t, ctype, cerr.Type)
}

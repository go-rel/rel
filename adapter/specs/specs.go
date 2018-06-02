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

func assertConstraint(t *testing.T, err error, kind errors.Kind, field string) {
	assert.NotNil(t, err)
	gerr, _ := err.(errors.Error)
	assert.True(t, strings.Contains(gerr.Field, field))
	assert.Equal(t, kind, gerr.Kind())
}

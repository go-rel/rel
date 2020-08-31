// Package specs defines test specifications for rel's adapter.
package specs

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/sql"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

// Flag for configuration.
type Flag int

func (f Flag) enabled(flags []Flag) bool {
	return len(flags) > 0 && f&flags[0] == 0
}

const (
	// SkipDropColumn spec.
	SkipDropColumn Flag = 1 << iota
	// SkipRenameColumn spec.
	SkipRenameColumn
)

// User defines users schema.
type User struct {
	ID             int64
	Name           string
	Gender         string
	Age            int
	Note           *string
	Addresses      []Address
	PrimaryAddress *Address
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Address defines addresses schema.
type Address struct {
	ID        int64
	User      User
	UserID    *int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Extra defines extra schema.
type Extra struct {
	ID     uint
	Slug   *string
	Score  int
	UserID int64
}

// Composite primaries example.
type Composite struct {
	Primary1 int `db:",primary"`
	Primary2 int `db:",primary"`
	Data     string
}

var (
	config = sql.Config{
		Placeholder: "?",
		EscapeChar:  "`",
	}
	builder = sql.NewBuilder(config)
)

func assertConstraint(t *testing.T, err error, ctype rel.ConstraintType, key string) {
	assert.NotNil(t, err)
	cerr, ok := err.(rel.ConstraintError)
	assert.True(t, ok)
	assert.True(t, strings.Contains(cerr.Key, key))
	assert.Equal(t, ctype, cerr.Type)
}

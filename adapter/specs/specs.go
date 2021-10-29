// Package specs defines test specifications for rel's adapter.
package specs

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/adapter/sql"
	"github.com/stretchr/testify/assert"
)

var ctx = context.TODO()

// Flag for configuration.
type Flag int

func (f Flag) disabled(flags []Flag) bool {
	for i := range flags {
		if f&flags[i] == 0 {
			return true
		}
	}

	return false
}

const (
	// SkipDropColumn spec.
	SkipDropColumn Flag = 1 << iota
	// SkipRenameColumn spec.
	SkipRenameColumn
	// SkipAllAndAnyKeyword spec.
	SkipAllAndAnyKeyword
)

// User defines users schema.
type User struct {
	ID             int64
	Name           string
	Gender         string
	Age            int
	Note           *string
	Addresses      []Address `autosave:"true"`
	PrimaryAddress *Address  `autosave:"true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Address defines addresses schema.
type Address struct {
	ID        int64
	User      User `autosave:"true"`
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
)

func assertConstraint(t *testing.T, err error, ctype rel.ConstraintType, key string) {
	assert.NotNil(t, err)
	cerr, ok := err.(rel.ConstraintError)
	assert.True(t, ok)
	assert.True(t, strings.Contains(cerr.Key, key))
	assert.Equal(t, ctype, cerr.Type)
}

func waitForReplication() {
	// wait for replication
	time.Sleep(5 * time.Millisecond)
}

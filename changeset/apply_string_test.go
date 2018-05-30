package changeset

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyString(t *testing.T) {
	type User struct {
		Name string
	}

	user := User{}
	params := map[string]interface{}{
		"name": "¡¡¡Hello, Gophers!!!",
	}

	ch := Cast(user, params, []string{"name"})
	ApplyString(ch, "name", func(s string) string {
		return strings.TrimPrefix(s, "¡¡¡Hello, ")
	})

	assert.Equal(t, "Gophers!!!", ch.Changes()["name"])
}

func TestApplyString_ignored(t *testing.T) {
	type User struct {
		Name int
	}

	user := User{}
	params := map[string]interface{}{
		"name": 1,
	}

	ch := Cast(user, params, []string{"name"})
	ApplyString(ch, "name", func(s string) string {
		return strings.TrimPrefix(s, "¡¡¡Hello, ")
	})

	assert.Equal(t, 1, ch.Changes()["name"])
}

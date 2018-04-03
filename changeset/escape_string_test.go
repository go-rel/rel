package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {
	type User struct {
		Name string
	}

	user := User{}
	params := map[string]interface{}{
		"name": `"Fran & Freddie's Diner" <tasty@example.com>`,
	}

	ch := Cast(user, params, []string{"name"})
	EscapeString(ch, "name")

	assert.Equal(t, "&#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt;", ch.Changes()["name"])
}

package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnescapeString(t *testing.T) {
	type User struct {
		Name string
	}

	user := User{}
	params := map[string]interface{}{
		"name": `&quot;Fran &amp; Freddie&#39;s Diner&quot; &lt;tasty@example.com&gt;`,
	}

	ch := Cast(user, params, []string{"name"})
	UnescapeString(ch, "name")

	assert.Equal(t, `"Fran & Freddie's Diner" <tasty@example.com>`, ch.Changes()["name"])
}

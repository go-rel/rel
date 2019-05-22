package changeset

import (
	"testing"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {
	type User struct {
		Name     string
		Username string
	}

	user := User{}
	input := params.Map{
		"name":     `"Fran & Freddie's Diner" <tasty@example.com>`,
		"username": `"Fran & Freddie's Diner" <tasty@example.com>`,
	}

	ch := Cast(user, input, []string{"name", "username"})
	EscapeString(ch, "name", "username")

	assert.Equal(t, "&#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt;", ch.Changes()["name"])
	assert.Equal(t, "&#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt;", ch.Changes()["username"])
}

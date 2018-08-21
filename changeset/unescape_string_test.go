package changeset

import (
	"testing"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestUnescapeString(t *testing.T) {
	type User struct {
		Name string
	}

	user := User{}
	input := params.Map{
		"name": `&quot;Fran &amp; Freddie&#39;s Diner&quot; &lt;tasty@example.com&gt;`,
	}

	ch := Cast(user, input, []string{"name"})
	UnescapeString(ch, "name")

	assert.Equal(t, `"Fran & Freddie's Diner" <tasty@example.com>`, ch.Changes()["name"])
}

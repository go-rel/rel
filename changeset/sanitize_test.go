package changeset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field": `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
		},
	}

	Sanitize(ch, []string{"field"})

	expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
	assert.Equal(t, expected, ch.changes["field"])
}

package changeset

import (
	"testing"
)

func TestSanitize(t *testing.T) {
	ch := &Changeset{
		changes: map[string]interface{}{
			"field": `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
		},
	}

	Sanitize(ch, "field")

	expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
	if ch.changes["field"] != expected {
		t.Error("Expected", expected, "but got", ch.changes["field"])
	}
}

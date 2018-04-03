package changeset

import (
	"html"
)

// UnescapeString unescapes entities like "&lt;" to become "<". this is helper for html.UnescapeString.
func UnescapeString(ch *Changeset, field string) {
	ApplyString(ch, field, html.UnescapeString)
}

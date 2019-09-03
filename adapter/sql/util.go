package sql

import (
	"strings"
)

// ExtractString between two string.
// TODO: use strings.Index
func ExtractString(s, left, right string) string {
	parts := strings.Split(s, left)
	if len(parts) <= 1 {
		return s
	}

	return strings.Split(parts[1], right)[0]
}

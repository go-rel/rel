package internal

import (
	"strings"
)

// ExtractString between two string.
func ExtractString(s, left, right string) string {
	return strings.Split(strings.Split(s, left)[1], right)[0]
}

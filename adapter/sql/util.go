package sql

import (
	"strings"
)

// ExtractString between two string.
func ExtractString(s, left, right string) string {
	var (
		start = strings.Index(s, left)
		end   = strings.LastIndex(s, right)
	)

	if start < 0 || end < 0 || start+len(left) >= end {
		return s
	}

	return s[start+len(left) : end]
}

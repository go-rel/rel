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

func toInt64(i interface{}) int64 {
	var result int64

	switch s := i.(type) {
	case int:
		result = int64(s)
	case int64:
		result = s
	case int32:
		result = int64(s)
	case int16:
		result = int64(s)
	case int8:
		result = int64(s)
	case uint:
		result = int64(s)
	case uint64:
		result = int64(s)
	case uint32:
		result = int64(s)
	case uint16:
		result = int64(s)
	case uint8:
		result = int64(s)
	}

	return result
}

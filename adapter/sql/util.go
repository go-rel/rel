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

type fieldCacheKey struct {
	field  string
	escape string
}

// Escape field or table name.
func Escape(config Config, field string) string {
	if config.EscapeChar == "" || field == "*" {
		return field
	}

	key := fieldCacheKey{field: field, escape: config.EscapeChar}
	escapedField, ok := fieldCache.Load(key)
	if ok {
		return escapedField.(string)
	}

	if len(field) > 0 && field[0] == UnescapeCharacter {
		escapedField = field[1:]
	} else if i := strings.Index(strings.ToLower(field), " as "); i > -1 {
		escapedField = Escape(config, field[:i]) + " AS " + Escape(config, field[i+4:])
	} else if start, end := strings.IndexRune(field, '('), strings.IndexRune(field, ')'); start >= 0 && end >= 0 && end > start {
		escapedField = field[:start+1] + Escape(config, field[start+1:end]) + field[end:]
	} else if strings.HasSuffix(field, "*") {
		escapedField = config.EscapeChar + strings.Replace(field, ".", config.EscapeChar+".", 1)
	} else {
		escapedField = config.EscapeChar +
			strings.Replace(field, ".", config.EscapeChar+"."+config.EscapeChar, 1) +
			config.EscapeChar
	}

	fieldCache.Store(key, escapedField)
	return escapedField.(string)
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

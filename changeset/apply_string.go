package changeset

import (
	"reflect"
)

// ApplyString apply a function for string value.
func ApplyString(ch *Changeset, field string, fn func(string) string) {
	if val, ok := ch.changes[field]; ok && val != nil && ch.types[field].Kind() == reflect.String {
		ch.changes[field] = fn(val.(string))
	}
}

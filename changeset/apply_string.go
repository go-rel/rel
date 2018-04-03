package changeset

import (
	"reflect"
)

// ApplyString apply a function for string value.
func ApplyString(ch *Changeset, field string, fn func(string) string) {
	if ch.types[field].Kind() == reflect.String {
		ch.changes[field] = fn(ch.changes[field].(string))
	}
}

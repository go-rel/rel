package changeset

import (
	"reflect"
	"strings"
)

// PutDefaultErrorMessage is the default error message for PutDefault.
var PutDefaultErrorMessage = "{field} is invalid"

// PutDefault to changeset.
func PutDefault(ch *Changeset, field string, value interface{}, opts ...Option) {
	options := Options{
		message: PutDefaultErrorMessage,
	}
	options.apply(opts)

	if typ, exist := ch.types[field]; exist {
		valTyp := reflect.TypeOf(value)
		if valTyp.ConvertibleTo(typ) {
			if ch.changes[field] == nil {
				ch.changes[field] = value
			}
			return
		}
	}

	msg := strings.Replace(options.message, "{field}", field, 1)
	AddError(ch, field, msg)
}

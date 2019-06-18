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
			if (ch.params == nil || !ch.params.Exists(field)) && // no input
				ch.changes[field] == nil && // no change
				isZero(ch.values[field]) { // existing value is zero value
				ch.changes[field] = value
			}
			return
		}
	}

	msg := strings.Replace(options.message, "{field}", field, 1)
	AddError(ch, field, msg)
}

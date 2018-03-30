package changeset

import (
	"reflect"
	"strings"
)

var PutChangeErrorMessage = "{field} is invalid"

func PutChange(ch *Changeset, field string, value interface{}, opts ...Option) {
	options := Options{
		Message:    PutChangeErrorMessage,
		AllowError: false,
	}
	options.Apply(opts)

	if typ, exist := ch.types[field]; exist && reflect.TypeOf(value).ConvertibleTo(typ) {
		ch.changes[field] = value
	} else if !options.AllowError {
		msg := strings.Replace(options.Message, "{field}", field, 1)
		AddError(ch, field, msg)
	}
}

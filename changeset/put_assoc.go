package changeset

import (
	"reflect"
	"strings"
)

// PutAssocErrorMessage is the default error message for PutAssoc.
var PutAssocErrorMessage = "{field} is invalid"

// PutAssoc to changeset.
func PutAssoc(ch *Changeset, field string, value interface{}, opts ...Option) {
	options := Options{
		message: PutAssocErrorMessage,
	}
	options.apply(opts)
	if typ, exist := ch.types[field]; exist {
		if typ.Kind() == reflect.Struct {
			valTyp := reflect.TypeOf(value)
			if valTyp == reflect.TypeOf(ch) {
				ch.changes[field] = value
				return
			}
		} else if typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Struct {
			valTyp := reflect.TypeOf(value)
			if valTyp.Kind() == reflect.Slice && valTyp.Elem().ConvertibleTo(reflect.TypeOf(ch)) {
				ch.changes[field] = value
				return
			}
		}
	}
	msg := strings.Replace(options.message, "{field}", field, 1)
	AddError(ch, field, msg)
}

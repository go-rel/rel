package changeset

import (
	"reflect"
	"strings"

	"github.com/Fs02/grimoire/errors"
)

var CastOneErrorMessage = "{field} is invalid"

func CastOne(ch *Changeset, field string, fn func(interface{}, map[string]interface{}) *Changeset, opts ...Option) {
	options := Options{
		Message: CastOneErrorMessage,
	}
	options.Apply(opts)

	if par, exist := ch.params[field]; exist {
		mpar, isMap := par.(map[string]interface{})
		if typ, exist := ch.types[field]; exist && isMap && typ.Kind() == reflect.Struct {
			var onech *Changeset

			if val, exist := ch.values[field]; exist && val != nil {
				onech = fn(val, mpar)
			} else {
				onech = fn(reflect.Zero(typ).Interface(), mpar)
			}

			ch.changes[field] = onech

			// add errors to main errors
			for _, err := range onech.errors {
				e := err.(errors.Error)
				AddError(ch, field+"."+e.Field, e.Message)
			}
		} else {
			msg := strings.Replace(options.Message, "{field}", field, 1)
			AddError(ch, field, msg)
		}
	}
}

package changeset

import (
	"reflect"
	"strings"
)

var CastErrorMessage = "{field} is invalid"

func Cast(entity interface{}, params map[string]interface{}, fields []string, opts ...Option) *Changeset {
	options := Options{
		Message: CastErrorMessage,
	}
	options.Apply(opts)

	ch := &Changeset{}
	ch.changes = make(map[string]interface{})
	ch.schema = mapFields(entity)

	for _, f := range fields {
		val, pexist := params[f]
		sf, sexist := ch.schema[f]
		if pexist && sexist {
			if reflect.TypeOf(val).ConvertibleTo(sf.Type) {
				ch.changes[f] = val
			} else {
				msg := strings.Replace(options.Message, "{field}", f, 1)
				AddError(ch, f, msg)
			}
		}
	}

	return ch
}

package changeset

import (
	"reflect"
	"strings"

	"github.com/Fs02/grimoire/params"
	"github.com/Fs02/grimoire/schema"
)

// CastErrorMessage is the default error message for Cast.
var CastErrorMessage = "{field} is invalid"

// Cast params as changes for the given data according to the permitted fields. Returns a new changeset.
// params will only be added as changes if it does not have the same value as the field in the data.
func Cast(data interface{}, params params.Params, fields []string, opts ...Option) *Changeset {
	options := Options{
		message:     CastErrorMessage,
		emptyValues: []interface{}{""},
	}
	options.apply(opts)

	var ch *Changeset
	if existingCh, ok := data.(Changeset); ok {
		ch = &existingCh
	} else if existingCh, ok := data.(*Changeset); ok {
		ch = existingCh
	} else {
		ch = &Changeset{}
		ch.params = params
		ch.changes = make(map[string]interface{})
		ch.values, ch.types, ch.zero = mapSchema(data, true)
	}

	for _, field := range fields {
		typ, texist := ch.types[field]

		if !params.Exists(field) || !texist {
			continue
		}

		// ignore if it's an empty value
		if contains(options.emptyValues, params.Get(field)) {
			continue
		}

		if change, valid := params.GetWithType(field, typ); valid {
			value, vexist := ch.values[field]

			if (typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array) || (ch.zero && change != nil) || (!vexist && change != nil) || (vexist && value != change) {
				ch.changes[field] = change
			}
		} else {
			msg := strings.Replace(options.message, "{field}", field, 1)
			AddError(ch, field, msg)
		}
	}

	return ch
}

func mapSchema(data interface{}, zero bool) (map[string]interface{}, map[string]reflect.Type, bool) {
	var (
		fields    = schema.InferFields(data)
		types     = schema.InferTypes(data)
		values    = schema.InferValues(data)
		valuesMap = make(map[string]interface{}, len(fields))
		typesMap  = make(map[string]reflect.Type, len(fields))
	)

	for field, index := range fields {
		typesMap[field] = types[index]

		value := values[index]

		if value != nil {
			valuesMap[field] = values[index]

			if zero {
				zero = isZero(value)
			}
		}
	}

	return valuesMap, typesMap, zero
}

func contains(vs []interface{}, v interface{}) bool {
	for i := range vs {
		if vs[i] == v {
			return true
		}
	}

	return false
}

// isZero shallowly check wether a field in struct is zero or not
func isZero(i interface{}) bool {
	zero := true

	switch v := i.(type) {
	case bool:
		zero = v == false
	case string:
		zero = v == ""
	case int:
		zero = v == 0
	case int8:
		zero = v == 0
	case int16:
		zero = v == 0
	case int32:
		zero = v == 0
	case int64:
		zero = v == 0
	case uint:
		zero = v == 0
	case uint8:
		zero = v == 0
	case uint16:
		zero = v == 0
	case uint32:
		zero = v == 0
	case uint64:
		zero = v == 0
	case uintptr:
		zero = v == 0
	case float32:
		zero = v == 0
	case float64:
		zero = v == 0
	}

	return zero
}

package changeset

import (
	"reflect"
	"strings"

	"github.com/Fs02/grimoire/params"
	"github.com/azer/snakecase"
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
		ch.values, ch.types, ch.zero = mapSchema(data)
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

func mapSchema(data interface{}) (map[string]interface{}, map[string]reflect.Type, bool) {
	mvalues := make(map[string]interface{})
	mtypes := make(map[string]reflect.Type)
	zero := true

	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()

	if rv.Kind() != reflect.Struct {
		panic("data must be a struct")
	}

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		var name string
		if tag := ft.Tag.Get("db"); tag != "" {
			if tag == "-" {
				continue
			}
			name = tag
		} else {
			name = snakecase.SnakeCase(ft.Name)
		}

		if ft.Type.Kind() == reflect.Ptr {
			mtypes[name] = ft.Type.Elem()
			if !fv.IsNil() {
				mvalues[name] = fv.Elem().Interface()
			}
		} else if ft.Type.Kind() == reflect.Slice && ft.Type.Elem().Kind() == reflect.Ptr {
			mtypes[name] = reflect.SliceOf(ft.Type.Elem().Elem())
			mvalues[name] = fv.Interface()
		} else {
			mtypes[name] = fv.Type()
			mvalues[name] = fv.Interface()
		}

		if zero {
			zero = isZero(fv)
		}
	}

	return mvalues, mtypes, zero
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
func isZero(rv reflect.Value) bool {
	zero := true
	switch rv.Kind() {
	case reflect.Bool:
		zero = !rv.Bool()
	case reflect.Float32, reflect.Float64:
		zero = rv.Float() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		zero = rv.Int() == 0
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		zero = rv.IsNil()
	case reflect.String:
		zero = rv.String() == ""
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		zero = rv.Uint() == 0
	}

	return zero
}

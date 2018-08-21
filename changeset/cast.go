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
		message: CastErrorMessage,
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
		ch.values, ch.types = mapSchema(data)
	}

	for _, field := range fields {
		typ, texist := ch.types[field]
		currentValue, vexist := ch.values[field]

		if !params.Exists(field) || !texist {
			continue
		}

		if value, valid := params.GetWithType(field, typ); valid {
			if (typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array) || (!vexist && value != nil) || (vexist && currentValue != value) {
				ch.changes[field] = value
			}
		} else {
			msg := strings.Replace(options.message, "{field}", field, 1)
			AddError(ch, field, msg)
		}
	}

	return ch
}

func mapSchema(data interface{}) (map[string]interface{}, map[string]reflect.Type) {
	mvalues := make(map[string]interface{})
	mtypes := make(map[string]reflect.Type)

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
	}

	return mvalues, mtypes
}

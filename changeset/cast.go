package changeset

import (
	"reflect"
	"strings"

	"github.com/azer/snakecase"
)

var CastErrorMessage = "{field} is invalid"

func Cast(entity interface{}, params map[string]interface{}, fields []string, opts ...Option) *Changeset {
	options := Options{
		Message: CastErrorMessage,
	}
	options.Apply(opts)

	ch := &Changeset{}
	ch.entity = entity
	ch.params = params
	ch.changes = make(map[string]interface{})
	ch.values, ch.types = mapSchema(ch.entity)

	for _, f := range fields {
		val, pexist := params[f]
		typ, texist := ch.types[f]
		if pexist && texist {
			if reflect.TypeOf(val).ConvertibleTo(typ) {
				ch.changes[f] = val
			} else {
				msg := strings.Replace(options.Message, "{field}", f, 1)
				AddError(ch, f, msg)
			}
		}
	}

	return ch
}

func mapSchema(entity interface{}) (map[string]interface{}, map[string]reflect.Type) {
	mvalues := make(map[string]interface{})
	mtypes := make(map[string]reflect.Type)

	rv := reflect.ValueOf(entity)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()

	if rv.Kind() != reflect.Struct {
		panic("entity must be a struct")
	}

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)

		var name string
		if tag := ft.Tag.Get("db"); tag != "" && tag != "-" {
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

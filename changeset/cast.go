package changeset

import (
	"reflect"
	"strings"

	"github.com/azer/snakecase"
)

// CastErrorMessage is the default error message for Cast.
var CastErrorMessage = "{field} is invalid"

// Cast params as changes for the given entity according to the given fields. Returns a new changeset.
func Cast(entity interface{}, params map[string]interface{}, fields []string, opts ...Option) *Changeset {
	options := Options{
		message: CastErrorMessage,
	}
	options.apply(opts)

	ch := &Changeset{}
	ch.entity = entity
	ch.params = params
	ch.changes = make(map[string]interface{})
	ch.values, ch.types = mapSchema(ch.entity)

	for _, f := range fields {
		val, pexist := params[f]
		typ, texist := ch.types[f]
		if pexist && texist {
			if val != (interface{})(nil) {
				valTyp := reflect.TypeOf(val)
				if valTyp.Kind() == reflect.Ptr {
					valTyp = valTyp.Elem()
				}

				if valTyp.ConvertibleTo(typ) {
					ch.changes[f] = val
					continue
				}
			} else {
				ch.changes[f] = reflect.Zero(reflect.PtrTo(typ)).Interface()
				continue
			}

			msg := strings.Replace(options.message, "{field}", f, 1)
			AddError(ch, f, msg)
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

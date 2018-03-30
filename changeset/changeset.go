package changeset

import (
	"reflect"

	"github.com/azer/snakecase"
)

type Changeset struct {
	errors  []error
	schema  map[string]Field
	changes map[string]interface{}
}

func (changeset *Changeset) Changes() map[string]interface{} {
	return changeset.changes
}

func (changeset *Changeset) Errors() []error {
	return changeset.errors
}

func (changeset *Changeset) Error() error {
	if changeset.errors != nil {
		return changeset.errors[0]
	}
	return nil
}

type Field struct {
	Value interface{}
	Type  reflect.Type
}

func mapFields(entity interface{}) map[string]Field {
	mf := make(map[string]Field)

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
		def := Field{}

		if fv.Kind() == reflect.Ptr {
			def.Type = ft.Type.Elem()
			if !fv.IsNil() {
				def.Value = fv.Elem().Interface()
			}
		} else {
			def.Type = fv.Type()
			def.Value = fv.Interface()
		}

		if tag := ft.Tag.Get("db"); tag != "" && tag != "-" {
			mf[tag] = def
		} else {
			mf[snakecase.SnakeCase(ft.Name)] = def
		}
	}

	return mf
}

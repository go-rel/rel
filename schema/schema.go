package schema

import (
	"reflect"
	"strings"
	"sync"

	"github.com/azer/snakecase"
)

type schema interface {
	Fields() map[string]int
	Types() []reflect.Type
}

var schemasCache sync.Map

type schemaData struct {
	fields map[string]int
	types  []reflect.Type
}

func InferFields(record interface{}) map[string]int {
	if s, ok := record.(schema); ok {
		return s.Fields()
	}

	fields, _ := inferFieldAndTypes(record)
	return fields
}

func InferTypes(record interface{}) []reflect.Type {
	if s, ok := record.(schema); ok {
		return s.Types()
	}

	_, types := inferFieldAndTypes(record)
	return types
}

func inferFieldAndTypes(record interface{}) (map[string]int, []reflect.Type) {
	rt := reflectTypePtr(record)

	// check for cache
	if v, cached := schemasCache.Load((rt)); cached {
		s := v.(schemaData)
		return s.fields, s.types
	}

	var (
		rv     = reflectValuePtr(record)
		index  = 0
		fields = make(map[string]int, rv.NumField())
		types  = make([]reflect.Type, 0, rv.NumField())
	)

	for i := 0; i < rv.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			ft   = sf.Type
			name = inferFieldName(sf)
		)

		if name == "" {
			continue
		}

		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		} else if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Ptr {
			ft = reflect.SliceOf(ft.Elem().Elem())
		}

		fields[name] = index
		types = append(types, ft)
		index++
	}

	schemasCache.Store(rt, schemaData{
		fields: fields,
		types:  types,
	})

	return fields, types
}

func inferFieldName(sf reflect.StructField) string {
	if tag := sf.Tag.Get("db"); tag != "" {
		if tag == "-" {
			return ""
		}

		return strings.Split(tag, ",")[0]
	}

	return snakecase.SnakeCase(sf.Name)
}

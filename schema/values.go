package schema

import (
	"reflect"
	"sync"
)

type values interface {
	Values() []interface{}
}

func InferValues(record interface{}) []interface{} {
	if v, ok := record.(values); ok {
		return v.Values()
	}

	var (
		rv      = reflectValuePtr(record)
		fields  = InferFields(record)
		mapping = inferFieldMapping(record)
		values  = make([]interface{}, len(fields))
	)

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			fv          = rv.Field(structIndex)
			ft          = fv.Type()
		)

		if ft.Kind() == reflect.Ptr {
			if !fv.IsNil() {
				values[index] = fv.Elem().Interface()
			}
		} else {
			values[index] = fv.Interface()
		}
	}

	return values
}

var fieldMappingCache sync.Map

func inferFieldMapping(record interface{}) map[string]int {
	rt := reflectTypePtr(record)

	// check for cache
	if v, cached := fieldMappingCache.Load((rt)); cached {
		return v.(map[string]int)
	}

	var (
		rv      = reflectValuePtr(record)
		mapping = make(map[string]int, rv.NumField())
	)

	for i := 0; i < rv.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			name = inferFieldName(sf)
		)

		if name == "" {
			continue
		}

		mapping[name] = i
	}

	fieldMappingCache.Store(rt, mapping)

	return mapping
}

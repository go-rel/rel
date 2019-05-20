package schema

import (
	"reflect"
)

type values interface {
	fields
	Values() []interface{}
}

// InferValues from struct.
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

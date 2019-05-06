package schema

import (
	"reflect"
)

type scanners interface {
	fields
	Scanners() []interface{}
}

func InferScanners(record interface{}) []interface{} {
	if v, ok := record.(scanners); ok {
		return v.Scanners()
	}

	var (
		rv      = reflectValuePtr(record)
		fields  = InferFields(record)
		mapping = inferFieldMapping(record)
		values  = make([]interface{}, len(fields))
	)

	if !rv.CanAddr() {
		panic("grimoire: not a pointer")
	}

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			fv          = rv.Field(structIndex)
			ft          = fv.Type()
		)

		if ft.Kind() == reflect.Ptr {
			values[index] = fv.Addr().Interface()
		} else {
			values[index] = Value(fv.Addr().Interface())
		}
	}

	return values
}

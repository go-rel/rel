package schema

import (
	"database/sql"
	"reflect"
)

type scanners interface {
	Scanners([]string) []interface{}
}

// InferScanners from struct.
func InferScanners(record interface{}, fields []string) []interface{} {
	if v, ok := record.(scanners); ok {
		return v.Scanners(fields)
	}

	if s, ok := record.(sql.Scanner); ok {
		return []interface{}{s}
	}

	var (
		rv        = reflectValuePtr(record)
		mapping   = inferFieldMapping(record)
		result    = make([]interface{}, len(fields))
		tempValue = sql.RawBytes{}
	)

	if !rv.CanAddr() {
		panic("grimoire: not a pointer")
	}

	for index, field := range fields {
		if structIndex, ok := mapping[field]; ok {
			var (
				fv = rv.Field(structIndex)
				ft = fv.Type()
			)

			if ft.Kind() == reflect.Ptr {
				result[index] = fv.Addr().Interface()
			} else {
				result[index] = Nullable(fv.Addr().Interface())
			}
		} else {
			result[index] = &tempValue
		}
	}

	return result
}

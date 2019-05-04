package internal

import (
	"reflect"
)

func reflectInternalType(record interface{}) reflect.Type {
	rt := reflect.TypeOf(record)

	for rt.Kind() == reflect.Ptr || rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
		rt = rt.Elem()
	}

	return rt
}

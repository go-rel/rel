package schema

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

func reflectTypePtr(record interface{}) reflect.Type {
	rt := reflect.TypeOf(record)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		panic("grimoire: record must be a struct")
	}

	return rt
}

func reflectValuePtr(record interface{}) reflect.Value {
	rv := reflect.ValueOf(record)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		panic("grimoire: record must be a struct")
	}

	return rv
}

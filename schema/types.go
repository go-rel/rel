package schema

import (
	"reflect"
	"sync"
	"time"
)

var (
	// Bool is alias for reflect.TypeOf(false)
	Bool = reflect.TypeOf(false)
	// String is alias reflect.TypeOf("")
	String = reflect.TypeOf("")
	// Int is alias for reflect.TypeOf(int(0))
	Int = reflect.TypeOf(int(0))
	// Int8 is alias for reflect.TypeOf(int8(0))
	Int8 = reflect.TypeOf(int8(0))
	// Int16 is alias for reflect.TypeOf(int16(0))
	Int16 = reflect.TypeOf(int16(0))
	// Int32 is alias for reflect.TypeOf(int32(0))
	Int32 = reflect.TypeOf(int32(0))
	// Int64 is alias for reflect.TypeOf(int64(0))
	Int64 = reflect.TypeOf(int64(0))
	// Uint is alias for reflect.TypeOf(uint(0))
	Uint = reflect.TypeOf(uint(0))
	// Uint8 is alias for reflect.TypeOf(uint8(0))
	Uint8 = reflect.TypeOf(uint8(0))
	// Uint16 is alias for reflect.TypeOf(uint16(0))
	Uint16 = reflect.TypeOf(uint16(0))
	// Uint32 is alias for reflect.TypeOf(uint32(0))
	Uint32 = reflect.TypeOf(uint32(0))
	// Uint64 is alias for reflect.TypeOf(uint64(0))
	Uint64 = reflect.TypeOf(uint64(0))
	// Uintptr is alias for reflect.TypeOf(uintptr(0))
	Uintptr = reflect.TypeOf(uintptr(0))
	// Byte is alias for reflect.TypeOf(byte(0))
	Byte = reflect.TypeOf(byte(0))
	// Rune is alias for reflect.TypeOf(rune(' '))
	Rune = reflect.TypeOf(rune(' '))
	// Float32 is alias for reflect.TypeOf(float32(0))
	Float32 = reflect.TypeOf(float32(0))
	// Float64 is alias for reflect.TypeOf(float64(0))
	Float64 = reflect.TypeOf(float64(0))
	// Bytes is alias for reflect.TypeOf([]byte{})
	Bytes = reflect.TypeOf([]byte{})
	// Time is alias for reflect.TypeOf(time.Time{})
	Time = reflect.TypeOf(time.Time{})
)

type types interface {
	fields
	Types() []reflect.Type
}

var typesCache sync.Map

// InferTypes from struct.
func InferTypes(record interface{}) []reflect.Type {
	if v, ok := record.(types); ok {
		return v.Types()
	}

	rt := reflectTypePtr(record)

	// check for cache
	if v, cached := typesCache.Load(rt); cached {
		return v.([]reflect.Type)
	}

	var (
		fields  = InferFields(record)
		mapping = inferFieldMapping(record)
		types   = make([]reflect.Type, len(fields))
	)

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			ft          = rt.Field(structIndex).Type
		)

		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		} else if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Ptr {
			ft = reflect.SliceOf(ft.Elem().Elem())
		}

		types[index] = ft
	}

	typesCache.Store(rt, types)

	return types
}

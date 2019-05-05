package internal

import (
	"reflect"
)

var (
	Bool    = reflect.TypeOf(false)
	String  = reflect.TypeOf("")
	Int     = reflect.TypeOf(int(0))
	Int8    = reflect.TypeOf(int8(0))
	Int16   = reflect.TypeOf(int16(0))
	Int32   = reflect.TypeOf(int32(0))
	Int64   = reflect.TypeOf(int64(0))
	Uint    = reflect.TypeOf(uint(0))
	Uint8   = reflect.TypeOf(uint8(0))
	Uint16  = reflect.TypeOf(uint16(0))
	Uint32  = reflect.TypeOf(uint32(0))
	Uint64  = reflect.TypeOf(uint64(0))
	Uintptr = reflect.TypeOf(uintptr(0))
	Byte    = reflect.TypeOf(byte(0))
	Rune    = reflect.TypeOf(rune(' '))
	Float32 = reflect.TypeOf(float32(0))
	Float64 = reflect.TypeOf(float64(0))
)

package rel

import (
	"math"
	"reflect"
)

func indirect(rv reflect.Value) interface{} {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}

		rv = rv.Elem()
	}

	return rv.Interface()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type isZeroer interface {
	IsZero() bool
}

// isZero shallowly check wether a field in struct is zero or not
func isZero(value interface{}) bool {
	var (
		zero bool
	)

	switch v := value.(type) {
	case nil:
		zero = true
	case bool:
		zero = v == false
	case string:
		zero = v == ""
	case int:
		zero = v == 0
	case int8:
		zero = v == 0
	case int16:
		zero = v == 0
	case int32:
		zero = v == 0
	case int64:
		zero = v == 0
	case uint:
		zero = v == 0
	case uint8:
		zero = v == 0
	case uint16:
		zero = v == 0
	case uint32:
		zero = v == 0
	case uint64:
		zero = v == 0
	case uintptr:
		zero = v == 0
	case float32:
		zero = v == 0
	case float64:
		zero = v == 0
	case isZeroer:
		zero = v.IsZero()
	default:
		zero = isDeepZero(reflect.ValueOf(value), 0)
	}

	return zero
}

// modified from https://golang.org/src/reflect/value.go?s=33807:33835#L1077
func isDeepZero(rv reflect.Value, depth int) bool {
	if depth < 0 {
		return true
	}

	switch rv.Kind() {
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(rv.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := rv.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if !isDeepZero(rv.Index(i), depth-1) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		return rv.IsNil()
	case reflect.Slice:
		return rv.IsNil() || rv.Len() == 0
	case reflect.String:
		return rv.Len() == 0
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			if !isDeepZero(rv.Field(i), depth-1) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

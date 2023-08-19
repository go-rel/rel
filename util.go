package rel

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func indirectInterface(rv reflect.Value) any {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}

		rv = rv.Elem()
	}

	return rv.Interface()
}

func indirectReflectType(rt reflect.Type) reflect.Type {
	if rt.Kind() == reflect.Ptr {
		return rt.Elem()
	}

	return rt
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustTrue(flag bool, msg string) {
	if !flag {
		panic(msg)
	}
}

type isZeroer interface {
	IsZero() bool
}

// isZero shallowly check wether a field in struct is zero or not
func isZero(value any) bool {
	var (
		zero bool
	)

	switch v := value.(type) {
	case nil:
		zero = true
	case bool:
		zero = !v
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
		// check one level deeper if it's an uuid ([16]byte)
		if rv.Type().Elem().Kind() == reflect.Uint8 && rv.Len() == 16 {
			depth += 1
		}

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

func setPointerValue(ft reflect.Type, fv reflect.Value, rt reflect.Type, rv reflect.Value) bool {
	if ft.Elem() != rt && !rt.AssignableTo(ft.Elem()) {
		return false
	}

	if fv.IsNil() {
		fv.Set(reflect.New(ft.Elem()))
	}
	fv.Elem().Set(rv)

	return true
}

func setConvertValue(ft reflect.Type, fv reflect.Value, rt reflect.Type, rv reflect.Value) bool {
	var (
		rk = rt.Kind()
		fk = ft.Kind()
	)

	// prevents unintentional conversion
	if (rk >= reflect.Int || rk <= reflect.Uint64) && fk == reflect.String {
		return false
	}

	fv.Set(rv.Convert(ft))
	return true
}

func fmtAny(v any) string {
	if str, ok := v.(string); ok {
		return "\"" + str + "\""
	}

	return fmt.Sprint(v)
}

func fmtAnys(v []any) string {
	var str strings.Builder
	for i := range v {
		if i > 0 {
			str.WriteString(", ")
		}
		str.WriteString(fmtAny(v[i]))
	}

	return str.String()
}

// Encode index slice into single string
func encodeIndices(indices []int) string {
	var sb strings.Builder
	for _, index := range indices {
		sb.WriteString("/")
		sb.WriteString(strconv.Itoa(index))
	}
	return sb.String()
}

// Get field by index and init pointers on path if flag is true
//
//	modified from: https://cs.opensource.google/go/go/+/refs/tags/go1.17.7:src/reflect/value.go;l=1228-1245;bpv
func reflectValueFieldByIndex(rv reflect.Value, index []int, init bool) reflect.Value {
	if len(index) == 1 {
		return rv.Field(index[0])
	}

	for depth := 0; depth < len(index)-1; depth += 1 {
		field := rv.Field(index[depth])

		if field.Kind() != reflect.Ptr {
			rv = field
			continue
		}

		if field.IsNil() {
			if !init {
				targetType := field.Type().Elem().FieldByIndex(index[depth+1:]).Type
				return reflect.Zero(reflect.PtrTo(targetType))
			}
			field.Set(reflect.New(field.Type().Elem()))
		}

		rv = field.Elem()
	}
	return rv.Field(index[len(index)-1])
}

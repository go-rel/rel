package grimoire

import (
	"reflect"
)

func transformError(err error) error {
	// TODO
	return err
}

func indirect(rv reflect.Value) interface{} {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}

		rv = rv.Elem()
	}

	return rv.Interface()
}

// must is grimoire version of paranoid.Panic without context, but only original error.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

type isZeroer interface {
	IsZero() bool
}

// isZero shallowly check wether a field in struct is zero or not
func isZero(i interface{}) bool {
	zero := true

	switch v := i.(type) {
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
	}

	return zero
}

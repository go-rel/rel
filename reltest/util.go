package reltest

import (
	"fmt"
	"reflect"
)

// compact sprint struct ignoring zero values
func csprint(v interface{}, parent bool) string {
	var (
		notEmpty bool
		str      string
		rv       = reflect.ValueOf(v)
		rt       = rv.Type()
	)

	if rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return ""
		}

		rv = rv.Elem()
		rt = rt.Elem()
		str += "&"
	}

	switch rt.Kind() {
	case reflect.Struct:
		str += fmt.Sprintf("%s{", rt.String())
		for i := 0; i < rt.NumField(); i++ {
			var (
				fv    = rv.Field(i)
				ft    = rt.Field(i)
				fvstr = csprint(fv.Interface(), false)
			)

			if fvstr != "" {
				if notEmpty && i > 0 {
					str += ", "
				}

				str += fmt.Sprintf("%s: %s", ft.Name, fvstr)
				notEmpty = true
			}
		}
		str += "}"
	case reflect.Slice:
		str += fmt.Sprintf("%s{", rt.String())
		for i := 0; i < rv.Len(); i++ {
			if i > 0 {
				str += ", "
			}

			var (
				fv    = rv.Index(i)
				fvstr = csprint(fv.Interface(), false)
			)

			if fvstr != "" {
				str += fvstr
				notEmpty = true
			} else {
				str += fmt.Sprintf("%s{}", fv.Type().String())
			}
		}
		str += "}"
	default:
		if !rv.IsZero() {
			str = fmt.Sprintf("%v", rv.Interface())
			notEmpty = true
		}
	}

	if !notEmpty && !parent {
		str = ""
	}

	return str
}

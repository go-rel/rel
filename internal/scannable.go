package internal

import (
	"database/sql"
	"reflect"
	"time"
)

// Scannable checks whether type is scannable.
func Scannable(rt reflect.Type) bool {
	if rt.Kind() == reflect.Ptr {
		fzeroval := reflect.New(rt.Elem()).Interface()
		kind := rt.Elem().Kind()
		_, isScanner := fzeroval.(sql.Scanner)
		_, isTime := fzeroval.(*time.Time)

		if (kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) && kind != reflect.Uint8 && !isScanner && !isTime {
			return false
		}
	} else {
		fzeroval := reflect.New(rt).Interface()
		kind := rt.Kind()
		_, isScanner := fzeroval.(sql.Scanner)
		_, isTime := fzeroval.(*time.Time)

		if (kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) && kind != reflect.Uint8 && !isScanner && !isTime {
			return false
		}
	}
	return true
}

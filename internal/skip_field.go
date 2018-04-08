package internal

import (
	"database/sql"
	"reflect"
	"time"
)

// SkipType checks whether type should be skipped when scanning or generating changes.
func SkipType(rt reflect.Type) bool {
	if rt.Kind() == reflect.Ptr {
		fzeroval := reflect.New(rt.Elem()).Interface()
		kind := rt.Elem().Kind()
		_, isScanner := fzeroval.(sql.Scanner)
		_, isTime := fzeroval.(*time.Time)

		if (kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) && kind != reflect.Uint8 && !isScanner && !isTime {
			return true
		}
	} else {
		fzeroval := reflect.New(rt).Interface()
		kind := rt.Kind()
		_, isScanner := fzeroval.(sql.Scanner)
		_, isTime := fzeroval.(*time.Time)

		if (kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) && kind != reflect.Uint8 && !isScanner && !isTime {
			return true
		}
	}
	return false
}

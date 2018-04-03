package sqlutil

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/azer/snakecase"
)

// minimum rows interface for test purpose
type Rows interface {
	Scan(dest ...interface{}) error
	Columns() ([]string, error)
	Next() bool
}

var typeScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func Scan(value interface{}, rows Rows) (int64, error) {
	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("value must be pointer")
	}

	count := int64(0)
	rv = rv.Elem()
	var index map[string]int
	isScanner := rv.Addr().Type().Implements(typeScanner)
	isSlice := rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() != reflect.Uint8 && !isScanner

	if !isScanner {
		if isSlice {
			index = fieldIndex(rv.Type().Elem())
		} else {
			index = fieldIndex(rv.Type())
		}
	}

	for rows.Next() {
		var elem reflect.Value
		if isSlice {
			elem = reflect.New(rv.Type().Elem()).Elem()
		} else {
			elem = rv
		}

		var ptr []interface{}
		if isScanner {
			ptr = []interface{}{elem.Addr().Interface()}
		} else {
			ptr = fieldPtr(elem, index, columns)
		}

		err = rows.Scan(ptr...)
		if err != nil {
			return 0, err
		}

		count += 1

		if isSlice {
			rv.Set(reflect.Append(rv, elem))
		} else {
			break
		}
	}

	return count, nil
}

func fieldPtr(rv reflect.Value, index map[string]int, columns []string) []interface{} {
	var ptr []interface{}

	dummy := sql.RawBytes{}
	for _, col := range columns {
		if id, exist := index[col]; exist {
			ptr = append(ptr, rv.Field(id).Addr().Interface())
		} else {
			ptr = append(ptr, &dummy)
		}
	}

	return ptr
}

func fieldIndex(rt reflect.Type) map[string]int {
	fields := make(map[string]int)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		// skip if struct or slice but not a scanner or time
		if f.Type.Kind() == reflect.Ptr {
			fzeroval := reflect.New(f.Type.Elem()).Interface()
			kind := f.Type.Elem().Kind()
			_, isScanner := fzeroval.(sql.Scanner)
			_, isTime := fzeroval.(*time.Time)

			if (kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) && kind != reflect.Uint8 && !isScanner && !isTime {
				continue
			}
		} else {
			fzeroval := reflect.New(f.Type).Interface()
			kind := f.Type.Kind()
			_, isScanner := fzeroval.(sql.Scanner)
			_, isTime := fzeroval.(*time.Time)

			if (kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) && kind != reflect.Uint8 && !isScanner && !isTime {
				continue
			}
		}

		if tag := f.Tag.Get("db"); tag != "" {
			if tag == "-" {
				continue
			}

			fields[tag] = i
		} else {
			fields[snakecase.SnakeCase(f.Name)] = i
		}
	}

	return fields
}

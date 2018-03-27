package sqlutil

import (
	"database/sql"
	"reflect"
	"unicode"
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

var tag = "db"

func fieldIndex(rt reflect.Type) map[string]int {
	fields := make(map[string]int)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if tag := f.Tag.Get(tag); tag != "" {
			if tag == "-" {
				continue
			}

			fields[tag] = i
		} else {
			fields[toSnake(f.Name)] = i
		}
	}

	return fields
}

// convert string to snake case
// https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func toSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

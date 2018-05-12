package sql

import (
	"database/sql"
	"reflect"

	"github.com/Fs02/grimoire/internal"
	"github.com/azer/snakecase"
)

// Rows is minimal rows interface for test purpose
type Rows interface {
	Scan(dest ...interface{}) error
	Columns() ([]string, error)
	Next() bool
}

var typeScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

// Scan rows into interface
func Scan(value interface{}, rows Rows) (int64, error) {
	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("grimoire: record parameter must be a pointer")
	}

	count := int64(0)
	rv = rv.Elem()
	var index map[string]int
	isScanner := rv.Addr().Type().Implements(typeScanner)
	isSlice := rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() != reflect.Uint8 && !isScanner

	if !isScanner {
		if isSlice {
			rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
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

		count++

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

		// skip if not scannable
		if !internal.Scannable(f.Type) {
			continue
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

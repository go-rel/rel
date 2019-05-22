package sql

import (
	"reflect"

	"github.com/Fs02/grimoire/schema"
)

// Rows is minimal rows interface for test purpose
type Rows interface {
	Scan(dest ...interface{}) error
	Columns() ([]string, error)
	Next() bool
}

// Scan rows into interface
func Scan(value interface{}, rows Rows) (int64, error) {
	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	rt := reflect.TypeOf(value)
	if rt.Kind() != reflect.Ptr {
		panic("grimoire: record parameter must be a pointer")
	}

	if rt.Elem().Kind() == reflect.Slice {
		return scanMany(value, rows, columns)
	}

	return scanOne(value, rows, columns)
}

func scanOne(out interface{}, rows Rows, columns []string) (int64, error) {
	if !rows.Next() {
		return 0, nil
	}

	var (
		ptr = schema.InferScanners(out, columns)
		err = rows.Scan(ptr...)
	)

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func scanMany(out interface{}, rows Rows, columns []string) (int64, error) {
	var (
		rv      = reflect.ValueOf(out).Elem()
		sliceRv = reflect.MakeSlice(rv.Type(), 0, 0)
		elemRt  = rv.Type().Elem()
		elemRv  = reflect.New(elemRt).Elem()
	)

	for rows.Next() {
		var (
			copyElem = elemRv
			ptr      = schema.InferScanners(copyElem.Addr().Interface(), columns)
			err      = rows.Scan(ptr...)
		)

		if err != nil {
			return 0, err
		}

		sliceRv = reflect.Append(sliceRv, copyElem)
	}

	rv.Set(sliceRv)

	return int64(sliceRv.Len()), nil
}

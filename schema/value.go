package schema

import (
	"database/sql"
	"reflect"
)

type value struct {
	dest interface{}
}

func (v value) Scan(src interface{}) error {
	return convertAssign(v.dest, src)
}

func Value(dest interface{}) interface{} {
	if s, ok := dest.(sql.Scanner); ok {
		return s
	}

	rt := reflect.TypeOf(dest)
	if rt.Kind() != reflect.Ptr {
		panic("grimoire: destination must be a pointer")
	}

	if rt.Elem().Kind() == reflect.Ptr {
		return dest
	}

	return value{
		dest: dest,
	}
}

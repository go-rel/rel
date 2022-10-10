package rel

import (
	"database/sql"
	"reflect"
)

type nullable struct {
	dest any
}

var _ sql.Scanner = (*nullable)(nil)

func (n nullable) Scan(src any) error {
	return convertAssign(n.dest, src)
}

// Nullable wrap value as a nullable sql.Scanner.
// If value returned from database is nil, nullable scanner will set dest to zero value.
func Nullable(dest any) any {
	if s, ok := dest.(sql.Scanner); ok {
		return s
	}

	rt := reflect.TypeOf(dest)
	if rt.Kind() != reflect.Ptr {
		panic("rel: destination must be a pointer")
	}

	if rt.Elem().Kind() == reflect.Ptr {
		return dest
	}

	return nullable{
		dest: dest,
	}
}

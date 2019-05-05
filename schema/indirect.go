package schema

import (
	"database/sql"
)

type indirect struct {
	dest interface{}
}

func (i indirect) Scan(src interface{}) error {
	return convertAssign(i.dest, src)
}

func Indirect(dest interface{}) sql.Scanner {
	if s, ok := dest.(sql.Scanner); ok {
		return s
	}

	// rt := reflect.TypeOf(dest)
	// if rt.Kind() != reflect.Ptr || rt.Elem().Kind() == reflect.Ptr{
	// 	panic("grimoire: destination must be a pointer to a value")
	// }

	return indirect{
		dest: dest,
	}
}

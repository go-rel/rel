package schema

import (
	"database/sql"
)

type indirect struct {
	dest interface{}
}

func (i indirect) Scan(src interface{}) error {
	return nil
}

func Indirect(dest interface{}) sql.Scanner {
	return indirect{
		dest: dest,
	}
}

package sql

import (
	"database/sql"
)

type Cursor struct {
	*sql.Rows
}

func (c *Cursor) Fields() ([]string, error) {
	return c.Columns()
}

func (c *Cursor) NopScanner() interface{} {
	return &sql.RawBytes{}
}

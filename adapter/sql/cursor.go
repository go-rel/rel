package sql

import (
	"database/sql"
)

// Cursor used for retrieving result.
type Cursor struct {
	*sql.Rows
}

// Fields returned in the result.
func (c *Cursor) Fields() ([]string, error) {
	return c.Columns()
}

// NopScanner for this adapter.
func (c *Cursor) NopScanner() interface{} {
	return &sql.RawBytes{}
}

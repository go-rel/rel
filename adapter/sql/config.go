package sql

import (
	"github.com/Fs02/rel"
)

// Config holds configuration for adapter.
type Config struct {
	Placeholder         string
	Ordinal             bool
	InsertDefaultValues bool
	EscapeChar          string
	ErrorFunc           func(error) error
	IncrementFunc       func(Adapter) int
	MapColumnTypeFunc   func(column rel.Column) (string, int, int)
}

// MapColumnType func.
func MapColumnType(column rel.Column) (string, int, int) {
	var (
		typ  string
		m, n int
	)

	switch column.Type {
	case rel.ID:
		typ = "INT UNSIGNED AUTO_INCREMENT PRIMARY KEY"
	case rel.Bool:
		typ = "BOOL"
	case rel.Int:
		typ = "INT"
		m = column.Limit
	case rel.BigInt:
		typ = "BIGINT"
		m = column.Limit
	case rel.Float:
		typ = "FLOAT"
		m = column.Precision
	case rel.Decimal:
		typ = "DECIMAL"
		m = column.Precision
		n = column.Scale
	case rel.String:
		typ = "VARCHAR"
		m = column.Limit
		if m == 0 {
			m = 255
		}
	case rel.Text:
		typ = "TEXT"
		m = column.Limit
	case rel.Binary:
		typ = "BINARY"
		m = column.Limit
	case rel.Date:
		typ = "DATE"
	case rel.DateTime:
		typ = "DATETIME"
	case rel.Time:
		typ = "TIME"
	case rel.Timestamp:
		// TODO: mysql automatically add on update options.
		typ = "TIMESTAMP"
	default:
		typ = string(column.Type)
	}

	return typ, m, n
}

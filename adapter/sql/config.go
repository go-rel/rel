package sql

import (
	"time"

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
	MapColumnFunc       func(column *rel.Column) (string, int, int)
}

// MapColumn func.
func MapColumn(column *rel.Column) (string, int, int) {
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
	case rel.Date:
		typ = "DATE"
		if t, ok := column.Default.(time.Time); ok {
			column.Default = t.Format("2006-01-02")
		}
	case rel.DateTime:
		typ = "DATETIME"
		if t, ok := column.Default.(time.Time); ok {
			column.Default = t.Format("2006-01-02 15:04:05")
		}
	case rel.Time:
		typ = "TIME"
		if t, ok := column.Default.(time.Time); ok {
			column.Default = t.Format("15:04:05")
		}
	case rel.Timestamp:
		// TODO: mysql automatically add on update options.
		typ = "TIMESTAMP"
		if t, ok := column.Default.(time.Time); ok {
			column.Default = t.Format("2006-01-02 15:04:05")
		}
	default:
		typ = string(column.Type)
	}

	return typ, m, n
}

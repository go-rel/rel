package rel

import (
	"context"
)

// Adapter interface
type Adapter interface {
	Ping(ctx context.Context) error
	Aggregate(ctx context.Context, query Query, mode string, field string, loggers ...Logger) (int, error)
	Query(ctx context.Context, query Query, loggers ...Logger) (Cursor, error)
	Insert(ctx context.Context, query Query, modifies map[string]Modify, loggers ...Logger) (interface{}, error)
	InsertAll(ctx context.Context, query Query, fields []string, bulkModifies []map[string]Modify, loggers ...Logger) ([]interface{}, error)
	Update(ctx context.Context, query Query, modifies map[string]Modify, loggers ...Logger) (int, error)
	Delete(ctx context.Context, query Query, loggers ...Logger) (int, error)

	Begin(ctx context.Context) (Adapter, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

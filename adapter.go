package rel

import (
	"context"
)

// Adapter interface
type Adapter interface {
	Instrumentation(instrumenter Instrumenter)
	Ping(ctx context.Context) error
	Aggregate(ctx context.Context, query Query, mode string, field string) (int, error)
	Query(ctx context.Context, query Query) (Cursor, error)
	Insert(ctx context.Context, query Query, modifies map[string]Modify) (interface{}, error)
	InsertAll(ctx context.Context, query Query, fields []string, bulkModifies []map[string]Modify) ([]interface{}, error)
	Update(ctx context.Context, query Query, modifies map[string]Modify) (int, error)
	Delete(ctx context.Context, query Query) (int, error)

	Begin(ctx context.Context) (Adapter, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

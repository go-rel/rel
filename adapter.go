package rel

import (
	"context"
)

// Adapter interface
type Adapter interface {
	Name() string
	Close() error

	Instrumentation(instrumenter Instrumenter)
	Ping(ctx context.Context) error
	Aggregate(ctx context.Context, query Query, mode string, field string) (int, error)
	Query(ctx context.Context, query Query) (Cursor, error)
	Insert(ctx context.Context, query Query, primaryField string, mutates map[string]Mutate, onConflict OnConflict) (any, error)
	InsertAll(ctx context.Context, query Query, primaryField string, fields []string, bulkMutates []map[string]Mutate, onConflict OnConflict) ([]any, error)
	Update(ctx context.Context, query Query, primaryField string, mutates map[string]Mutate) (int, error)
	Delete(ctx context.Context, query Query) (int, error)
	Exec(ctx context.Context, stmt string, args []any) (int64, int64, error)

	Begin(ctx context.Context) (Adapter, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	Apply(ctx context.Context, migration Migration) error
}

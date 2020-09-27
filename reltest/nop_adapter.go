package reltest

import (
	"context"

	"github.com/go-rel/rel"
)

type nopAdapter struct {
}

func (na *nopAdapter) Instrumentation(instrumenter rel.Instrumenter) {
}

func (na *nopAdapter) Ping(ctx context.Context) error {
	return nil
}

func (na *nopAdapter) Aggregate(ctx context.Context, query rel.Query, mode string, field string) (int, error) {
	return 0, nil
}

func (na *nopAdapter) Begin(ctx context.Context) (rel.Adapter, error) {
	return na, nil
}

func (na *nopAdapter) Commit(ctx context.Context) error {
	return nil
}

func (na *nopAdapter) Delete(ctx context.Context, query rel.Query) (int, error) {
	return 1, nil
}

func (na *nopAdapter) Insert(ctx context.Context, query rel.Query, primaryField string, mutates map[string]rel.Mutate) (interface{}, error) {
	return 1, nil
}

func (na *nopAdapter) InsertAll(ctx context.Context, query rel.Query, primaryField string, fields []string, bulkMutates []map[string]rel.Mutate) ([]interface{}, error) {
	var (
		ids = make([]interface{}, len(bulkMutates))
	)

	for i := range bulkMutates {
		ids[i] = i + 1
	}

	return ids, nil
}

func (na *nopAdapter) Query(ctx context.Context, query rel.Query) (rel.Cursor, error) {
	return &nopCursor{count: 1}, nil
}

func (na *nopAdapter) Rollback(ctx context.Context) error {
	return nil
}

func (na *nopAdapter) Update(ctx context.Context, query rel.Query, mutates map[string]rel.Mutate) (int, error) {
	return 1, nil
}

func (na *nopAdapter) Apply(ctx context.Context, migration rel.Migration) error {
	return nil
}

type nopCursor struct {
	count int
}

func (nc *nopCursor) Close() error {
	return nil
}

func (nc *nopCursor) Fields() ([]string, error) {
	return nil, nil
}

func (nc *nopCursor) Next() bool {
	nc.count--
	return nc.count >= 0
}

func (nc *nopCursor) Scan(...interface{}) error {
	nc.NopScanner()
	return nil
}

func (nc *nopCursor) NopScanner() interface{} {
	return nil
}

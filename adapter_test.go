package rel

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type testAdapter struct {
	mock.Mock
	result interface{}
}

var _ Adapter = (*testAdapter)(nil)

func (ta *testAdapter) Open(dsn string) error {
	args := ta.Called(dsn)
	return args.Error(0)
}

func (ta *testAdapter) Close() error {
	args := ta.Called()
	return args.Error(0)
}

func (ta *testAdapter) Instrumentation(instrumenter Instrumenter) {
}

func (ta *testAdapter) Ping(ctx context.Context) error {
	args := ta.Called()
	return args.Error(0)
}

func (ta *testAdapter) Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error) {
	args := ta.Called(query, aggregate, field)
	return args.Int(0), args.Error(1)
}

func (ta *testAdapter) Query(ctx context.Context, query Query) (Cursor, error) {
	args := ta.Called(query)
	return args.Get(0).(Cursor), args.Error(1)
}

func (ta *testAdapter) Insert(ctx context.Context, query Query, primaryField string, mutates map[string]Mutate) (interface{}, error) {
	args := ta.Called(query, mutates)
	return args.Get(0), args.Error(1)
}

func (ta *testAdapter) InsertAll(ctx context.Context, query Query, primaryField string, fields []string, mutates []map[string]Mutate) ([]interface{}, error) {
	args := ta.Called(query, fields, mutates)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (ta *testAdapter) Update(ctx context.Context, query Query, primaryField string, mutates map[string]Mutate) (int, error) {
	args := ta.Called(query, primaryField, mutates)
	return args.Int(0), args.Error(1)
}

func (ta *testAdapter) Delete(ctx context.Context, query Query) (int, error) {
	args := ta.Called(query)
	return args.Int(0), args.Error(1)
}

func (ta *testAdapter) Begin(ctx context.Context) (Adapter, error) {
	args := ta.Called()
	return ta, args.Error(0)
}

func (ta *testAdapter) Commit(ctx context.Context) error {
	args := ta.Called()
	return args.Error(0)
}

func (ta *testAdapter) Rollback(ctx context.Context) error {
	args := ta.Called()
	return args.Error(0)
}

func (ta *testAdapter) Apply(ctx context.Context, migration Migration) error {
	args := ta.Called(migration)
	return args.Error(0)
}

func (ta *testAdapter) Result(result interface{}) *testAdapter {
	ta.result = result
	return ta
}

func (ta *testAdapter) Exec(ctx context.Context, stmt string, args []interface{}) (int64, int64, error) {
	mockArgs := ta.Called(ctx, stmt, args)
	return int64(mockArgs.Int(0)), int64(mockArgs.Int(1)), mockArgs.Error(2)
}

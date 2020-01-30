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

func (ta *testAdapter) Aggregate(ctx context.Context, query Query, aggregate string, field string, logger ...Logger) (int, error) {
	args := ta.Called(query, aggregate, field)
	return args.Int(0), args.Error(1)
}

func (ta *testAdapter) Query(ctx context.Context, query Query, logger ...Logger) (Cursor, error) {
	args := ta.Called(query)
	return args.Get(0).(Cursor), args.Error(1)
}

func (ta *testAdapter) Insert(ctx context.Context, query Query, modifies map[string]Modify, logger ...Logger) (interface{}, error) {
	args := ta.Called(query, modifies)
	return args.Get(0), args.Error(1)
}

func (ta *testAdapter) InsertAll(ctx context.Context, query Query, fields []string, modifies []map[string]Modify, logger ...Logger) ([]interface{}, error) {
	args := ta.Called(query, fields, modifies)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (ta *testAdapter) Update(ctx context.Context, query Query, modifies map[string]Modify, logger ...Logger) (int, error) {
	args := ta.Called(query, modifies)
	return args.Int(0), args.Error(1)
}

func (ta *testAdapter) Delete(ctx context.Context, query Query, logger ...Logger) (int, error) {
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

func (ta *testAdapter) Result(result interface{}) *testAdapter {
	ta.result = result
	return ta
}

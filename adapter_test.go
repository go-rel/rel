package grimoire

import (
	"github.com/Fs02/grimoire/change"
	"github.com/Fs02/grimoire/query"
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

func (ta *testAdapter) Aggregate(queries query.Query, out interface{}, mode string, field string, logger ...Logger) error {
	args := ta.Called(queries, out, mode, field)
	return args.Error(0)
}

func (ta *testAdapter) Query(queries query.Query, logger ...Logger) (Cursor, error) {
	args := ta.Called(queries)
	return args.Get(0).(Cursor), args.Error(1)
}

func (ta *testAdapter) Insert(queries query.Query, changes change.Changes, logger ...Logger) (interface{}, error) {
	args := ta.Called(queries, changes)
	return args.Get(0), args.Error(1)
}

func (ta *testAdapter) InsertAll(queries query.Query, fields []string, changess []change.Changes, logger ...Logger) ([]interface{}, error) {
	args := ta.Called(queries, changess)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (ta *testAdapter) Update(queries query.Query, changes change.Changes, logger ...Logger) error {
	args := ta.Called(queries, changes)
	return args.Error(0)
}

func (ta *testAdapter) Delete(queries query.Query, logger ...Logger) error {
	args := ta.Called(queries)
	return args.Error(0)
}

func (ta *testAdapter) Begin() (Adapter, error) {
	args := ta.Called()
	return ta, args.Error(0)
}

func (ta *testAdapter) Commit() error {
	args := ta.Called()
	return args.Error(0)
}

func (ta *testAdapter) Rollback() error {
	args := ta.Called()
	return args.Error(0)
}

func (ta *testAdapter) Result(result interface{}) *testAdapter {
	ta.result = result
	return ta
}

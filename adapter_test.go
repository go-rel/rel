package grimoire

import (
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

func (ta *testAdapter) Aggregate(query Query, aggregate string, field string, logger ...Logger) (int, error) {
	args := ta.Called(query, aggregate, field)
	return args.Int(0), args.Error(1)
}

func (ta *testAdapter) Query(query Query, logger ...Logger) (Cursor, error) {
	args := ta.Called(query)
	return args.Get(0).(Cursor), args.Error(1)
}

func (ta *testAdapter) Insert(query Query, changes Changes, logger ...Logger) (interface{}, error) {
	args := ta.Called(query, changes)
	return args.Get(0), args.Error(1)
}

func (ta *testAdapter) InsertAll(query Query, fields []string, changess []Changes, logger ...Logger) ([]interface{}, error) {
	args := ta.Called(query, changess)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (ta *testAdapter) Update(query Query, changes Changes, logger ...Logger) error {
	args := ta.Called(query, changes)
	return args.Error(0)
}

func (ta *testAdapter) Delete(query Query, logger ...Logger) error {
	args := ta.Called(query)
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

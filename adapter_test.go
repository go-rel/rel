package grimoire

import (
	"github.com/stretchr/testify/mock"
)

type TestAdapter struct {
	mock.Mock
}

var _ Adapter = (*TestAdapter)(nil)

func (adapter TestAdapter) Open(dsn string) error {
	args := adapter.Called(dsn)
	return args.Error(0)
}

func (adapter TestAdapter) Close() error {
	args := adapter.Called()
	return args.Error(0)
}

func (adapter TestAdapter) Count(query Query, logger ...Logger) (int, error) {
	args := adapter.Called(query)
	return args.Int(0), args.Error(1)
}

func (adapter TestAdapter) All(query Query, doc interface{}, logger ...Logger) (int, error) {
	args := adapter.Called(query, doc)
	return args.Int(0), args.Error(1)
}

func (adapter TestAdapter) Insert(query Query, ch map[string]interface{}, logger ...Logger) (interface{}, error) {
	args := adapter.Called(query, ch)
	return args.Get(0), args.Error(1)
}

func (adapter TestAdapter) InsertAll(query Query, fields []string, chs []map[string]interface{}, logger ...Logger) ([]interface{}, error) {
	args := adapter.Called(query, chs)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (adapter TestAdapter) Update(query Query, ch map[string]interface{}, logger ...Logger) error {
	args := adapter.Called(query, ch)
	return args.Error(0)
}

func (adapter TestAdapter) Delete(query Query, logger ...Logger) error {
	args := adapter.Called(query)
	return args.Error(0)
}

func (adapter TestAdapter) Begin() (Adapter, error) {
	args := adapter.Called()
	return adapter, args.Error(0)
}

func (adapter TestAdapter) Commit() error {
	args := adapter.Called()
	return args.Error(0)
}

func (adapter TestAdapter) Rollback() error {
	args := adapter.Called()
	return args.Error(0)
}

package grimoire

import (
	"github.com/Fs02/grimoire/change"
	"github.com/Fs02/grimoire/query"
	"github.com/stretchr/testify/mock"
)

type TestAdapter struct {
	mock.Mock
	result interface{}
}

var _ Adapter = (*TestAdapter)(nil)

func (adapter *TestAdapter) Open(dsn string) error {
	args := adapter.Called(dsn)
	return args.Error(0)
}

func (adapter *TestAdapter) Close() error {
	args := adapter.Called()
	return args.Error(0)
}

func (adapter *TestAdapter) Aggregate(queries query.Query, out interface{}, mode string, field string, logger ...Logger) error {
	args := adapter.Called(queries, out, mode, field)
	return args.Error(0)
}

func (adapter *TestAdapter) All(queries query.Query, collection Collection, logger ...Logger) (int, error) {
	args := adapter.Called(queries, collection)

	// if adapter.result != nil {
	// 	switch doc.(type) {
	// 	case *[]Address:
	// 		*doc.(*[]Address) = adapter.result.([]Address)
	// 	case *[]Transaction:
	// 		*doc.(*[]Transaction) = adapter.result.([]Transaction)
	// 	case *[]User:
	// 		*doc.(*[]User) = adapter.result.([]User)
	// 	default:
	// 		panic("not implemented")
	// 	}
	// }

	return args.Int(0), args.Error(1)
}

func (adapter *TestAdapter) Insert(queries query.Query, changes change.Changes, logger ...Logger) (interface{}, error) {
	args := adapter.Called(queries, changes)
	return args.Get(0), args.Error(1)
}

func (adapter *TestAdapter) InsertAll(queries query.Query, fields []string, changess []change.Changes, logger ...Logger) ([]interface{}, error) {
	args := adapter.Called(queries, changess)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (adapter *TestAdapter) Update(queries query.Query, changes change.Changes, logger ...Logger) error {
	args := adapter.Called(queries, changes)
	return args.Error(0)
}

func (adapter *TestAdapter) Delete(queries query.Query, logger ...Logger) error {
	args := adapter.Called(queries)
	return args.Error(0)
}

func (adapter *TestAdapter) Begin() (Adapter, error) {
	args := adapter.Called()
	return adapter, args.Error(0)
}

func (adapter *TestAdapter) Commit() error {
	args := adapter.Called()
	return args.Error(0)
}

func (adapter *TestAdapter) Rollback() error {
	args := adapter.Called()
	return args.Error(0)
}

func (adapter *TestAdapter) Result(result interface{}) *TestAdapter {
	adapter.result = result
	return adapter
}

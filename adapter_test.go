package grimoire

import (
	"github.com/Fs02/grimoire/changeset"
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

func (adapter TestAdapter) Find(query Query) (string, []interface{}) {
	args := adapter.Called(query)
	return args.String(0), args.Get(1).([]interface{})
}

func (adapter TestAdapter) Insert(query Query, ch changeset.Changeset) (string, []interface{}) {
	args := adapter.Called(query, ch)
	return args.String(0), args.Get(1).([]interface{})
}

func (adapter TestAdapter) Update(query Query, ch changeset.Changeset) (string, []interface{}) {
	args := adapter.Called(query, ch)
	return args.String(0), args.Get(1).([]interface{})
}

func (adapter TestAdapter) Delete(query Query) (string, []interface{}) {
	args := adapter.Called(query)
	return args.String(0), args.Get(1).([]interface{})
}

func (adapter TestAdapter) Begin() error {
	args := adapter.Called()
	return args.Error(0)
}

func (adapter TestAdapter) Commit() error {
	args := adapter.Called()
	return args.Error(0)
}

func (adapter TestAdapter) Rollback() error {
	args := adapter.Called()
	return args.Error(0)
}

func (adapter TestAdapter) Query(out interface{}, qs string, qargs []interface{}) (int64, error) {
	args := adapter.Called(out, qs, qargs)
	return args.Get(0).(int64), args.Error(1)
}

func (adapter TestAdapter) Exec(qs string, qargs []interface{}) (int64, int64, error) {
	args := adapter.Called(qs, qargs)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

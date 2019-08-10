package grimoire

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testCursor struct {
	mock.Mock
}

var _ Cursor = (*testCursor)(nil)

func (tc *testCursor) Close() error {
	ret := tc.Called()
	return ret.Error(0)
}

func (tc *testCursor) Fields() ([]string, error) {
	ret := tc.Called()
	return ret.Get(0).([]string), ret.Error(1)
}

func (tc *testCursor) Next() bool {
	ret := tc.Called()
	return ret.Get(0).(bool)
}

func (tc *testCursor) NopScanner() interface{} {
	return &sql.RawBytes{}
}

func (tc *testCursor) Scan(scanners ...interface{}) error {
	ret := tc.Called(scanners...)

	var err error
	if fn, ok := ret.Get(0).(func(...interface{}) error); ok {
		err = fn(scanners...)
	} else {
		err = ret.Error(0)
	}

	return err
}

func (tc *testCursor) SetScan(times int, ret ...interface{}) *mock.Call {
	args := make([]interface{}, len(ret))
	for i := 0; i < len(args); i++ {
		args[i] = mock.MatchedBy(func(arg interface{}) bool {
			_, ok := arg.(sql.Scanner)
			return ok
		})
	}

	return tc.On("Scan", args...).
		Return(func(scanners ...interface{}) error {
			for i := 0; i < len(scanners); i++ {
				scanners[i].(sql.Scanner).Scan(ret[i])
			}

			return nil
		}).Times(times)
}

func (tc *testCursor) SetScanHead(times int, ret ...interface{}) *mock.Call {
	args := make([]interface{}, len(ret))

	args[0] = mock.MatchedBy(func(arg interface{}) bool {
		_, ok := arg.(*int)
		return ok
	})

	for i := 1; i < len(args); i++ {
		args[i] = mock.Anything
	}

	return tc.On("Scan", args...).
		Return(func(scanners ...interface{}) error {
			for i := 0; i < len(scanners); i++ {
				if v, ok := scanners[i].(*int); ok {
					*v = ret[i].(int)
				}
			}

			return nil
		}).Once()
}

func TestScanOne(t *testing.T) {
	var (
		user User
		cur  = &testCursor{}
		doc  = newDocument(&user)
		now  = time.Now()
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.SetScan(1, 10, "del piero", nil, now, nil)

	err := scanOne(cur, doc)
	assert.Nil(t, err)

	assert.Equal(t, User{
		ID:        10,
		Name:      "del piero",
		CreatedAt: now,
	}, user)

	cur.AssertExpectations(t)
}

func TestScanMany(t *testing.T) {
	var (
		users []User
		cur   = &testCursor{}
		col   = newCollection(&users)
		now   = time.Now()
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()

	cur.On("Next").Return(true).Twice()
	cur.SetScan(1, 10, "del piero", nil, now, nil)
	cur.SetScan(1, 11, "nedved", 46, now, now)
	cur.On("Next").Return(false).Once()

	err := scanMany(cur, col)
	assert.Nil(t, err)
	assert.Len(t, users, 2)

	assert.Equal(t, User{
		ID:        10,
		Name:      "del piero",
		CreatedAt: now,
	}, users[0])

	assert.Equal(t, User{
		ID:        11,
		Name:      "nedved",
		Age:       46,
		CreatedAt: now,
		UpdatedAt: now,
	}, users[1])

	cur.AssertExpectations(t)
}

func TestScanMulti(t *testing.T) {
	var (
		users1   []User
		users2   []User
		users3   []User
		cur      = &testCursor{}
		keyField = "id"
		keyType  = reflect.TypeOf(0)
		cols     = map[interface{}][]Collection{
			10: {newCollection(&users1), newCollection(&users2)},
			11: {newCollection(&users3)},
		}
		now = time.Now()
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()

	cur.On("Next").Return(true).Twice()
	cur.SetScanHead(1, 10, nil, nil, nil, nil, nil)
	cur.SetScan(2, 10, "del piero", nil, now, nil)
	cur.SetScanHead(1, 11, nil, nil, nil, nil, nil)
	cur.SetScan(1, 11, "nedved", 46, now, now)
	cur.On("Next").Return(false).Once()

	err := scanMulti(cur, keyField, keyType, cols)
	assert.Nil(t, err)

	assert.Len(t, users1, 1)
	assert.Equal(t, User{
		ID:        10,
		Name:      "del piero",
		CreatedAt: now,
	}, users1[0])

	assert.Len(t, users2, 1)
	assert.Equal(t, User{
		ID:        10,
		Name:      "del piero",
		CreatedAt: now,
	}, users2[0])

	assert.Len(t, users3, 1)
	assert.Equal(t, User{
		ID:        11,
		Name:      "nedved",
		Age:       46,
		CreatedAt: now,
		UpdatedAt: now,
	}, users3[0])

	cur.AssertExpectations(t)
}

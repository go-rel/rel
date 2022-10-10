package rel

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

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

func (tc *testCursor) NopScanner() any {
	return &sql.RawBytes{}
}

func (tc *testCursor) Scan(scanners ...any) error {
	ret := tc.Called(scanners...)

	var err error
	if fn, ok := ret.Get(0).(func(...any) error); ok {
		err = fn(scanners...)
	} else {
		err = ret.Error(0)
	}

	return err
}

func (tc *testCursor) MockScan(ret ...any) *mock.Call {
	args := make([]any, len(ret))
	for i := 0; i < len(args); i++ {
		args[i] = mock.Anything
	}

	return tc.On("Scan", args...).
		Return(func(scanners ...any) error {
			for i := 0; i < len(scanners); i++ {
				if v, ok := scanners[i].(sql.Scanner); ok {
					v.Scan(ret[i])
				} else {
					convertAssign(scanners[i], ret[i])
				}
			}

			return nil
		})
}

func TestScanOne(t *testing.T) {
	var (
		user User
		cur  = &testCursor{}
		doc  = NewDocument(&user)
		now  = Now()
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(10, "Del Piero", nil, now, nil).Once()

	assert.Nil(t, scanOne(cur, doc))
	assert.Equal(t, User{
		ID:        10,
		Name:      "Del Piero",
		CreatedAt: now,
	}, user)

	cur.AssertExpectations(t)
}

func TestScanOne_fieldsError(t *testing.T) {
	var (
		user User
		cur  = &testCursor{}
		doc  = NewDocument(&user)
		err  = errors.New("field error")
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{}, err).Once()

	assert.Equal(t, err, scanOne(cur, doc))
	cur.AssertExpectations(t)
}

func TestScanAll(t *testing.T) {
	var (
		users []User
		cur   = &testCursor{}
		col   = NewCollection(&users)
		now   = Now()
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()

	cur.On("Next").Return(true).Twice()
	cur.MockScan(10, "Del Piero", nil, now, nil).Once()
	cur.MockScan(11, "Nedved", 46, now, now).Once()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, scanAll(cur, col))
	assert.Len(t, users, 2)
	assert.Equal(t, User{
		ID:        10,
		Name:      "Del Piero",
		CreatedAt: now,
	}, users[0])
	assert.Equal(t, User{
		ID:        11,
		Name:      "Nedved",
		Age:       46,
		CreatedAt: now,
		UpdatedAt: now,
	}, users[1])

	cur.AssertExpectations(t)
}

func TestScanAll_scanError(t *testing.T) {
	var (
		users []User
		cur   = &testCursor{}
		col   = NewCollection(&users)
		err   = errors.New("scan error")
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id"}, nil).Once()

	cur.On("Next").Return(true).Once()
	cur.On("Scan", mock.Anything).Return(err).Once()

	assert.Equal(t, err, scanAll(cur, col))
	cur.AssertExpectations(t)
}

func TestScanAll_fieldsError(t *testing.T) {
	var (
		users []User
		cur   = &testCursor{}
		col   = NewCollection(&users)
		err   = errors.New("field error")
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{}, err).Once()

	assert.Equal(t, err, scanAll(cur, col))
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
		cols     = map[any][]slice{
			10: {NewCollection(&users1), NewCollection(&users2)},
			11: {NewCollection(&users3)},
		}
		now = Now()
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()

	cur.On("Next").Return(true).Twice()
	cur.MockScan(10, "Del Piero", nil, now, nil).Times(3)
	cur.MockScan(11, "Nedved", 46, now, now).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, scanMulti(cur, keyField, keyType, cols))
	assert.Len(t, users1, 1)
	assert.Equal(t, User{
		ID:        10,
		Name:      "Del Piero",
		CreatedAt: now,
	}, users1[0])
	assert.Len(t, users2, 1)
	assert.Equal(t, User{
		ID:        10,
		Name:      "Del Piero",
		CreatedAt: now,
	}, users2[0])
	assert.Len(t, users3, 1)
	assert.Equal(t, User{
		ID:        11,
		Name:      "Nedved",
		Age:       46,
		CreatedAt: now,
		UpdatedAt: now,
	}, users3[0])

	cur.AssertExpectations(t)
}

func TestScanMulti_scanError(t *testing.T) {
	var (
		users    []User
		cur      = &testCursor{}
		keyField = "id"
		keyType  = reflect.TypeOf(0)
		cols     = map[any][]slice{
			11: {NewCollection(&users)},
		}
		err = errors.New("scan error")
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()

	cur.On("Next").Return(true).Once()
	cur.MockScan(11, "Nedved", 46, Now, Now).Once()
	cur.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(err).Once()

	assert.Equal(t, err, scanMulti(cur, keyField, keyType, cols))
	cur.AssertExpectations(t)
}

func TestScanMulti_scanKeyError(t *testing.T) {
	var (
		users    []User
		cur      = &testCursor{}
		keyField = "id"
		keyType  = reflect.TypeOf(0)
		cols     = map[any][]slice{
			11: {NewCollection(&users)},
		}
		err = errors.New("scan key error")
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id"}, nil).Once()

	cur.On("Next").Return(true).Once()
	cur.On("Scan", mock.Anything).Return(err).Once()

	assert.Equal(t, err, scanMulti(cur, keyField, keyType, cols))
	cur.AssertExpectations(t)
}

func TestScanMulti_keyFieldsNotExists(t *testing.T) {
	var (
		users    []User
		cur      = &testCursor{}
		keyField = "id"
		keyType  = reflect.TypeOf(0)
		cols     = map[any][]slice{
			11: {NewCollection(&users)},
		}
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{}, nil).Once()

	assert.Panics(t, func() {
		scanMulti(cur, keyField, keyType, cols)
	})
	cur.AssertExpectations(t)
}

func TestScanMulti_fieldsError(t *testing.T) {
	var (
		users    []User
		cur      = &testCursor{}
		keyField = "id"
		keyType  = reflect.TypeOf(0)
		cols     = map[any][]slice{
			11: {NewCollection(&users)},
		}
		err = errors.New("fields error")
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{}, err).Once()

	assert.Equal(t, err, scanMulti(cur, keyField, keyType, cols))
	cur.AssertExpectations(t)
}

func TestScanMulti_multipleTimes(t *testing.T) {
	var (
		users    = make([][]User, 6)
		cur      = &testCursor{}
		keyField = "id"
		keyType  = reflect.TypeOf(0)
		cols     = map[any][]slice{
			10: {NewCollection(&users[0]), NewCollection(&users[1])},
			11: {NewCollection(&users[2])},
			12: {NewCollection(&users[3]), NewCollection(&users[4])},
			13: {NewCollection(&users[5])},
		}
		now = Now()
	)

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()

	cur.On("Next").Return(true).Twice()
	cur.MockScan(10, "Del Piero", nil, now, nil).Times(3)
	cur.MockScan(11, "Nedved", 46, now, now).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, scanMulti(cur, keyField, keyType, cols))
	assert.Len(t, users[0], 1)
	assert.Equal(t, User{
		ID:        10,
		Name:      "Del Piero",
		CreatedAt: now,
	}, users[0][0])
	assert.Len(t, users[1], 1)
	assert.Equal(t, User{
		ID:        10,
		Name:      "Del Piero",
		CreatedAt: now,
	}, users[1][0])
	assert.Len(t, users[2], 1)
	assert.Equal(t, User{
		ID:        11,
		Name:      "Nedved",
		Age:       46,
		CreatedAt: now,
		UpdatedAt: now,
	}, users[2][0])

	cur.AssertExpectations(t)

	// Continue with a new cursor but the same cols -> works only if the ids in
	// the subsequent calls did not occur yet.
	cur = &testCursor{}

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name", "age", "created_at", "updated_at"}, nil).Once()

	cur.On("Next").Return(true).Twice()
	cur.MockScan(12, "Linus Torvalds", 52, now, nil).Times(3)
	cur.MockScan(13, "Tim Cook", 61, now, now).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, scanMulti(cur, keyField, keyType, cols))
	assert.Len(t, users[3], 1)
	assert.Equal(t, User{
		ID:        12,
		Name:      "Linus Torvalds",
		Age:       52,
		CreatedAt: now,
	}, users[3][0])
	assert.Len(t, users[4], 1)
	assert.Equal(t, User{
		ID:        12,
		Age:       52,
		Name:      "Linus Torvalds",
		CreatedAt: now,
	}, users[4][0])
	assert.Len(t, users[5], 1)
	assert.Equal(t, User{
		ID:        13,
		Name:      "Tim Cook",
		Age:       61,
		CreatedAt: now,
		UpdatedAt: now,
	}, users[5][0])

	cur.AssertExpectations(t)
}

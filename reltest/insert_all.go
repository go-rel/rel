package reltest

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type insertAll []*MockInsertAll

func (ia *insertAll) register(ctxData ctxData) *MockInsertAll {
	mia := &MockInsertAll{
		assert: &Assert{ctxData: ctxData},
	}
	*ia = append(*ia, mia)
	return mia
}

func (ia insertAll) execute(ctx context.Context, records interface{}) error {
	for _, mia := range ia {
		if (mia.argRecord == nil || reflect.DeepEqual(mia.argRecord, records)) &&
			(mia.argRecordType == "" || mia.argRecordType == reflect.TypeOf(records).String()) &&
			(mia.argRecordTable == "" || mia.argRecordTable == rel.NewCollection(records, true).Table()) &&
			mia.assert.call(ctx) {
			return mia.retError
		}
	}

	panic(failExecuteMessage(MockInsertAll{argRecord: records}, ia))
}

func (ia *insertAll) assert(t T) bool {
	for _, mia := range *ia {
		if !mia.assert.assert(t, mia) {
			return false
		}
	}

	*ia = nil
	return true
}

// MockInsertAll asserts and simulate Insert function for test.
type MockInsertAll struct {
	assert         *Assert
	argRecord      interface{}
	argRecordType  string
	argRecordTable string
	retError       error
}

// For assert calls for given record.
func (mia *MockInsertAll) For(record interface{}) *MockInsertAll {
	mia.argRecord = record
	return mia
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (mia *MockInsertAll) ForType(typ string) *MockInsertAll {
	mia.argRecordType = "*" + strings.TrimPrefix(typ, "*")
	return mia
}

// ForTable assert calls for given table.
func (mia *MockInsertAll) ForTable(typ string) *MockInsertAll {
	mia.argRecordTable = typ
	return mia
}

// Error sets error to be returned.
func (mia *MockInsertAll) Error(err error) *Assert {
	mia.retError = err
	return mia.assert
}

// Success sets no error to be returned.
func (mia *MockInsertAll) Success() *Assert {
	return mia.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (mia *MockInsertAll) ConnectionClosed() *Assert {
	return mia.Error(ErrConnectionClosed)
}

// NotUnique sets not unique error to be returned.
func (mia *MockInsertAll) NotUnique(key string) *Assert {
	return mia.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

// String representation of mocked call.
func (mia MockInsertAll) String() string {
	argRecord := "<Any>"
	if mia.argRecord != nil {
		argRecord = csprint(mia.argRecord, true)
	} else if mia.argRecordType != "" {
		argRecord = fmt.Sprintf("<Type: %s>", mia.argRecordType)
	} else if mia.argRecordTable != "" {
		argRecord = fmt.Sprintf("<Table: %s>", mia.argRecordTable)
	}

	return fmt.Sprintf("InsertAll(ctx, %s)", argRecord)
}

// ExpectString representation of mocked call.
func (mia MockInsertAll) ExpectString() string {
	return fmt.Sprintf("InsertAll().ForType(\"%T\")", mia.argRecord)
}

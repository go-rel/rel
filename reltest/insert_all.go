package reltest

import (
	"context"
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

	panic("TODO: Query doesn't match")
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
func (mm *MockInsertAll) For(record interface{}) *MockInsertAll {
	mm.argRecord = record
	return mm
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (mm *MockInsertAll) ForType(typ string) *MockInsertAll {
	mm.argRecordType = "*" + strings.TrimPrefix(typ, "*")
	return mm
}

// ForTable assert calls for given table.
func (mm *MockInsertAll) ForTable(typ string) *MockInsertAll {
	mm.argRecordTable = typ
	return mm
}

// Error sets error to be returned.
func (mm *MockInsertAll) Error(err error) *Assert {
	mm.retError = err
	return mm.assert
}

// Success sets no error to be returned.
func (mm *MockInsertAll) Success() *Assert {
	return mm.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (mm *MockInsertAll) ConnectionClosed() *Assert {
	return mm.Error(ErrConnectionClosed)
}

// NotUnique sets not unique error to be returned.
func (mm *MockInsertAll) NotUnique(key string) *Assert {
	return mm.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

package reltest

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type deleteAll []*MockDeleteAll

func (da *deleteAll) register(ctxData ctxData, options ...rel.Cascade) *MockDeleteAll {
	mda := &MockDeleteAll{
		assert:     &Assert{ctxData: ctxData},
		argOptions: options,
	}
	*da = append(*da, mda)
	return mda
}

func (da deleteAll) execute(ctx context.Context, record interface{}) error {
	for _, mda := range da {
		if (mda.argRecord == nil || reflect.DeepEqual(mda.argRecord, record)) &&
			(mda.argRecordType == "" || mda.argRecordType == reflect.TypeOf(record).String()) &&
			(mda.argRecordTable == "" || mda.argRecordTable == rel.NewCollection(record, true).Table()) &&
			mda.assert.call(ctx) {
			return mda.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockDeleteAll asserts and simulate Delete function for test.
type MockDeleteAll struct {
	assert         *Assert
	argRecord      interface{}
	argRecordType  string
	argRecordTable string
	argOptions     []rel.Cascade
	retError       error
}

// For assert calls for given record.
func (mda *MockDeleteAll) For(record interface{}) *MockDeleteAll {
	mda.argRecord = record
	return mda
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (mda *MockDeleteAll) ForType(typ string) *MockDeleteAll {
	mda.argRecordType = "*" + strings.TrimPrefix(typ, "*")
	return mda
}

// ForTable assert calls for given table.
func (mda *MockDeleteAll) ForTable(typ string) *MockDeleteAll {
	mda.argRecordTable = typ
	return mda
}

// Error sets error to be returned.
func (mda *MockDeleteAll) Error(err error) *Assert {
	mda.retError = err
	return mda.assert
}

// Success sets no error to be returned.
func (mda *MockDeleteAll) Success() *Assert {
	return mda.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (mda *MockDeleteAll) ConnectionClosed() *Assert {
	return mda.Error(ErrConnectionClosed)
}

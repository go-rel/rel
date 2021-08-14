package reltest

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type deleteAll []*MockDeleteAll

func (da *deleteAll) register(ctxData ctxData, options ...rel.Cascade) *MockDeleteAll {
	mda := &MockDeleteAll{ctxData: ctxData, argOptions: options}
	*da = append(*da, mda)
	return mda
}

func (da deleteAll) execute(ctx context.Context, record interface{}) error {
	for _, mda := range da {
		if fetchContext(ctx) == mda.ctxData &&
			(mda.argRecord == nil || reflect.DeepEqual(mda.argRecord, record)) &&
			(mda.argRecordType == "" || mda.argRecordType == reflect.TypeOf(record).String()) &&
			(mda.argRecordTable == "" || mda.argRecordTable == rel.NewDocument(record, true).Table()) {
			return mda.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockDeleteAll asserts and simulate Delete function for test.
type MockDeleteAll struct {
	ctxData        ctxData
	argRecord      interface{}
	argRecordType  string
	argRecordTable string
	argOptions     []rel.Cascade
	retError       error
}

// For expect calls for given record.
func (mda *MockDeleteAll) For(record interface{}) *MockDeleteAll {
	mda.argRecord = record
	return mda
}

// ForType expect calls for given type.
// Type must include package name, example: `model.User`.
func (mda *MockDeleteAll) ForType(typ string) *MockDeleteAll {
	mda.argRecordType = "*" + strings.TrimPrefix(typ, "*")
	return mda
}

// ForTable expect calls for given table.
func (mda *MockDeleteAll) ForTable(typ string) *MockDeleteAll {
	mda.argRecordTable = typ
	return mda
}

// Error sets error to be returned.
func (mda *MockDeleteAll) Error(err error) {
	mda.retError = err
}

// ConnectionClosed sets this error to be returned.
func (mda *MockDeleteAll) ConnectionClosed() {
	mda.Error(ErrConnectionClosed)
}

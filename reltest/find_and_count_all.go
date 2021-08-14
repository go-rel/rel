package reltest

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-rel/rel"
)

type findAndCountAll []*MockFindAndCountAll

func (fca *findAndCountAll) register(ctxData ctxData, queriers ...rel.Querier) *MockFindAndCountAll {
	mfca := &MockFindAndCountAll{
		assert:   &Assert{ctxData: ctxData},
		argQuery: rel.Build("", queriers...),
	}
	*fca = append(*fca, mfca)
	return mfca
}

func (fca findAndCountAll) execute(ctx context.Context, records interface{}, queriers ...rel.Querier) (int, error) {
	query := rel.Build("", queriers...)
	for _, mfca := range fca {
		if matchQuery(mfca.argQuery, query) &&
			mfca.assert.call(ctx) {
			if mfca.argRecords != nil {
				reflect.ValueOf(records).Elem().Set(reflect.ValueOf(mfca.argRecords))
			}

			return mfca.retCount, mfca.retError
		}
	}

	panic(failExecuteMessage(MockFindAndCountAll{argQuery: query, argRecords: records}, fca))
}

func (fca *findAndCountAll) assert(t T) bool {
	for _, mfca := range *fca {
		if !mfca.assert.assert(t, mfca) {
			return false
		}
	}

	*fca = nil
	return true
}

// MockFindAndCountAll asserts and simulate find and count all function for test.
type MockFindAndCountAll struct {
	assert     *Assert
	argQuery   rel.Query
	argRecords interface{}
	retCount   int
	retError   error
}

// Result sets the result of this query.
func (mfca *MockFindAndCountAll) Result(result interface{}, count int) *Assert {
	mfca.argQuery.Table = rel.NewCollection(result, true).Table()
	mfca.argRecords = result
	mfca.retCount = count
	return mfca.assert
}

// Error sets error to be returned.
func (mfca *MockFindAndCountAll) Error(err error) *Assert {
	mfca.retError = err
	return mfca.assert
}

// ConnectionClosed sets this error to be returned.
func (mfca *MockFindAndCountAll) ConnectionClosed() *Assert {
	return mfca.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (mfca MockFindAndCountAll) String() string {
	return fmt.Sprintf("FindAndCountAll(ctx, <Any>, %s)", mfca.argQuery)
}

// ExpectString representation of mocked call.
func (mfca MockFindAndCountAll) ExpectString() string {
	return fmt.Sprintf("ExpectFindAndCountAll(%s)", mfca.argQuery)
}

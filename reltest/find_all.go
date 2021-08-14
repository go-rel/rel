package reltest

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-rel/rel"
)

type findAll []*MockFindAll

func (fa *findAll) register(ctxData ctxData, queriers ...rel.Querier) *MockFindAll {
	mfa := &MockFindAll{
		assert:   &Assert{ctxData: ctxData},
		argQuery: rel.Build("", queriers...),
	}
	*fa = append(*fa, mfa)
	return mfa
}

func (fa findAll) execute(ctx context.Context, records interface{}, queriers ...rel.Querier) error {
	query := rel.Build("", queriers...)
	for _, mfa := range fa {
		if matchQuery(mfa.argQuery, query) &&
			mfa.assert.call(ctx) {
			if mfa.argRecords != nil {
				reflect.ValueOf(records).Elem().Set(reflect.ValueOf(mfa.argRecords))
			}

			return mfa.retError
		}
	}

	panic(failExecuteMessage(MockFindAll{argQuery: query, argRecords: records}, fa))
}

func (fa *findAll) assert(t T) bool {
	for _, mfa := range *fa {
		if !mfa.assert.assert(t, mfa) {
			return false
		}
	}

	*fa = nil
	return true
}

// MockFindAll asserts and simulate find all function for test.
type MockFindAll struct {
	assert     *Assert
	argQuery   rel.Query
	argRecords interface{}
	retError   error
}

// Result sets the result of this query.
func (mfa *MockFindAll) Result(result interface{}) *Assert {
	mfa.argQuery.Table = rel.NewCollection(result, true).Table()
	mfa.argRecords = result
	return mfa.assert
}

// Error sets error to be returned.
func (mfa *MockFindAll) Error(err error) *Assert {
	mfa.retError = err
	return mfa.assert
}

// ConnectionClosed sets this error to be returned.
func (mfa *MockFindAll) ConnectionClosed() *Assert {
	return mfa.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (mfa MockFindAll) String() string {
	return fmt.Sprintf("FindAll(ctx, <Any>, %s)", mfa.argQuery)
}

// ExpectString representation of mocked call.
func (mfa MockFindAll) ExpectString() string {
	return fmt.Sprintf("ExpectFindAll(%s)", mfa.argQuery)
}

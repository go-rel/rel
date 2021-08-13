package reltest

import (
	"context"
	"reflect"

	"github.com/go-rel/rel"
)

type findAll []*MockFindAll

func (fa *findAll) register(ctxData ctxData, queriers ...rel.Querier) *MockFindAll {
	mfa := &MockFindAll{ctxData: ctxData, argQuery: rel.Build("", queriers...)}
	*fa = append(*fa, mfa)
	return mfa
}

func (fa findAll) execute(ctx context.Context, records interface{}, queriers ...rel.Querier) error {
	query := rel.Build("", queriers...)
	for _, mfa := range fa {
		if fetchContext(ctx) == mfa.ctxData && matchQuery(mfa.argQuery, query) {
			if mfa.argRecords != nil {
				reflect.ValueOf(records).Elem().Set(reflect.ValueOf(mfa.argRecords))
			}

			return mfa.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockFindAll asserts and simulate find all function for test.
type MockFindAll struct {
	ctxData    ctxData
	argQuery   rel.Query
	argRecords interface{}
	retError   error
}

// Result sets the result of this query.
func (mfa *MockFindAll) Result(result interface{}) {
	mfa.argQuery.Table = rel.NewCollection(result, true).Table()
	mfa.argRecords = result
}

// Error sets error to be returned.
func (mfa *MockFindAll) Error(err error) {
	mfa.retError = err
}

// ConnectionClosed sets this error to be returned.
func (mfa *MockFindAll) ConnectionClosed() {
	mfa.Error(ErrConnectionClosed)
}

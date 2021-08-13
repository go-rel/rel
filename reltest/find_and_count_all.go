package reltest

import (
	"context"
	"reflect"

	"github.com/go-rel/rel"
)

type findAndCountAll []*MockFindAndCountAll

func (fca *findAndCountAll) register(ctxData ctxData, queriers ...rel.Querier) *MockFindAndCountAll {
	mfca := &MockFindAndCountAll{ctxData: ctxData, argQuery: rel.Build("", queriers...)}
	*fca = append(*fca, mfca)
	return mfca
}

func (fca findAndCountAll) execute(ctx context.Context, records interface{}, queriers ...rel.Querier) (int, error) {
	query := rel.Build("", queriers...)
	for _, mfca := range fca {
		if fetchContext(ctx) == mfca.ctxData && matchQuery(mfca.argQuery, query) {
			if mfca.argRecords != nil {
				reflect.ValueOf(records).Elem().Set(reflect.ValueOf(mfca.argRecords))
			}

			return mfca.retCount, mfca.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockFindAndCountAll asserts and simulate find and count all function for test.
type MockFindAndCountAll struct {
	ctxData    ctxData
	argQuery   rel.Query
	argRecords interface{}
	retCount   int
	retError   error
}

// Result sets the result of this query.
func (mfca *MockFindAndCountAll) Result(result interface{}, count int) {
	mfca.argQuery.Table = rel.NewCollection(result, true).Table()
	mfca.argRecords = result
	mfca.retCount = count
}

// Error sets error to be returned.
func (mfca *MockFindAndCountAll) Error(err error) {
	mfca.retError = err
}

// ConnectionClosed sets this error to be returned.
func (mfca *MockFindAndCountAll) ConnectionClosed() {
	mfca.Error(ErrConnectionClosed)
}

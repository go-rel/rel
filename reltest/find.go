package reltest

import (
	"context"
	"reflect"

	"github.com/go-rel/rel"
)

type find []*MockFind

func (f *find) register(ctxData ctxData, queriers ...rel.Querier) *MockFind {
	mf := &MockFind{ctxData: ctxData, argQuery: rel.Build("", queriers...)}
	*f = append(*f, mf)
	return mf
}

func (f find) execute(ctx context.Context, records interface{}, queriers ...rel.Querier) error {
	query := rel.Build("", queriers...)
	for _, mf := range f {
		if fetchContext(ctx) == mf.ctxData &&
			matchQuery(mf.argQuery, query) {
			if mf.argRecord != nil {
				reflect.ValueOf(records).Elem().Set(reflect.ValueOf(mf.argRecord))
			}

			return mf.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockFind asserts and simulate find function for test.
type MockFind struct {
	ctxData   ctxData
	argQuery  rel.Query
	argRecord interface{}
	retError  error
}

// Result sets the result of this query.
func (mf *MockFind) Result(result interface{}) {
	mf.argQuery.Table = rel.NewDocument(result, true).Table()
	mf.argRecord = result
}

// Error sets error to be returned.
func (mf *MockFind) Error(err error) {
	mf.retError = err
}

// ConnectionClosed sets this error to be returned.
func (mf *MockFind) ConnectionClosed() {
	mf.Error(ErrConnectionClosed)
}

// NotFound sets NotFoundError to be returned.
func (mf *MockFind) NotFound() {
	mf.Error(rel.NotFoundError{})
}

// ExpectFind to be called with given field and queries.
func ExpectFind(queriers []rel.Querier) *MockFind {
	return &MockFind{
		argQuery: rel.Build("", queriers...),
	}
}

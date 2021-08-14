package reltest

import (
	"context"

	"github.com/go-rel/rel"
)

type count []*MockCount

func (c *count) register(ctxData ctxData, collection string, queriers ...rel.Querier) *MockCount {
	mc := &MockCount{ctxData: ctxData, argCollection: collection, argQuery: rel.Build(collection, queriers...)}
	*c = append(*c, mc)
	return mc
}

func (e count) execute(ctx context.Context, collection string, queriers ...rel.Querier) (int, error) {
	query := rel.Build(collection, queriers...)
	for _, me := range e {
		if fetchContext(ctx) == me.ctxData &&
			me.argCollection == collection &&
			matchQuery(me.argQuery, query) {
			return me.retCount, me.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockCount asserts and simulate UpdateAny function for test.
type MockCount struct {
	ctxData       ctxData
	argCollection string
	argQuery      rel.Query
	retCount      int
	retError      error
}

// Result sets the result of this query.
func (me *MockCount) Result(count int) {
	me.retCount = count
}

// Error sets error to be returned.
func (me *MockCount) Error(err error) {
	me.retError = err
}

// ConnectionClosed sets this error to be returned.
func (me *MockCount) ConnectionClosed() {
	me.Error(ErrConnectionClosed)
}

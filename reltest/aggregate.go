package reltest

import (
	"context"

	"github.com/go-rel/rel"
)

type aggregate []*MockAggregate

func (a *aggregate) register(ctxData ctxData, query rel.Query, aggregate string, field string) *MockAggregate {
	ma := &MockAggregate{
		assert:       &Assert{ctxData: ctxData},
		argQuery:     query,
		argAggregate: aggregate,
		argField:     field,
	}
	*a = append(*a, ma)
	return ma
}

func (a aggregate) execute(ctx context.Context, query rel.Query, aggregate string, field string) (int, error) {
	for _, ma := range a {
		if matchQuery(ma.argQuery, query) &&
			ma.argAggregate == aggregate &&
			ma.argField == field &&
			ma.assert.call(ctx) {
			return ma.retCount, ma.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockAggregate asserts and simulate UpdateAny function for test.
type MockAggregate struct {
	assert       *Assert
	argQuery     rel.Query
	argAggregate string
	argField     string
	retCount     int
	retError     error
}

// Result sets the result of this query.
func (me *MockAggregate) Result(count int) *Assert {
	me.retCount = count
	return me.assert
}

// Error sets error to be returned.
func (me *MockAggregate) Error(err error) *Assert {
	me.retError = err
	return me.assert
}

// ConnectionClosed sets this error to be returned.
func (me *MockAggregate) ConnectionClosed() *Assert {
	me.Error(ErrConnectionClosed)
	return me.assert
}

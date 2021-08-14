package reltest

import (
	"context"
	"fmt"
	"testing"

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

	ma := MockAggregate{argQuery: query, argAggregate: aggregate, argField: field}
	mocks := ""
	for i := range a {
		mocks += "\n\t" + a[i].ExpectString()
	}
	panic(fmt.Sprintf("FAIL: this call is not mocked:\n\t%s\nMaybe try adding mock:\t\n%s\n\nAvailable mocks:%s", ma, ma.ExpectString(), mocks))
}

func (a aggregate) assert(t *testing.T) bool {
	for _, ma := range a {
		if !ma.assert.assert() {
			t.Errorf("FAIL: The code you are testing needs to make %d more call(s):\n\t%s", ma.assert.repeatability-ma.assert.totalCalls, ma)
			return false
		}
	}

	return true
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
func (ma *MockAggregate) Result(count int) *Assert {
	ma.retCount = count
	return ma.assert
}

// Error sets error to be returned.
func (ma *MockAggregate) Error(err error) *Assert {
	ma.retError = err
	return ma.assert
}

// ConnectionClosed sets this error to be returned.
func (ma *MockAggregate) ConnectionClosed() *Assert {
	ma.Error(ErrConnectionClosed)
	return ma.assert
}

// String representation of mocked call.
func (ma MockAggregate) String() string {
	return fmt.Sprintf("Aggregate(ctx, %s, %s, %s)", ma.argQuery, ma.argAggregate, ma.argField)
}

// ExpectString representation of mocked call.
func (ma MockAggregate) ExpectString() string {
	return fmt.Sprintf("ExpectAggregate(%s, %s, %s)", ma.argQuery, ma.argAggregate, ma.argField)
}

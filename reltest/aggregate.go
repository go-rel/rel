package reltest

import "github.com/Fs02/rel"

// Aggregate asserts and simulate aggregate function for test.
type Aggregate struct {
	*Expect
}

// Result sets the result of this query.
func (a *Aggregate) Result(count int) {
	a.Return(count, nil)
}

// Error sets error to be returned.
func (a *Aggregate) Error(err error) {
	a.Return(0, err)
}

// ConnectionClosed sets this error to be returned.
func (a *Aggregate) ConnectionClosed() {
	a.Error(ErrConnectionClosed)
}

// ExpectAggregate to be called with given field and queries.
func ExpectAggregate(r *Repository, query rel.Query, aggregate string, field string) *Aggregate {
	return &Aggregate{
		Expect: newExpect(r, "Aggregate",
			[]interface{}{r.ctxData, query, aggregate, field},
			[]interface{}{0, nil},
		),
	}
}

// ExpectCount to be called with given field and queries.
func ExpectCount(r *Repository, collection string, queriers []rel.Querier) *Aggregate {
	return &Aggregate{
		Expect: newExpect(r, "Count",
			[]interface{}{r.ctxData, collection, queriers},
			[]interface{}{0, nil},
		),
	}
}

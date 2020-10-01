package reltest

import (
	"github.com/go-rel/rel"
	"github.com/stretchr/testify/mock"
)

// Find asserts and simulate find function for test.
type Find struct {
	*FindAll
}

// NotFound sets NotFoundError to be returned.
func (f *Find) NotFound() {
	f.Error(rel.NotFoundError{})
}

// ExpectFind to be called with given field and queries.
func ExpectFind(r *Repository, queriers []rel.Querier) *Find {
	return &Find{
		FindAll: &FindAll{
			Expect: newExpect(r, "Find",
				[]interface{}{r.ctxData, mock.Anything, queriers},
				[]interface{}{nil},
			),
		},
	}
}

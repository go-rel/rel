package reltest

import (
	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

// Find asserts and simulate find function for test.
type Find struct {
	*FindAll
}

// NoResult sets NoResultError to be returned.
func (f *Find) NoResult() {
	f.Error(rel.NoResultError{})
}

// ExpectFind to be called with given field and queries.
func ExpectFind(r *Repository, queriers []rel.Querier) *Find {
	return &Find{
		FindAll: &FindAll{
			Expect: newExpect(r, "Find",
				[]interface{}{mock.Anything, queriers},
				[]interface{}{nil},
			),
		},
	}
}

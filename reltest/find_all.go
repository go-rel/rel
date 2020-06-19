package reltest

import (
	"fmt"
	"reflect"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

// FindAll asserts and simulate find all function for test.
type FindAll struct {
	*Expect
}

// Result sets the result of this query.
func (fa *FindAll) Result(records interface{}) {
	fa.Arguments[1] = mock.AnythingOfType(fmt.Sprintf("*%T", records))

	fa.Run(func(args mock.Arguments) {
		reflect.ValueOf(args[1]).Elem().Set(reflect.ValueOf(records))
	})
}

// ExpectFindAll to be called with given field and queries.
func ExpectFindAll(r *Repository, queriers []rel.Querier) *FindAll {
	return &FindAll{
		Expect: newExpect(r, "FindAll",
			[]interface{}{r.ctxData, mock.Anything, queriers},
			[]interface{}{nil},
		),
	}
}

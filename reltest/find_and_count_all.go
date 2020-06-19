package reltest

import (
	"fmt"
	"reflect"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

// FindAndCountAll asserts and simulate find all function for test.
type FindAndCountAll struct {
	*Expect
}

// Result sets the result of this query.
func (fa *FindAndCountAll) Result(records interface{}, count int) {
	fa.Arguments[1] = mock.AnythingOfType(fmt.Sprintf("*%T", records))

	fa.Run(func(args mock.Arguments) {
		reflect.ValueOf(args[1]).Elem().Set(reflect.ValueOf(records))
	}).Return(count, nil)
}

// Error sets error to be returned.
func (fa *FindAndCountAll) Error(err error) {
	fa.Return(0, err)
}

// ConnectionClosed sets this error to be returned.
func (fa *FindAndCountAll) ConnectionClosed() {
	fa.Error(ErrConnectionClosed)
}

// ExpectFindAndCountAll to be called with given field and queries.
func ExpectFindAndCountAll(r *Repository, queriers []rel.Querier) *FindAndCountAll {
	return &FindAndCountAll{
		Expect: newExpect(r, "FindAndCountAll",
			[]interface{}{r.ctxData, mock.Anything, queriers},
			[]interface{}{0, nil},
		),
	}
}

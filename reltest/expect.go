package reltest

import (
	"database/sql"
	"reflect"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

type ExpectFindAll struct {
	call *mock.Call
}

// Result sets the result of Find query.
func (efa *ExpectFindAll) Result(record interface{}) {
	// TODO: mock anything of type
	efa.call.Return(func(out interface{}, queriers ...rel.Querier) error {
		// TODO: check type
		reflect.ValueOf(out).Elem().Set(reflect.ValueOf(record))
		return nil
	}).Once()
}

// Error sets error to be returned by Find Query.
func (efa *ExpectFindAll) Error(err error) {
	efa.call.Return(err).Once()
}

func (efa *ExpectFindAll) ConnectionClosed() {
	efa.Error(sql.ErrConnDone)
}

func NewExpectFindAll(r *Repository, queriers []rel.Querier) *ExpectFindAll {
	return &ExpectFindAll{
		call: r.On("FindAll", querierArgs(queriers)...).Return(nil),
	}
}

type ExpectFind struct {
	*ExpectFindAll
}

// NoResult sets NoResultError to be returned by Find Query.
func (ef ExpectFind) NoResult() {
	ef.Error(rel.NoResultError{})
}

func NewExpectFind(r *Repository, queriers []rel.Querier) ExpectFind {
	return ExpectFind{
		ExpectFindAll: &ExpectFindAll{
			call: r.On("Find", querierArgs(queriers)...).Return(nil),
		},
	}
}

func querierArgs(queriers []rel.Querier) []interface{} {
	args := make([]interface{}, len(queriers)+1)
	args[0] = mock.Anything // first records argument
	for i := range queriers {
		args[i+1] = queriers[i]
	}

	return args
}

package reltest

import (
	"reflect"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

type ExpectFindAll struct {
	call *mock.Call
}

// Result sets the result of Find query.
func (efa *ExpectFindAll) Result(record interface{}) {
	efa.call.Return(func(out interface{}, queriers ...rel.Querier) error {
		// TODO: check type
		reflect.ValueOf(out).Elem().Set(reflect.ValueOf(record))
		return nil
	})
}

// Error sets error to be returned by Find Query.
func (efa *ExpectFindAll) Error(err error) {
	efa.call.Return(err)
}

func NewExpectFindAll(r *Repository, queriers []rel.Querier) *ExpectFindAll {
	return &ExpectFindAll{
		call: r.On("FindAll", querierArgs(queriers)),
	}
}

type ExpectFind ExpectFindAll

// NoResult sets NoResultError to be returned by Find Query.
func (ef *ExpectFind) NoResult() {
	ef.call.Return(rel.NoResultError{})
}

func NewExpectFind(r *Repository, queriers []rel.Querier) *ExpectFind {
	return &ExpectFind{
		call: r.On("Find", querierArgs(queriers)),
	}
}

func querierArgs(queriers []rel.Querier) []interface{} {
	args := make([]interface{}, len(queriers)+1)
	args[0] = mock.Anything
	for i := range queriers {
		args[i+1] = queriers[i]
	}

	return args
}

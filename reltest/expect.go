package reltest

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

type Expect struct {
	*mock.Call
}

// Error sets error to be returned by Find Query.
func (e *Expect) Error(err error) {
	e.Return(err)
}

func (e *Expect) ConnectionClosed() {
	e.Error(sql.ErrConnDone)
}

func newExpect(r *Repository, methodName string, args []interface{}, rets []interface{}) *Expect {
	return &Expect{
		Call: r.On(methodName, args...).Return(rets...).Once(),
	}
}

type ExpectAggregate struct {
	*Expect
}

func (ea *ExpectAggregate) Result(count int) {
	ea.Return(func(query rel.Query, aggregate string, field string) int {
		return count
	}, nil)
}

// Error sets error to be returned by Find Query.
func (ea *ExpectAggregate) Error(err error) {
	ea.Return(0, err)
}

func (ea *ExpectAggregate) ConnectionClosed() {
	ea.Error(sql.ErrConnDone)
}

func NewExpectAggregate(r *Repository, query rel.Query, aggregate string, field string) *ExpectAggregate {
	return &ExpectAggregate{
		Expect: newExpect(r, "Aggregate",
			[]interface{}{query, aggregate, field},
			[]interface{}{0, nil},
		),
	}
}

type ExpectFindAll struct {
	*Expect
}

// Result sets the result of Find query.
func (efa *ExpectFindAll) Result(records interface{}) {
	// adjust arguments
	efa.Arguments[0] = mock.AnythingOfType(fmt.Sprintf("*%T", records))

	efa.Return(func(out interface{}, queriers ...rel.Querier) error {
		reflect.ValueOf(out).Elem().Set(reflect.ValueOf(records))
		return nil
	})
}

func newExpectFindAll(r *Repository, methodName string, queriers []rel.Querier) *ExpectFindAll {
	return &ExpectFindAll{
		Expect: newExpect(r, methodName,
			[]interface{}{mock.Anything, rel.BuildQuery("", queriers...)},
			[]interface{}{nil},
		),
	}
}

func NewExpectFindAll(r *Repository, queriers []rel.Querier) *ExpectFindAll {
	return newExpectFindAll(r, "FindAll", queriers)
}

type ExpectFind struct {
	*ExpectFindAll
}

// NoResult sets NoResultError to be returned by Find Query.
func (ef *ExpectFind) NoResult() {
	ef.Error(rel.NoResultError{})
}

func NewExpectFind(r *Repository, queriers []rel.Querier) *ExpectFind {
	return &ExpectFind{
		ExpectFindAll: newExpectFindAll(r, "Find", queriers),
	}
}

type ExpectDelete struct {
	*Expect
}

func (ed *ExpectDelete) Record(record interface{}) {
	// adjust arguments
	ed.Arguments[0] = record
}

func NewExpectDelete(r *Repository) *ExpectDelete {
	return &ExpectDelete{
		Expect: newExpect(r, "Delete", []interface{}{mock.Anything}, []interface{}{nil}),
	}
}

type ExpectDeleteAll struct {
	*Expect
}

// Unsafe allows for unsafe delete that doesn't contains where cloause.
func (eda *ExpectDeleteAll) Unsafe() {
	eda.RunFn = nil // clear validation
}

func NewExpectDeleteAll(r *Repository, queriers []rel.Querier) *ExpectDeleteAll {
	query := rel.BuildQuery("", queriers...)
	if query.Collection == "" {
		panic("reltest: cannot call DeleteAll without specifying table name. use rel.From(tableName)")
	}

	eda := &ExpectDeleteAll{
		Expect: newExpect(r, "DeleteAll",
			[]interface{}{query},
			[]interface{}{nil},
		),
	}

	// validation
	eda.Run(func(args mock.Arguments) {
		query := args[0].(rel.Query)

		if query.WhereQuery.None() {
			panic("reltest: unsafe delete all detected. if you want to delete all records without filter, please use ExpectDeleteAll().Unsafe()")
		}
	})

	return eda
}

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

// Error sets error to be returned by Find query.
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
	ea.Return(count, nil)
}

// Error sets error to be returned by Find query.
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

func NewExpectAggregateCount(r *Repository, collection string, queriers []rel.Querier) *ExpectAggregate {
	return &ExpectAggregate{
		Expect: newExpect(r, "Count",
			[]interface{}{collection, queriers},
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

	efa.Run(func(args mock.Arguments) {
		reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(records))
	})
}

func newExpectFindAll(r *Repository, methodName string, queriers []rel.Querier) *ExpectFindAll {
	return &ExpectFindAll{
		Expect: newExpect(r, methodName,
			[]interface{}{mock.Anything, queriers},
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

// NoResult sets NoResultError to be returned by Find query.
func (ef *ExpectFind) NoResult() {
	ef.Error(rel.NoResultError{})
}

func NewExpectFind(r *Repository, queriers []rel.Querier) *ExpectFind {
	return &ExpectFind{
		ExpectFindAll: newExpectFindAll(r, "Find", queriers),
	}
}

type ExpectModify struct {
	*Expect
}

func (em *ExpectModify) Record(record interface{}) {
	// adjust arguments
	em.Arguments[0] = record
}

func applyChanges(rv reflect.Value, changes rel.Changes) {
	var (
		doc   = rel.NewDocument(rv)
		index = doc.Index()
	)

	// TODO: check update without id

	rv = rv.Elem()

	// TODO: for insertion, always set id to 1.
	pv := rv.Field(index[doc.PrimaryField()])
	if pv.Kind() == reflect.Int && pv.Int() == 0 {
		pv.SetInt(1)
	}

	if changes.Count() == 0 {
		return
	}

	for _, ch := range changes.All() {
		if i, ok := index[ch.Field]; ok {
			// TODO: other types
			switch ch.Type {
			case rel.ChangeSetOp:
				rv.Field(i).Set(reflect.ValueOf(ch.Value))
			}
		} else {
			panic("reltest: cannot apply changes, field " + ch.Field + " is not defined")
		}
	}

	if changes.AssocCount() == 0 {
		return
	}

	for _, field := range doc.BelongsTo() {
		if assoc, ok := changes.GetAssoc(field); ok {
			applyChanges(rv.Field(index[field]).Addr(), assoc.Changes[0])
		}
	}

	for _, field := range doc.HasOne() {
		if assoc, ok := changes.GetAssoc(field); ok {
			applyChanges(rv.Field(index[field]).Addr(), assoc.Changes[0])
		}
	}

	for _, field := range doc.HasMany() {
		if assoc, ok := changes.GetAssoc(field); ok {
			var (
				fv    = rv.Field(index[field])
				elTyp = fv.Type().Elem()
			)

			fv.Set(reflect.MakeSlice(fv.Type(), 0, len(assoc.Changes)))

			for _, ch := range assoc.Changes {
				el := reflect.New(elTyp)
				applyChanges(el, ch)

				fv.Set(reflect.Append(fv, el.Elem()))
			}
		}
	}
}

func NewExpectModify(r *Repository, methodName string, changers []rel.Changer) *ExpectModify {
	em := &ExpectModify{
		Expect: newExpect(r, methodName,
			[]interface{}{mock.Anything, changers},
			[]interface{}{nil},
		),
	}

	em.Run(func(args mock.Arguments) {
		changes := rel.BuildChanges(args[1].([]rel.Changer)...)
		applyChanges(reflect.ValueOf(args[0]), changes)
	})

	return em
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

// Unsafe allows for unsafe delete that doesn't contains where clause.
func (eda *ExpectDeleteAll) Unsafe() {
	eda.RunFn = nil // clear validation
}

func NewExpectDeleteAll(r *Repository, queriers []rel.Querier) *ExpectDeleteAll {
	eda := &ExpectDeleteAll{
		Expect: newExpect(r, "DeleteAll",
			[]interface{}{queriers},
			[]interface{}{nil},
		),
	}

	// validation
	eda.Run(func(args mock.Arguments) {
		query := rel.BuildQuery("", args[0].([]rel.Querier)...)

		if query.Collection == "" {
			panic("reltest: cannot call DeleteAll without specifying table name. use rel.From(tableName)")
		}

		if query.WhereQuery.None() {
			panic("reltest: unsafe delete all detected. if you want to delete all records without filter, please use ExpectDeleteAll().Unsafe()")
		}
	})

	return eda
}

type ExpectPreload struct {
	*Expect
}

// Result sets the result of Find query.
func (ep *ExpectPreload) Result(records interface{}) {
	// adjust arguments
	ep.Arguments[0] = mock.AnythingOfType(fmt.Sprintf("*%T", records))

	ep.Run(func(args mock.Arguments) {
		reflect.ValueOf(args[0]).Elem().Set(reflect.ValueOf(records))
	})
}

func NewExpectPreload(r *Repository, field string, queriers []rel.Querier) *ExpectPreload {
	return &ExpectPreload{
		Expect: newExpect(r, "Preload",
			[]interface{}{mock.Anything, field, queriers},
			[]interface{}{nil},
		),
	}
}

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
		Call: r.mock.On(methodName, args...).Return(rets...).Once(),
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

func newExpectAggregate(r *Repository, query rel.Query, aggregate string, field string) *ExpectAggregate {
	return &ExpectAggregate{
		Expect: newExpect(r, "Aggregate",
			[]interface{}{query, aggregate, field},
			[]interface{}{0, nil},
		),
	}
}

func newExpectAggregateCount(r *Repository, collection string, queriers []rel.Querier) *ExpectAggregate {
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

func newExpectFindAll(r *Repository, queriers []rel.Querier) *ExpectFindAll {
	return &ExpectFindAll{
		Expect: newExpect(r, "FindAll",
			[]interface{}{mock.Anything, queriers},
			[]interface{}{nil},
		),
	}
}

func newExpectPreload(r *Repository, field string, queriers []rel.Querier) *ExpectFindAll {
	return &ExpectFindAll{
		Expect: newExpect(r, "Preload",
			[]interface{}{mock.Anything, field, queriers},
			[]interface{}{nil},
		),
	}
}

type ExpectFind struct {
	*ExpectFindAll
}

// NoResult sets NoResultError to be returned by Find query.
func (ef *ExpectFind) NoResult() {
	ef.Error(rel.NoResultError{})
}

func newExpectFind(r *Repository, queriers []rel.Querier) *ExpectFind {
	return &ExpectFind{
		ExpectFindAll: &ExpectFindAll{
			Expect: newExpect(r, "Find",
				[]interface{}{mock.Anything, queriers},
				[]interface{}{nil},
			),
		},
	}
}

type ExpectModify struct {
	*Expect
}

func (em *ExpectModify) For(record interface{}) {
	// adjust arguments
	em.Arguments[0] = record
}

func (em *ExpectModify) NotUnique(key string) {
	em.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

func applyChanges(doc *rel.Document, changes rel.Changes, insertion bool, includeAssoc bool) {
	if doc.PrimaryValue() == 0 {
		if insertion {
			doc.SetValue(doc.PrimaryField(), 1)
		} else {
			panic("reltest: cannot update a record without using primary key")
		}
	}

	if includeAssoc {
		applyBelongsToChanges(doc, &changes)
	}

	for _, ch := range changes.All() {
		if !doc.SetValue(ch.Field, ch.Value) {
			panic("reltest: cannot apply changes, field " + ch.Field + " is not defined")
		}
	}

	if includeAssoc {
		applyHasOneChanges(doc, &changes)
		applyHasManyChanges(doc, &changes, insertion)
	}
}

// this logic should be similar to repository.saveBelongsTo
func applyBelongsToChanges(doc *rel.Document, changes *rel.Changes) {
	for _, field := range doc.BelongsTo() {
		ac, changed := changes.GetAssoc(field)
		if !changed || len(ac.Changes) == 0 {
			continue
		}

		var (
			assocChanges = ac.Changes[0]
			assoc        = doc.Association(field)
			fValue       = assoc.ForeignValue()
			doc, loaded  = assoc.Document()
		)

		if loaded {
			// update
			var (
				pField = doc.PrimaryField()
				pValue = doc.PrimaryValue()
			)

			if pch, exist := assocChanges.Get(pField); exist && pch.Value != pValue {
				panic("reltest: inconsistent primary value of belongs to assoc")
			}

			if assoc.ReferenceValue() != fValue {
				panic("reltest: inconsistent referenced foreign key of belongs to assoc")
			}

			applyChanges(doc, assocChanges, false, true)
		} else {
			// insert
			applyChanges(doc, assocChanges, true, true)

			changes.SetValue(assoc.ReferenceField(), assoc.ForeignValue())
		}
	}
}

// This logic should be similar to repository.saveHasOne
func applyHasOneChanges(doc *rel.Document, changes *rel.Changes) {
	for _, field := range doc.HasOne() {
		ac, changed := changes.GetAssoc(field)
		if !changed || len(ac.Changes) == 0 {
			continue
		}

		var (
			assocChanges = ac.Changes[0]
			assoc        = doc.Association(field)
			fField       = assoc.ForeignField()
			rValue       = assoc.ReferenceValue()
			doc, loaded  = assoc.Document()
			pField       = doc.PrimaryField()
			pValue       = doc.PrimaryValue()
		)

		if loaded {
			// update
			if pch, exist := assocChanges.Get(pField); exist && pch.Value != pValue {
				panic("cannot update assoc: inconsistent primary key")
			}

			if assoc.ForeignValue() != rValue {
				panic("reltest: inconsistent referenced foreign key of has one assoc")
			}

			applyChanges(doc, assocChanges, false, true)
		} else {
			// insert
			assocChanges.SetValue(fField, rValue)

			applyChanges(doc, assocChanges, true, true)
		}
	}
}

func applyHasManyChanges(doc *rel.Document, changes *rel.Changes, insertion bool) {
	for _, field := range doc.HasMany() {
		ac, changed := changes.GetAssoc(field)
		if !changed {
			continue
		}

		var (
			assoc       = doc.Association(field)
			col, loaded = assoc.Collection()
			pField      = col.PrimaryField()
			fField      = assoc.ForeignField()
			rValue      = assoc.ReferenceValue()
			pIndex      = make(map[interface{}]int)
			pValues     = col.PrimaryValue().([]interface{})
		)

		for i, v := range pValues {
			pIndex[v] = i
		}

		if !insertion && !loaded {
			panic("rel: association must be loaded to update")
		}

		var (
			curr    = 0
			inserts []rel.Changes
		)

		// update
		for _, ch := range ac.Changes {
			if pChange, changed := ch.Get(pField); changed {
				// update
				pId, ok := pIndex[pChange.Value]
				if !ok {
					panic("reltest: cannot update has many assoc that is not loaded or doesn't belong to this record")
				}

				if pId != curr {
					col.Swap(pId, curr)
					pValues[pId], pValues[curr] = pValues[curr], pValues[pId]
				}

				var (
					doc       = col.Get(curr)
					fValue, _ = doc.Value(fField)
				)

				if fValue != rValue {
					panic("reltest: inconsistent foreign key when updating has many")
				}

				applyChanges(doc, ch, false, true)

				delete(pIndex, pChange.Value)
				curr++
			} else {
				inserts = append(inserts, ch)
			}
		}

		// delete stales
		if curr < col.Len() {
			col.Truncate(0, curr)
		}

		// inserts remaining
		for _, ch := range inserts {
			ch.SetValue(fField, rValue)

			applyChanges(col.Add(), ch, true, true)
		}
	}
}

func newExpectModify(r *Repository, methodName string, changers []rel.Changer, insertion bool) *ExpectModify {
	em := &ExpectModify{
		Expect: newExpect(r, methodName,
			[]interface{}{mock.Anything, changers},
			[]interface{}{nil},
		),
	}

	em.Run(func(args mock.Arguments) {
		var (
			changes  rel.Changes
			changers = args[1].([]rel.Changer)
		)

		if len(changers) == 0 {
			changes = rel.BuildChanges(rel.NewStructset(args[0]))
		} else {
			changes = rel.BuildChanges(changers...)
		}

		applyChanges(rel.NewDocument(args[0]), changes, insertion, true)
	})

	return em
}

func newExpectInsertAll(r *Repository, changes []rel.Changes) *ExpectModify {
	em := &ExpectModify{
		Expect: newExpect(r, "InsertAll",
			[]interface{}{mock.Anything, changes},
			[]interface{}{nil},
		),
	}

	em.Run(func(args mock.Arguments) {
		var (
			records = args[0]
			col     = rel.NewCollection(records)
		)

		if len(changes) == 0 {
			// just set primary keys
			for i := 0; i < col.Len(); i++ {
				doc := col.Get(i)
				doc.SetValue(doc.PrimaryField(), 1)
			}
		} else {
			col.Reset()

			for i := range changes {
				doc := col.Add()
				applyChanges(doc, changes[i], true, false)
			}
		}
	})

	return em
}

type ExpectDelete struct {
	*Expect
}

func (ed *ExpectDelete) For(record interface{}) {
	// adjust arguments
	ed.Arguments[0] = record
}

func newExpectDelete(r *Repository) *ExpectDelete {
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

func newExpectDeleteAll(r *Repository, queriers []rel.Querier) *ExpectDeleteAll {
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

package reltest

import (
	"strings"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

// Modify asserts and simulate insert or update function for test.
type Modify struct {
	*Expect
}

// For match expect calls for given record.
func (m *Modify) For(record interface{}) *Modify {
	m.Arguments[0] = record
	return m
}

// ForType match expect calls for given type.
// Type must include package name, example: `model.User`.
func (m *Modify) ForType(typ string) *Modify {
	return m.For(mock.AnythingOfType("*" + strings.TrimPrefix(typ, "*")))
}

// NotUnique sets not unique error to be returned.
func (m *Modify) NotUnique(key string) {
	m.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

// ExpectModify to be called with given field and queries.
func ExpectModify(r *Repository, methodName string, changers []rel.Changer, insertion bool) *Modify {
	em := &Modify{
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

// ExpectInsertAll to be called with given field and queries.
func ExpectInsertAll(r *Repository, changes []rel.Changes) *Modify {
	em := &Modify{
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
				pID, ok := pIndex[pChange.Value]
				if !ok {
					panic("reltest: cannot update has many assoc that is not loaded or doesn't belong to this record")
				}

				if pID != curr {
					col.Swap(pID, curr)
					pValues[pID], pValues[curr] = pValues[curr], pValues[pID]
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

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
func ExpectModify(r *Repository, methodName string, modifiers []rel.Modifier, insertion bool) *Modify {
	em := &Modify{
		Expect: newExpect(r, methodName,
			[]interface{}{mock.Anything, modifiers},
			[]interface{}{nil},
		),
	}

	em.Run(func(args mock.Arguments) {
		var (
			modification rel.Modification
			modifiers    = args[1].([]rel.Modifier)
		)

		if len(modifiers) == 0 {
			modification = rel.BuildModification(rel.NewStructset(args[0]))
		} else {
			modification = rel.BuildModification(modifiers...)
		}

		applyModification(rel.NewDocument(args[0]), modification, insertion, true)
	})

	return em
}

// ExpectInsertAll to be called with given field and queries.
func ExpectInsertAll(r *Repository, modification []rel.Modification) *Modify {
	em := &Modify{
		Expect: newExpect(r, "InsertAll",
			[]interface{}{mock.Anything, modification},
			[]interface{}{nil},
		),
	}

	em.Run(func(args mock.Arguments) {
		var (
			records = args[0]
			col     = rel.NewCollection(records)
		)

		if len(modification) == 0 {
			// just set primary keys
			for i := 0; i < col.Len(); i++ {
				doc := col.Get(i)
				doc.SetValue(doc.PrimaryField(), 1)
			}
		} else {
			col.Reset()

			for i := range modification {
				doc := col.Add()
				applyModification(doc, modification[i], true, false)
			}
		}
	})

	return em
}

func applyModification(doc *rel.Document, modification rel.Modification, insertion bool, includeAssoc bool) {
	if doc.PrimaryValue() == 0 {
		if insertion {
			// FIXME: support other field types.
			doc.SetValue(doc.PrimaryField(), 1)
		} else {
			panic("reltest: cannot update a record without using primary key")
		}
	}

	if includeAssoc {
		applyBelongsToModification(doc, &modification)
	}

	for _, mod := range modification.All() {
		switch mod.Type {
		case rel.ChangeSetOp:
			if !doc.SetValue(mod.Field, mod.Value) {
				panic("reltest: cannot apply modification, field " + mod.Field + " is not defined or not assignable")
			}
		case rel.ChangeIncOp:
			if insertion {
				panic("reltest: increment is not supported for insertion")
			}

			applyIncOrDec(doc, mod.Field, mod.Value.(int))
		case rel.ChangeDecOp:
			if insertion {
				panic("reltest: decrement is not supported for insertion")
			}

			applyIncOrDec(doc, mod.Field, -mod.Value.(int))
		}
	}

	if includeAssoc {
		applyHasOneModification(doc, &modification)
		applyHasManyModification(doc, &modification, insertion)
	}
}

func applyIncOrDec(doc *rel.Document, field string, n int) {
	v, ok := doc.Value(field)
	if !ok {
		panic("reltest: " + field + " field doesn't exists")
	}

	vi, ok := v.(int)
	if !ok {
		panic("reltest: can't inc/dec " + field + " field")
	}

	doc.SetValue(field, vi+n)
}

// this logic should be similar to repository.saveBelongsTo
func applyBelongsToModification(doc *rel.Document, modification *rel.Modification) {
	for _, field := range doc.BelongsTo() {
		ac, changed := modification.GetAssoc(field)
		if !changed || len(ac.Modification) == 0 {
			continue
		}

		var (
			assocModification = ac.Modification[0]
			assoc             = doc.Association(field)
			fValue            = assoc.ForeignValue()
			doc, loaded       = assoc.Document()
		)

		if loaded {
			// update
			var (
				pField = doc.PrimaryField()
				pValue = doc.PrimaryValue()
			)

			if pch, exist := assocModification.Get(pField); exist && pch.Value != pValue {
				panic("reltest: inconsistent primary value of belongs to assoc")
			}

			if assoc.ReferenceValue() != fValue {
				panic("reltest: inconsistent referenced foreign key of belongs to assoc")
			}

			applyModification(doc, assocModification, false, true)
		} else {
			// insert
			applyModification(doc, assocModification, true, true)

			modification.SetValue(assoc.ReferenceField(), assoc.ForeignValue())
		}
	}
}

// This logic should be similar to repository.saveHasOne
func applyHasOneModification(doc *rel.Document, modification *rel.Modification) {
	for _, field := range doc.HasOne() {
		ac, changed := modification.GetAssoc(field)
		if !changed || len(ac.Modification) == 0 {
			continue
		}

		var (
			assocModification = ac.Modification[0]
			assoc             = doc.Association(field)
			fField            = assoc.ForeignField()
			rValue            = assoc.ReferenceValue()
			doc, loaded       = assoc.Document()
			pField            = doc.PrimaryField()
			pValue            = doc.PrimaryValue()
		)

		if loaded {
			// update
			if pch, exist := assocModification.Get(pField); exist && pch.Value != pValue {
				panic("cannot update assoc: inconsistent primary key")
			}

			if assoc.ForeignValue() != rValue {
				panic("reltest: inconsistent referenced foreign key of has one assoc")
			}

			applyModification(doc, assocModification, false, true)
		} else {
			// insert
			assocModification.SetValue(fField, rValue)

			applyModification(doc, assocModification, true, true)
		}
	}
}

func applyHasManyModification(doc *rel.Document, modification *rel.Modification, insertion bool) {
	for _, field := range doc.HasMany() {
		ac, changed := modification.GetAssoc(field)
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
			inserts []rel.Modification
		)

		// update
		for _, mod := range ac.Modification {
			if pChange, changed := mod.Get(pField); changed {
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

				applyModification(doc, mod, false, true)

				delete(pIndex, pChange.Value)
				curr++
			} else {
				inserts = append(inserts, mod)
			}
		}

		// delete stales
		if curr < col.Len() {
			col.Truncate(0, curr)
		}

		// inserts remaining
		for _, mod := range inserts {
			mod.SetValue(fField, rValue)

			applyModification(col.Add(), mod, true, true)
		}
	}
}

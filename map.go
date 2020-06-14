package rel

import (
	"fmt"
)

// Map can be used as mutation for repository insert or update operation.
// This allows inserting or updating only on specified field.
// Insert/Update of has one or belongs to can be done using other Map as a value.
// Insert/Update of has many can be done using slice of Map as a value.
// Map is intended to be used internally within application, and not to be exposed directly as an APIs.
type Map map[string]interface{}

// Apply mutation.
func (m Map) Apply(doc *Document, mutation *Mutation) {
	var (
		pField = doc.PrimaryField()
		pValue = doc.PrimaryValue()
	)

	for field, value := range m {
		switch v := value.(type) {
		case Map:
			if !mutation.Cascade {
				continue
			}

			var (
				assoc = doc.Association(field)
			)

			if assoc.Type() != HasOne && assoc.Type() != BelongsTo {
				panic(fmt.Sprint("rel: cannot associate has many", v, "as", field, "into", doc.Table()))
			}

			var (
				assocDoc, _   = assoc.Document()
				assocMutation = Apply(assocDoc, v)
			)

			mutation.SetAssoc(field, assocMutation)
		case []Map:
			if !mutation.Cascade {
				continue
			}
			var (
				assoc            = doc.Association(field)
				muts, deletedIDs = applyMaps(v, assoc)
			)

			mutation.SetAssoc(field, muts...)
			mutation.SetDeletedIDs(field, deletedIDs)
		default:
			if field == pField {
				if v != pValue {
					panic(fmt.Sprint("rel: replacing primary value (", pValue, " become ", v, ") is not allowed"))
				} else {
					continue
				}
			}

			if !doc.SetValue(field, v) {
				panic(fmt.Sprint("rel: cannot assign ", v, " as ", field, " into ", doc.Table()))
			}

			mutation.Add(Set(field, v))
		}
	}
}

func applyMaps(maps []Map, assoc Association) ([]Mutation, []interface{}) {
	var (
		deletedIDs []interface{}
		muts       = make([]Mutation, len(maps))
		col, _     = assoc.Collection()
		pField     = col.PrimaryField()
		pIndex     = make(map[interface{}]int)
		pValues    = col.PrimaryValue().([]interface{})
	)

	for i, v := range pValues {
		pIndex[v] = i
	}

	var (
		curr    = 0
		inserts []Map
	)

	for _, m := range maps {
		if pChange, changed := m[pField]; changed {
			// update
			pID, ok := pIndex[pChange]
			if !ok {
				panic("rel: cannot update has many assoc that is not loaded or doesn't belong to this record")
			}

			if pID != curr {
				col.Swap(pID, curr)
				pValues[pID], pValues[curr] = pValues[curr], pValues[pID]
			}

			muts[curr] = Apply(col.Get(curr), m)
			delete(pIndex, pChange)
			curr++
		} else {
			inserts = append(inserts, m)
		}
	}

	// delete stales
	if curr < col.Len() {
		deletedIDs = pValues[curr:]
		col.Truncate(0, curr)
	} else {
		deletedIDs = []interface{}{}
	}

	// inserts remaining
	for i, m := range inserts {
		muts[curr+i] = Apply(col.Add(), m)
	}

	return muts, deletedIDs

}

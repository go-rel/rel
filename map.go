package rel

import (
	"fmt"
)

// Map can be used as modification for repository insert or update operation.
// This allows inserting or updating only on specified field.
// Insert/Update of has one or belongs to can be done using other Map as a value.
// Insert/Update of has many can be done using slice of Map as a value.
type Map map[string]interface{}

// Apply modification.
func (m Map) Apply(doc *Document, modification *Modification) {
	for field, value := range m {
		switch v := value.(type) {
		case Map:
			var (
				assoc = doc.Association(field)
			)

			if assoc.Type() != HasOne && assoc.Type() != BelongsTo {
				panic(fmt.Sprint("rel: cannot associate has many", v, "as", field, "into", doc.Table()))
			}

			var (
				assocDoc, _       = assoc.Document()
				assocModification = Apply(assocDoc, v)
			)

			modification.SetAssoc(field, assocModification)
		case []Map:
			var (
				mods   = make([]Modification, len(v))
				assoc  = doc.Association(field)
				col, _ = assoc.Collection()
			)

			col.Reset()

			for i := range v {
				mods[i] = Apply(col.Add(), v[i])
			}
			modification.SetAssoc(field, mods...)
		default:
			if !doc.SetValue(field, v) {
				panic(fmt.Sprint("rel: cannot assign", v, "as", field, "into", doc.Table()))
			}

			modification.SetValue(field, v)
		}
	}
}

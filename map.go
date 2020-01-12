package rel

import (
	"fmt"
)

// Map can be used as changes for repository insert or update operation.
// This allows inserting or updating only on specified field.
// Insert/Update of has one or belongs to can be done using other Map as a value.
// Insert/Update of has many can be done using slice of Map as a value.
type Map map[string]interface{}

// Apply changes.
func (m Map) Apply(doc *Document, changes *Changes) {
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
				assocDoc, _  = assoc.Document()
				assocChanges = ApplyChanges(assocDoc, v)
			)

			changes.SetAssoc(field, assocChanges)
		case []Map:
			var (
				chs    = make([]Changes, len(v))
				assoc  = doc.Association(field)
				col, _ = assoc.Collection()
				doc    = col.Add() // Note: it's reseted again in has to many
			)

			for i := range v {
				chs[i] = ApplyChanges(doc, v[i])
			}
			changes.SetAssoc(field, chs...)
			changes.reload = true // TODO: optimistic create/update, also wrong reload.
		default:
			if !doc.SetValue(field, v) {
				panic(fmt.Sprint("rel: cannot assign", v, "as", field, "into", doc.Table()))
			}

			changes.SetValue(field, v)
		}
	}
}

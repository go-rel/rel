package grimoire

import (
	"reflect"
)

type Changeset struct {
	doc    *Document
	fields map[string]int
	types  []reflect.Type
	values []interface{}
}

type emptyChangeset struct{}

func (c Changeset) Build(changes *Changes) {
	var (
		values = c.doc.Values()
	)

	for f, i := range c.fields {
		var (
			cur = values[i]
		)

		switch v := c.values[i].(type) {
		case Changeset:
			changes.SetAssoc(f, BuildChanges(v))
		case map[interface{}]Changeset:
			c.buildAssocMany(f, changes, v)
		case emptyChangeset:
			// do nothing
		default:
			if c.types[i].Comparable() && v == cur {
				continue
			} else if reflect.DeepEqual(v, cur) {
				continue
			} else {
				changes.SetValue(f, cur)
			}
		}
	}
}

func (c Changeset) buildAssocMany(field string, changes *Changes, changemap map[interface{}]Changeset) {
	var (
		assoc         = c.doc.Association(field)
		col, _        = assoc.Collection()
		lenght        = col.Len()
		chs           = make([]Changes, lenght)
		staleIDs      = []interface{}{}
		unstaleIDsMap = make(map[interface{}]struct{}, lenght)
	)

	for i := range chs {
		var (
			doc    = col.Get(i)
			pField = doc.PrimaryField()
			pValue = doc.PrimaryValue()
		)

		if isZero(pValue) {
			chs[i] = BuildChanges(newStructset(doc))
		} else if cs, ok := changemap[pValue]; ok {
			chs[i] = BuildChanges(cs)
			chs[i].SetValue(pField, pValue)
			unstaleIDsMap[pValue] = struct{}{}
		} else {
			panic("grimoire: cannot update unloaded association")
		}
	}
	changes.SetAssoc(field, chs...)

	// add stale ids
	for id := range changemap {
		if _, ok := unstaleIDsMap[id]; !ok {
			staleIDs = append(staleIDs, id)
		}
	}

	changes.SetStaleAssoc(field, staleIDs)
}

// stores old document values
func newChangeset(doc *Document) Changeset {
	var (
		c = Changeset{
			doc:    doc,
			fields: doc.Fields(),
			types:  doc.Types(),
			values: doc.Values(),
		}
	)

	// replace assoc values
	for _, f := range append(doc.BelongsTo(), doc.HasOne()...) {
		var (
			assoc = doc.Association(f)
		)

		if doc, loaded := assoc.Document(); loaded {
			c.values[c.fields[f]] = newChangeset(doc)
		} else {
			c.values[c.fields[f]] = emptyChangeset{}
		}
	}

	for _, f := range doc.HasMany() {
		if col, loaded := doc.Association(f).Collection(); loaded {
			var (
				docCount  = col.Len()
				changemap = make(map[interface{}]Changeset, docCount)
			)

			for i := 0; i < docCount; i++ {
				var (
					doc    = col.Get(i)
					pValue = doc.PrimaryValue()
				)

				if !isZero(pValue) {
					changemap[pValue] = newChangeset(doc)
				}
			}

			c.values[c.fields[f]] = changemap
		} else {
			c.values[c.fields[f]] = emptyChangeset{}
		}
	}

	return c
}

func NewChangeset(record interface{}) Changeset {
	return newChangeset(newDocument(record))
}

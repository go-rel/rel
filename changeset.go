package rel

import (
	"reflect"
)

type Changeset struct {
	doc    *Document
	fields []string
	types  []reflect.Type
	values []interface{}

	// assocs array is ordered based on the following entry.
	belongsTo []string
	hasOne    []string
	hasMany   []string
	assocs    []interface{}
}

type emptyChangeset struct{}

func (c Changeset) Build(changes *Changes) {
	for i, field := range c.fields {
		var (
			prev   = c.values[i]
			cur, _ = c.doc.Value(field)
		)

		if c.types[i].Comparable() && prev == cur {
			continue
		} else if reflect.DeepEqual(prev, cur) {
			continue
		} else {
			changes.SetValue(field, cur)
		}
	}

	offset := 0
	for i, f := range c.belongsTo {
		if ac, ok := c.assocs[offset+i].(Changeset); ok {
			if ch := BuildChanges(ac); !ch.Empty() {
				changes.SetAssoc(f, ch)
			}
		}
	}

	offset = len(c.belongsTo)
	for i, f := range c.hasOne {
		if ac, ok := c.assocs[offset+i].(Changeset); ok {
			if ch := BuildChanges(ac); !ch.Empty() {
				changes.SetAssoc(f, ch)
			}
		}
	}

	offset = len(c.hasOne)
	for i, f := range c.hasMany {
		if acm, ok := c.assocs[offset+i].(map[interface{}]Changeset); ok {
			c.buildAssocMany(f, changes, acm)
		}
	}
}

func (c Changeset) buildAssocMany(field string, changes *Changes, changemap map[interface{}]Changeset) {
	var (
		assoc         = c.doc.Association(field)
		col, _        = assoc.Collection()
		colCount      = col.Len()
		chs           = []Changes{}
		unstaleIDsMap = make(map[interface{}]struct{})
		staleIDs      = []interface{}{}
	)

	for i := 0; i < colCount; i++ {
		var (
			doc    = col.Get(i)
			pField = doc.PrimaryField()
			pValue = doc.PrimaryValue()
		)

		if isZero(pValue) {
			chs = append(chs, BuildChanges(newStructset(doc)))
		} else if cs, ok := changemap[pValue]; ok {
			if ch := BuildChanges(cs); !ch.Empty() {
				ch.SetValue(pField, pValue)
				chs = append(chs, ch)
			}
			unstaleIDsMap[pValue] = struct{}{}
		} else {
			panic("rel: cannot update unloaded association")
		}
	}

	if len(chs) > 0 {
		changes.SetAssoc(field, chs...)
	}

	// add stale ids
	for id := range changemap {
		if _, ok := unstaleIDsMap[id]; !ok {
			staleIDs = append(staleIDs, id)
		}
	}

	if len(staleIDs) > 0 {
		changes.SetStaleAssoc(field, staleIDs)
	}
}

// stores old document values
func newChangeset(doc *Document) Changeset {
	var (
		fields    = doc.Fields()
		belongsTo = doc.BelongsTo()
		hasOne    = doc.HasOne()
		hasMany   = doc.HasMany()
		changeset = Changeset{
			doc:       doc,
			fields:    fields,
			types:     make([]reflect.Type, len(fields)),
			values:    make([]interface{}, len(fields)),
			belongsTo: belongsTo,
			hasOne:    hasOne,
			hasMany:   hasMany,
			assocs:    make([]interface{}, len(belongsTo)+len(hasOne)+len(hasMany)),
		}
	)

	for i, field := range fields {
		changeset.values[i], _ = doc.Value(field)
		changeset.types[i], _ = doc.Type(field)
	}

	offset := 0
	for i, f := range belongsTo {
		if doc, loaded := doc.Association(f).Document(); loaded {
			changeset.assocs[offset+i] = newChangeset(doc)
		}
	}

	offset = len(belongsTo)
	for i, f := range hasOne {
		if doc, loaded := doc.Association(f).Document(); loaded {
			changeset.assocs[offset+i] = newChangeset(doc)
		}
	}

	offset = len(hasOne)
	for i, f := range hasMany {
		if col, loaded := doc.Association(f).Collection(); loaded {
			var (
				docCount  = col.Len()
				changemap = make(map[interface{}]Changeset, docCount)
			)

			for j := 0; j < docCount; j++ {
				var (
					doc    = col.Get(j)
					pValue = doc.PrimaryValue()
				)

				if !isZero(pValue) {
					changemap[pValue] = newChangeset(doc)
				}
			}

			changeset.assocs[offset+i] = changemap
		}
	}

	return changeset
}

func NewChangeset(record interface{}) Changeset {
	return newChangeset(newDocument(record))
}

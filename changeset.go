package grimoire

import (
	"reflect"
)

type Changeset struct {
	doc    Document
	fields map[string]int
	types  []reflect.Type
	values []interface{}
}

func (c Changeset) Build(changes *Changes) {
	var (
		values    = c.doc.Values()
		structset = newStructset(c.doc)
	)

	for f, i := range c.fields {
		var (
			cur = values[i]
		)

		switch v := c.values[i].(type) {
		case Changeset:
			changes.SetAssoc(f, BuildChanges(v))
		case []Changeset:
			structset.buildAssocMany(f, changes)
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

// stores old document values
func newChangeset(doc Document) Changeset {
	var (
		c = Changeset{
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

		if col, loaded := assoc.Target(); loaded {
			c.values[c.fields[f]] = newChangeset(col.Get(0))
		}
	}

	// don't track has many
	for _, f := range doc.HasMany() {
		c.values[c.fields[f]] = []Changeset{}
	}

	return c
}

func NewChangeset(record interface{}) Changeset {
	return newChangeset(newDocument(record))
}

package rel

import (
	"reflect"
	"time"
)

type pair = [2]interface{}

// Changeset modifier for structs.
// This allows REL to efficiently to perform update operation only on updated fields and association.
// The catch is, enabling changeset will duplicates the original struct values which consumes more memory.
type Changeset struct {
	doc       *Document
	snapshot  []interface{}
	assoc     map[string]Changeset
	assocMany map[string]map[interface{}]Changeset
}

func (c Changeset) changed(typ reflect.Type, old interface{}, new interface{}) bool {
	if oeq, ok := old.(interface{ Equal(interface{}) bool }); ok {
		return !oeq.Equal(new)
	}

	if ot, ok := old.(time.Time); ok {
		return !ot.Equal(new.(time.Time))
	}

	return !(typ.Comparable() && old == new)
}

// Changes returns map of changes, with field names as the keys and an array of old and new values.
// TODO: also returns assoc changes.
func (c Changeset) Changes() map[string]pair {
	changes := make(map[string]pair)

	for i, field := range c.doc.Fields() {
		var (
			typ, _ = c.doc.Type(field)
			old    = c.snapshot[i]
			new, _ = c.doc.Value(field)
		)

		if c.changed(typ, old, new) {
			changes[field] = pair{old, new}
		}
	}

	return changes
}

// Apply modification.
func (c Changeset) Apply(doc *Document, mod *Modification) {
	var (
		t = now().Truncate(time.Second)
	)

	for i, field := range c.doc.Fields() {
		var (
			typ, _ = c.doc.Type(field)
			old    = c.snapshot[i]
			new, _ = c.doc.Value(field)
		)

		if c.changed(typ, old, new) {
			mod.Add(Set(field, new))
		}
	}

	if len(mod.Modifies) > 0 && c.doc.Flag(HasUpdatedAt) && c.doc.SetValue("updated_at", t) {
		mod.Add(Set("updated_at", t))
	}

	for _, field := range doc.BelongsTo() {
		c.applyAssoc(field, mod)
	}

	for _, field := range doc.HasOne() {
		c.applyAssoc(field, mod)
	}

	for _, field := range doc.HasMany() {
		c.applyAssocMany(field, mod)
	}
}

func (c Changeset) applyAssoc(field string, mod *Modification) {
	assoc := c.doc.Association(field)
	if assoc.IsZero() {
		return
	}

	doc, _ := assoc.Document()

	if ch, ok := c.assoc[field]; ok {
		if amod := Apply(doc, ch); len(amod.Modifies) > 0 || len(amod.Assoc) > 0 {
			mod.SetAssoc(field, amod)
		}
	} else {
		amod := Apply(doc, newStructset(doc, false))
		mod.SetAssoc(field, amod)
	}
}

func (c Changeset) applyAssocMany(field string, mod *Modification) {
	if dirties, ok := c.assocMany[field]; ok {
		var (
			assoc      = c.doc.Association(field)
			col, _     = assoc.Collection()
			mods       = make([]Modification, 0, col.Len())
			updatedIDs = make(map[interface{}]struct{})
			deletedIDs []interface{}
		)

		for i := 0; i < col.Len(); i++ {
			var (
				doc    = col.Get(i)
				pValue = doc.PrimaryValue()
			)

			if ch, ok := dirties[pValue]; ok {
				updatedIDs[pValue] = struct{}{}

				if amod := Apply(doc, ch); len(amod.Modifies) > 0 || len(amod.Assoc) > 0 {
					mods = append(mods, amod)
				}
			} else {
				mods = append(mods, Apply(doc, newStructset(doc, false)))
			}
		}

		// leftover snapshot.
		if len(updatedIDs) != len(dirties) {
			for i := range dirties {
				if _, ok := updatedIDs[i]; !ok {
					deletedIDs = append(deletedIDs, i)
				}
			}
		}

		if len(mods) > 0 || len(deletedIDs) > 0 {
			mod.SetAssoc(field, mods...)
			mod.SetDeletedIDs(field, deletedIDs)
		}
	} else {
		newStructset(c.doc, false).buildAssocMany(field, mod)
	}
}

// NewChangeset returns new changeset modifier for given record.
func NewChangeset(record interface{}) Changeset {
	// TODO: make it readonly
	return newChangeset(NewDocument(record, false))
}

func newChangeset(doc *Document) Changeset {
	c := Changeset{
		doc:       doc,
		snapshot:  make([]interface{}, len(doc.Fields())),
		assoc:     make(map[string]Changeset),
		assocMany: make(map[string]map[interface{}]Changeset),
	}

	for i, field := range doc.Fields() {
		c.snapshot[i], _ = doc.Value(field)
	}

	for _, field := range doc.BelongsTo() {
		initChangesetAssoc(doc, c.assoc, field)
	}

	for _, field := range doc.HasOne() {
		initChangesetAssoc(doc, c.assoc, field)
	}

	for _, field := range doc.HasMany() {
		initChangesetAssocMany(doc, c.assocMany, field)
	}

	return c
}

func initChangesetAssoc(doc *Document, assoc map[string]Changeset, field string) {
	doc, loaded := doc.Association(field).Document()
	if !loaded {
		return
	}

	assoc[field] = newChangeset(doc)
}

func initChangesetAssocMany(doc *Document, assoc map[string]map[interface{}]Changeset, field string) {
	col, loaded := doc.Association(field).Collection()
	if !loaded {
		return
	}

	assoc[field] = make(map[interface{}]Changeset)

	for i := 0; i < col.Len(); i++ {
		var (
			doc    = col.Get(i)
			pValue = doc.PrimaryValue()
		)

		if !isZero(pValue) {
			assoc[field][pValue] = newChangeset(doc)
		}
	}
}

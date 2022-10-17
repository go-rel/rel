package rel

import (
	"bytes"
	"reflect"
	"time"
)

type pair = [2]any

// Changeset mutator for structs.
// This allows REL to efficiently to perform update operation only on updated fields and association.
// The catch is, enabling changeset will duplicates the original struct values which consumes more memory.
type Changeset struct {
	doc       *Document
	snapshot  []any
	assoc     map[string]Changeset
	assocMany map[string]map[any]Changeset
}

func (c Changeset) valueChanged(typ reflect.Type, old any, new any) bool {
	if oeq, ok := old.(interface{ Equal(any) bool }); ok {
		return !oeq.Equal(new)
	}

	if ot, ok := old.(time.Time); ok {
		return !ot.Equal(new.(time.Time))
	}

	if typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Uint8 {
		return !bytes.Equal(reflect.ValueOf(old).Bytes(), reflect.ValueOf(new).Bytes())
	}

	return !(typ.Comparable() && old == new)
}

// FieldChanged returns true if field exists and it's already changed.
// returns false otherwise.
func (c Changeset) FieldChanged(field string) bool {
	for i, f := range c.doc.Fields() {
		if f == field {
			var (
				typ, _ = c.doc.Type(field)
				old    = c.snapshot[i]
				new, _ = c.doc.Value(field)
			)

			return c.valueChanged(typ, old, new)
		}
	}

	return false
}

// Changes returns map of changes.
func (c Changeset) Changes() map[string]any {
	return buildChanges(c.doc, c)
}

// Apply mutation.
func (c Changeset) Apply(doc *Document, mut *Mutation) {
	var (
		t = Now()
	)

	for i, field := range c.doc.Fields() {
		var (
			typ, _ = c.doc.Type(field)
			old    = c.snapshot[i]
			new, _ = c.doc.Value(field)
		)

		if c.valueChanged(typ, old, new) {
			mut.Add(Set(field, new))
		}
	}

	if !mut.IsMutatesEmpty() && c.doc.Flag(HasUpdatedAt) && c.doc.SetValue("updated_at", t) {
		mut.Add(Set("updated_at", t))
	}

	if mut.Cascade {
		for _, field := range doc.BelongsTo() {
			c.applyAssoc(field, mut)
		}

		for _, field := range doc.HasOne() {
			c.applyAssoc(field, mut)
		}

		for _, field := range doc.HasMany() {
			c.applyAssocMany(field, mut)
		}
	}
}

func (c Changeset) applyAssoc(field string, mut *Mutation) {
	assoc := c.doc.Association(field)
	if assoc.IsZero() {
		return
	}

	doc, _ := assoc.Document()

	if ch, ok := c.assoc[field]; ok {
		if amod := Apply(doc, ch); !amod.IsEmpty() {
			mut.SetAssoc(field, amod)
		}
	} else {
		amod := Apply(doc, newStructset(doc, false))
		mut.SetAssoc(field, amod)
	}
}

func (c Changeset) applyAssocMany(field string, mut *Mutation) {
	if chs, ok := c.assocMany[field]; ok {
		var (
			assoc      = c.doc.Association(field)
			col, _     = assoc.Collection()
			muts       = make([]Mutation, 0, col.Len())
			updatedIDs = make(map[any]struct{})
			deletedIDs []any
		)

		for i := 0; i < col.Len(); i++ {
			var (
				doc    = col.Get(i)
				pValue = doc.PrimaryValue()
			)

			if ch, ok := chs[pValue]; ok {
				updatedIDs[pValue] = struct{}{}

				if amod := Apply(doc, ch); !amod.IsEmpty() {
					muts = append(muts, amod)
				}
			} else {
				muts = append(muts, Apply(doc, newStructset(doc, false)))
			}
		}

		// leftover snapshot.
		if len(updatedIDs) != len(chs) {
			for id := range chs {
				if _, ok := updatedIDs[id]; !ok {
					deletedIDs = append(deletedIDs, id)
				}
			}
		}

		if len(muts) > 0 || len(deletedIDs) > 0 {
			mut.SetAssoc(field, muts...)
			mut.SetDeletedIDs(field, deletedIDs)
		}
	} else {
		newStructset(c.doc, false).buildAssocMany(field, mut)
	}
}

// NewChangeset returns new changeset mutator for given entity.
func NewChangeset(entity any) Changeset {
	return newChangeset(NewDocument(entity))
}

func newChangeset(doc *Document) Changeset {
	c := Changeset{
		doc:       doc,
		snapshot:  make([]any, len(doc.Fields())),
		assoc:     make(map[string]Changeset),
		assocMany: make(map[string]map[any]Changeset),
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

func initChangesetAssocMany(doc *Document, assoc map[string]map[any]Changeset, field string) {
	col, loaded := doc.Association(field).Collection()
	if !loaded {
		return
	}

	assoc[field] = make(map[any]Changeset)

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

func buildChanges(doc *Document, c Changeset) map[string]any {
	var (
		changes = make(map[string]any)
		fields  []string
	)

	if doc != nil {
		fields = doc.Fields()
	} else {
		fields = c.doc.Fields()
	}

	for i, field := range fields {
		switch {
		case doc == nil:
			if old := c.snapshot[i]; old != nil {
				changes[field] = pair{old, nil}
			}
		case i >= len(c.snapshot):
			if new, _ := doc.Value(field); new != nil {
				changes[field] = pair{nil, new}
			}
		default:
			old := c.snapshot[i]
			new, _ := doc.Value(field)
			if typ, _ := doc.Type(field); c.valueChanged(typ, old, new) {
				changes[field] = pair{old, new}
			}
		}
	}

	if doc == nil || len(c.snapshot) == 0 {
		return changes
	}

	for _, field := range doc.BelongsTo() {
		buildChangesAssoc(changes, c, field)
	}

	for _, field := range doc.HasOne() {
		buildChangesAssoc(changes, c, field)
	}

	for _, field := range doc.HasMany() {
		buildChangesAssocMany(changes, c, field)
	}

	return changes
}

func buildChangesAssoc(out map[string]any, c Changeset, field string) {
	assoc := c.doc.Association(field)
	if assoc.IsZero() {
		return
	}

	doc, _ := assoc.Document()
	if changes := buildChanges(doc, c.assoc[field]); len(changes) != 0 {
		out[field] = changes
	}
}

func buildChangesAssocMany(out map[string]any, c Changeset, field string) {
	var (
		changes    []map[string]any
		chs        = c.assocMany[field]
		assoc      = c.doc.Association(field)
		col, _     = assoc.Collection()
		updatedIDs = make(map[any]struct{})
	)

	for i := 0; i < col.Len(); i++ {
		var (
			doc          = col.Get(i)
			pValue       = doc.PrimaryValue()
			ch, isUpdate = chs[pValue]
		)

		if isUpdate {
			updatedIDs[pValue] = struct{}{}
		}

		if dChanges := buildChanges(doc, ch); len(dChanges) != 0 {
			changes = append(changes, dChanges)
		}
	}

	// leftover snapshot.
	if len(updatedIDs) != len(chs) {
		for id, ch := range chs {
			if _, ok := updatedIDs[id]; !ok {
				changes = append(changes, buildChanges(nil, ch))
			}
		}
	}

	if len(changes) != 0 {
		out[field] = changes
	}
}

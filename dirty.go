package rel

import (
	"reflect"
	"time"
)

// Dirty tracking for golang structs.
// This allows REL to efficiently to perform update operation only on updated fields and association.
// The catch is, enabling dirty will duplicates the original struct values which consume more memory.
// Dirty tracking will automatically initialized by Find/FindAll/FindAndCountAll/Preload methods when embedded.
type Dirty struct {
	doc       *Document
	snapshot  []interface{}
	assoc     map[string]*Dirty
	assocMany map[string]map[interface{}]*Dirty
}

// init dirty states.
func (d *Dirty) init(doc *Document) {
	d.doc = doc
	d.snapshot = make([]interface{}, len(doc.Fields()))
	d.assoc = make(map[string]*Dirty)
	d.assocMany = make(map[string]map[interface{}]*Dirty)

	for i, field := range doc.Fields() {
		d.snapshot[i], _ = doc.Value(field)
	}

	for _, field := range doc.BelongsTo() {
		d.initAssoc(field)
	}

	for _, field := range doc.HasOne() {
		d.initAssoc(field)
	}

	for _, field := range doc.HasMany() {
		d.initAssocMany(field)
	}
}

func (d *Dirty) initAssoc(field string) {
	var (
		assoc = d.doc.Association(field)
	)

	if doc, loaded := assoc.Document(); loaded {
		dirty := doc.Dirty()
		if dirty == nil {
			dirty = &Dirty{}
			dirty.init(doc)
		}

		d.assoc[field] = dirty
	}
}

func (d *Dirty) initAssocMany(field string) {
	var (
		assoc = d.doc.Association(field)
	)

	if col, loaded := assoc.Collection(); loaded {
		d.assocMany[field] = make(map[interface{}]*Dirty)

		for i := 0; i < col.Len(); i++ {
			var (
				doc    = col.Get(i)
				pValue = doc.PrimaryValue()
			)

			if !isZero(pValue) {
				dirty := doc.Dirty()
				if dirty == nil {
					dirty = &Dirty{}
					dirty.init(doc)
				}

				d.assocMany[field][pValue] = dirty
			}
		}
	}
}

func (d *Dirty) changed(typ reflect.Type, old interface{}, new interface{}) bool {
	if oeq, ok := old.(interface{ Equal(interface{}) bool }); ok {
		return !oeq.Equal(new)
	}

	if ot, ok := old.(time.Time); ok {
		return !ot.Equal(new.(time.Time))
	}

	return !(typ.Comparable() && old == new)
}

// Changes returns map of changes, with field names as the keys and an array of old and new values.
func (d Dirty) Changes() map[string][2]interface{} {
	changes := make(map[string][2]interface{})

	for i, field := range d.doc.Fields() {
		var (
			typ, _ = d.doc.Type(field)
			old    = d.snapshot[i]
			new, _ = d.doc.Value(field)
		)

		if d.changed(typ, old, new) {
			changes[field] = [2]interface{}{old, new}
		}
	}

	return changes
}

// Apply modification.
func (d Dirty) Apply(doc *Document, mod *Modification) {
	// Not initialized, fallback to structset.
	if len(d.snapshot) == 0 {
		newStructset(doc, false).Apply(doc, mod)
		return
	}

	var (
		t = now().Truncate(time.Second)
	)

	for i, field := range d.doc.Fields() {
		var (
			typ, _ = d.doc.Type(field)
			old    = d.snapshot[i]
			new, _ = d.doc.Value(field)
		)

		if d.changed(typ, old, new) {
			mod.Add(Set(field, new))
		}
	}

	if len(mod.Modifies) > 0 && d.doc.Flag(HasUpdatedAt) && d.doc.SetValue("updated_at", t) {
		mod.Add(Set("updated_at", t))
	}

	for _, field := range doc.BelongsTo() {
		d.applyAssoc(field, mod)
	}

	for _, field := range doc.HasOne() {
		d.applyAssoc(field, mod)
	}

	for _, field := range doc.HasMany() {
		d.applyAssocMany(field, mod)
	}

	// TODO: reinitialize dirty
	// d.init(d.doc)
}

func (d Dirty) applyAssoc(field string, mod *Modification) {
	if dirty, ok := d.assoc[field]; ok {
		assoc := d.doc.Association(field)
		if assoc.IsZero() {
			return
		}

		doc, _ := assoc.Document()
		mod.SetAssoc(field, Apply(doc, dirty))
	} else {
		newStructset(d.doc, false).buildAssoc(field, mod)
	}
}

func (d Dirty) applyAssocMany(field string, mod *Modification) {
	if dirties, ok := d.assocMany[field]; ok {
		assoc := d.doc.Association(field)
		if assoc.IsZero() {
			return
		}

		var (
			col, _ = assoc.Collection()
			mods   = make([]Modification, col.Len())
		)

		for i := 0; i < col.Len(); i++ {
			var (
				doc    = col.Get(i)
				pValue = doc.PrimaryValue()
			)

			if dirty, ok := dirties[pValue]; ok {
				mods[i] = Apply(doc, dirty)
			} else {
				mods[i] = Apply(doc, newStructset(doc, false))
			}
		}
	} else {
		newStructset(d.doc, false).buildAssocMany(field, mod)
	}
}

package rel

import (
	"fmt"
	"time"
)

var (
	now = time.Now
)

// Structset can be used as modification for repository insert or update operation.
// This will save every field in struct and it's association as long as it's loaded.
// This is the default modifier used by repository.
type Structset struct {
	doc *Document
}

// Apply modification.
func (s Structset) Apply(doc *Document, mod *Modification) {
	var (
		pField = s.doc.PrimaryField()
		t      = now().Truncate(time.Second)
	)

	for _, field := range s.doc.Fields() {
		switch field {
		case pField:
			continue
		case "created_at", "inserted_at":
			if typ, ok := doc.Type(field); ok && typ == rtTime {
				if value, ok := doc.Value(field); ok && value.(time.Time).IsZero() {
					s.set(doc, mod, field, t, true)
				}
				continue
			}
		case "updated_at":
			if typ, ok := doc.Type(field); ok && typ == rtTime {
				s.set(doc, mod, field, t, true)
				continue
			}
		}

		if value, ok := s.doc.Value(field); ok && !isZero(value) {
			s.set(doc, mod, field, value, false)
		}
	}

	for _, field := range s.doc.BelongsTo() {
		s.buildAssoc(field, mod)
	}

	for _, field := range s.doc.HasOne() {
		s.buildAssoc(field, mod)
	}

	for _, field := range s.doc.HasMany() {
		s.buildAssocMany(field, mod)
	}
}

func (s Structset) set(doc *Document, mod *Modification, field string, value interface{}, force bool) {
	if (force || doc.v != s.doc.v) && !doc.SetValue(field, value) {
		panic(fmt.Sprint("rel: cannot assign ", value, " as ", field, " into ", doc.Table()))
	}

	mod.SetValue(field, value)
}

func (s Structset) buildAssoc(field string, mod *Modification) {
	var (
		assoc = s.doc.Association(field)
	)

	if !assoc.IsZero() {
		var (
			doc, _ = assoc.Document()
		)

		mod.SetAssoc(field, Apply(doc, Structset{doc: doc}))
	}
}

func (s Structset) buildAssocMany(field string, mod *Modification) {
	var (
		assoc = s.doc.Association(field)
	)

	if !assoc.IsZero() {
		var (
			col, _ = assoc.Collection()
			mods   = make([]Modification, col.Len())
		)

		for i := range mods {
			var (
				doc = col.Get(i)
			)

			mods[i] = Apply(doc, newStructset(doc))
		}

		mod.SetAssoc(field, mods...)
	}
}

func newStructset(doc *Document) Structset {
	return Structset{
		doc: doc,
	}
}

// NewStructset from a struct.
func NewStructset(record interface{}) Structset {
	return newStructset(NewDocument(record))
}

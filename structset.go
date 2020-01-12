package rel

import (
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
// TODO: apply modification if it's applying to struct.
func (s Structset) Apply(doc *Document, modification *Modification) {
	var (
		pField = s.doc.PrimaryField()
		t      = now()
	)

	for _, field := range s.doc.Fields() {
		switch field {
		case pField:
			continue
		case "created_at", "inserted_at":
			// TODO: handle *time.Time
			if typ, ok := s.doc.Type(field); ok && typ == rtTime {
				if value, ok := s.doc.Value(field); ok && value.(time.Time).IsZero() {
					modification.SetValue(field, t)
				}
				continue
			}
		case "updated_at":
			if typ, ok := s.doc.Type(field); ok && typ == rtTime {
				modification.SetValue(field, t)
				continue
			}
		}

		if value, ok := s.doc.Value(field); ok && !isZero(value) {
			modification.SetValue(field, value)
		}
	}

	for _, field := range s.doc.BelongsTo() {
		s.buildAssoc(field, modification)
	}

	for _, field := range s.doc.HasOne() {
		s.buildAssoc(field, modification)
	}

	for _, field := range s.doc.HasMany() {
		s.buildAssocMany(field, modification)
	}
}

func (s Structset) buildAssoc(field string, modification *Modification) {
	var (
		assoc = s.doc.Association(field)
	)

	if !assoc.IsZero() {
		var (
			doc, _ = assoc.Document()
			mod    = Apply(doc, Structset{doc: doc})
		)

		modification.SetAssoc(field, mod)
	}
}

func (s Structset) buildAssocMany(field string, modification *Modification) {
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

		modification.SetAssoc(field, mods...)
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

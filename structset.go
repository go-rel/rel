package rel

import (
	"time"
)

var (
	now = time.Now
)

// Structset can be used as changes for repository insert or update operation.
// This will save every field in struct and it's association as long as it's loaded.
// This is the default changer used by repository.
type Structset struct {
	doc *Document
}

// Apply changes.
func (s Structset) Apply(doc *Document, changes *Changes) error {
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
					changes.SetValue(field, t)
				}
				continue
			}
		case "updated_at":
			if typ, ok := s.doc.Type(field); ok && typ == rtTime {
				changes.SetValue(field, t)
				continue
			}
		}

		if value, ok := s.doc.Value(field); ok && !isZero(value) {
			changes.SetValue(field, value)
		}
	}

	for _, field := range s.doc.BelongsTo() {
		s.buildAssoc(field, changes)
	}

	for _, field := range s.doc.HasOne() {
		s.buildAssoc(field, changes)
	}

	for _, field := range s.doc.HasMany() {
		s.buildAssocMany(field, changes)
	}

	return nil
}

func (s Structset) buildAssoc(field string, changes *Changes) {
	var (
		assoc = s.doc.Association(field)
	)

	if !assoc.IsZero() {
		var (
			doc, _ = assoc.Document()
			ch, _  = ApplyChanges(doc, Structset{doc: doc})
		)

		changes.SetAssoc(field, ch)
	}
}

func (s Structset) buildAssocMany(field string, changes *Changes) {
	var (
		assoc = s.doc.Association(field)
	)

	if !assoc.IsZero() {
		var (
			col, _ = assoc.Collection()
			chs    = make([]Changes, col.Len())
		)

		for i := range chs {
			var (
				doc = col.Get(i)
			)

			chs[i], _ = ApplyChanges(doc, newStructset(doc))
		}

		changes.SetAssoc(field, chs...)
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

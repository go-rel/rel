package rel

import (
	"time"
)

var (
	now = time.Now
)

type Structset struct {
	doc *document
}

func (s Structset) Build(changes *Changes) {
	var (
		pField = s.doc.PrimaryField()
		t      = now()
	)

	for _, field := range s.doc.Fields() {
		switch field {
		case pField:
			continue
		case "created_at", "inserted_at":
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
}

func (s Structset) buildAssoc(field string, changes *Changes) {
	var (
		assoc = s.doc.Association(field)
	)

	if !assoc.IsZero() {
		var (
			doc, _ = assoc.Document()
			ch     = BuildChanges(Structset{doc: doc})
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
			chs[i] = BuildChanges(newStructset(col.Get(i)))
		}

		changes.SetAssoc(field, chs...)
	}
}

func newStructset(doc *document) Structset {
	return Structset{
		doc: doc,
	}
}

func NewStructset(record interface{}) Structset {
	return newStructset(newDocument(record))
}

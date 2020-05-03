package rel

import (
	"fmt"
	"time"
)

var (
	now = time.Now
)

// Structset can be used as mutation for repository insert or update operation.
// This will save every field in struct and it's association as long as it's loaded.
// This is the default mutator used by repository.
type Structset struct {
	doc      *Document
	skipZero bool
}

// Apply mutation.
func (s Structset) Apply(doc *Document, mod *Mutation) {
	var (
		pField = s.doc.PrimaryField()
		t      = now().Truncate(time.Second)
	)

	for _, field := range s.doc.Fields() {
		switch field {
		case pField:
			continue
		case "created_at", "inserted_at":
			if doc.Flag(HasCreatedAt) {
				if value, ok := doc.Value(field); ok && value.(time.Time).IsZero() {
					s.set(doc, mod, field, t, true)
					continue
				}
			}
		case "updated_at":
			if doc.Flag(HasUpdatedAt) {
				s.set(doc, mod, field, t, true)
				continue
			}
		}

		s.applyValue(doc, mod, field)
	}

	s.applyAssoc(mod)
}

func (s Structset) set(doc *Document, mod *Mutation, field string, value interface{}, force bool) {
	if (force || doc.v != s.doc.v) && !doc.SetValue(field, value) {
		panic(fmt.Sprint("rel: cannot assign ", value, " as ", field, " into ", doc.Table()))
	}

	mod.Add(Set(field, value))
}

func (s Structset) applyValue(doc *Document, mod *Mutation, field string) {
	if value, ok := s.doc.Value(field); ok {
		if s.skipZero && isZero(value) {
			return
		}

		s.set(doc, mod, field, value, false)
	}
}

func (s Structset) applyAssoc(mod *Mutation) {
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

func (s Structset) buildAssoc(field string, mod *Mutation) {
	assoc := s.doc.Association(field)
	if assoc.IsZero() {
		return
	}

	var (
		doc, _ = assoc.Document()
	)

	mod.SetAssoc(field, Apply(doc, newStructset(doc, s.skipZero)))
}

func (s Structset) buildAssocMany(field string, mod *Mutation) {
	assoc := s.doc.Association(field)
	if assoc.IsZero() {
		return
	}

	var (
		col, _ = assoc.Collection()
		pField = col.PrimaryField()
		mods   = make([]Mutation, col.Len())
	)

	for i := range mods {
		var (
			doc = col.Get(i)
		)

		mods[i] = Apply(doc, newStructset(doc, s.skipZero))
		doc.SetValue(pField, nil) // reset id, since it'll be reinserted.
	}

	mod.SetAssoc(field, mods...)
}

func newStructset(doc *Document, skipZero bool) Structset {
	return Structset{
		doc:      doc,
		skipZero: skipZero,
	}
}

// NewStructset from a struct.
func NewStructset(record interface{}, skipZero bool) Structset {
	return newStructset(NewDocument(record), skipZero)
}

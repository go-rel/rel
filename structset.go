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
func (s Structset) Apply(doc *Document, mut *Mutation) {
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
					s.set(doc, mut, field, t, true)
					continue
				}
			}
		case "updated_at":
			if doc.Flag(HasUpdatedAt) {
				s.set(doc, mut, field, t, true)
				continue
			}
		}

		s.applyValue(doc, mut, field)
	}

	if mut.Cascade {
		s.applyAssoc(mut)
	}
}

func (s Structset) set(doc *Document, mut *Mutation, field string, value interface{}, force bool) {
	if (force || doc.v != s.doc.v) && !doc.SetValue(field, value) {
		panic(fmt.Sprint("rel: cannot assign ", value, " as ", field, " into ", doc.Table()))
	}

	mut.Add(Set(field, value))
}

func (s Structset) applyValue(doc *Document, mut *Mutation, field string) {
	if value, ok := s.doc.Value(field); ok {
		if s.skipZero && isZero(value) {
			return
		}

		s.set(doc, mut, field, value, false)
	}
}

func (s Structset) applyAssoc(mut *Mutation) {
	for _, field := range s.doc.BelongsTo() {
		s.buildAssoc(field, mut)
	}

	for _, field := range s.doc.HasOne() {
		s.buildAssoc(field, mut)
	}

	for _, field := range s.doc.HasMany() {
		s.buildAssocMany(field, mut)
	}
}

func (s Structset) buildAssoc(field string, mut *Mutation) {
	assoc := s.doc.Association(field)
	if assoc.IsZero() {
		return
	}

	var (
		doc, _ = assoc.Document()
	)

	mut.SetAssoc(field, Apply(doc, newStructset(doc, s.skipZero)))
}

func (s Structset) buildAssocMany(field string, mut *Mutation) {
	assoc := s.doc.Association(field)
	if assoc.IsZero() {
		return
	}

	var (
		col, _ = assoc.Collection()
		pField = col.PrimaryField()
		muts   = make([]Mutation, col.Len())
	)

	for i := range muts {
		var (
			doc = col.Get(i)
		)

		muts[i] = Apply(doc, newStructset(doc, s.skipZero))
		doc.SetValue(pField, nil) // reset id, since it'll be reinserted.
	}

	mut.SetAssoc(field, muts...)
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
